package agent

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/calehh/hac-app/state"
	hac_types "github.com/calehh/hac-app/types"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtlog "github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/rpc/client/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Client interface {
	IfProcessProposal(ctx context.Context, proposer uint64, data []byte) (bool, error)
	IfAcceptProposal(ctx context.Context, proposal uint64) (bool, error)
	IfGrantNewMember(ctx context.Context, amount uint64, statement string) (bool, error)
}

var _ Client = &MockClient{}

type MockClient struct {
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) IfAcceptProposal(ctx context.Context, proposal uint64) (bool, error) {
	return true, nil
}

func (m *MockClient) IfGrantNewMember(ctx context.Context, amount uint64, statement string) (bool, error) {
	return true, nil
}

func (m *MockClient) IfProcessProposal(ctx context.Context, proposer uint64, data []byte) (bool, error) {
	return true, nil
}

type ChainIndexer struct {
	logger        cmtlog.Logger
	Url           string
	Height        int64
	db            *gorm.DB
	cli           *http.HTTP
	eventHandlers map[string]eventHandler
}

func NewChainIndexer(logger cmtlog.Logger, dbPath string, url string) (*ChainIndexer, error) {
	logger.Info("NewChainIndexer", "dbPath", dbPath, "url", url)
	cli, err := http.New(url, "/websocket")
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Grant{}, &Discussion{}, &Proposal{}, &Height{}, &GrantVote{}, &ProposalVote{}).Error; err != nil {
		return nil, err
	}
	h := Height{Id: 1}
	if err = db.First(&h).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	c := ChainIndexer{
		logger:        logger.With("module", "indexer"),
		Url:           url,
		Height:        int64(h.Height + 1),
		db:            db,
		cli:           cli,
		eventHandlers: map[string]eventHandler{},
	}
	c.eventHandlers = map[string]eventHandler{
		hac_types.EventGrantType:          c.handleEventGrant,
		hac_types.EventDiscussionType:     c.handleEventDiscussion,
		hac_types.EventSettleProposalType: c.handleEventSettleProposal,
		hac_types.EventProposalType:       c.handleEventProposal,
	}
	return &c, nil
}

type eventHandler func(ctx context.Context, event abci.Event, height int64)

func (c *ChainIndexer) handleEvent(ctx context.Context, event abci.Event, height int64) {
	if h, ok := c.eventHandlers[event.Type]; ok {
		h(ctx, event, height)
	}
}

func (c *ChainIndexer) handleEventGrant(ctx context.Context, event abci.Event, height int64) {
	ev := hac_types.ParseEventGrant(event)
	if ev == nil {
		c.logger.Error("decode event fail", "event", event)
		return
	}
	account := Grant{
		Id:              ev.Validator,
		Address:         ev.Address,
		Height:          uint64(height),
		Stake:           ev.Amount,
		Proposer:        ev.ProposerIndex,
		ProposerAddress: ev.ProposerAddress,
		Grant:           ev.Grant,
	}
	if err := c.db.Save(&account).Error; err != nil {
		c.logger.Error("save account fail", "err", err)
	}
}

func (c *ChainIndexer) handleEventDiscussion(ctx context.Context, event abci.Event, height int64) {
	ev := hac_types.DecodeEventDiscussion(event)
	if ev == nil {
		c.logger.Error("decode event fail", "event", event)
		return
	}
	discusstion := Discussion{
		Proposal:       ev.Proposal,
		SpeakerIndex:   ev.Speaker,
		SpeakerAddress: ev.SpeakerAddress,
		Data:           ev.Data,
		Height:         uint64(height),
	}
	if err := c.db.Save(&discusstion).Error; err != nil {
		c.logger.Error("save discusstion fail", "err", err)
	}
}

func (c *ChainIndexer) handleEventSettleProposal(ctx context.Context, event abci.Event, height int64) {
	ev := hac_types.DecodeEventSettleProposal(event)
	if ev == nil {
		c.logger.Error("decode event fail", "event", event)
		return
	}
	var proposal Proposal
	if err := c.db.First(&proposal, ev.Proposal).Error; err != nil {
		c.logger.Error("get proposal fail", "err", err)
		return
	}
	proposal.Status = uint64(ev.State)
	proposal.SettleHeight = uint64(height)
	if err := c.db.Save(&proposal).Error; err != nil {
		c.logger.Error("save proposal fail", "err", err)
	}
}

func (c *ChainIndexer) handleEventProposal(ctx context.Context, event abci.Event, height int64) {
	ev := hac_types.DecodeEventProposal(event)
	if ev == nil {
		c.logger.Error("decode event fail", "event", event)
		return
	}
	proposal := Proposal{
		Id:              ev.ProposalIndex,
		ProposerIndex:   ev.Proposer,
		ProposerAddress: ev.ProposerAddress,
		Data:            ev.Data,
		NewHeight:       uint64(height),
		Status:          ev.Status,
	}
	if err := c.db.Save(&proposal).Error; err != nil {
		c.logger.Error("save proposal fail", "err", err)
	}
}

func (c *ChainIndexer) handleVote(ctx context.Context, height int64) error {
	res, err := c.cli.Commit(ctx, &height)
	if err != nil {
		c.logger.Error("get Commit fail", "err", err)
		if !c.cli.IsRunning() {
			c.cli.Stop()
			c.cli, err = http.New(c.Url, "/websocket")
			if err != nil {
				c.logger.Error("reconnect fail", "err", err)
				return err
			}
		}
	}
	voteHeight := res.Height
	// new proposal
	newProposel := Proposal{}
	if err := c.db.Where("new_height = ?", voteHeight).First(&newProposel).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}
	if newProposel.Id != 0 {
		for _, v := range res.Commit.Signatures {
			acc, err := c.queryAccount(ctx, 0, v.ValidatorAddress.String())
			if err != nil {
				return err
			}
			if acc == nil {
				return fmt.Errorf("commit sig address not exist address:%s", v.ValidatorAddress.String())
			}
			if err := c.db.Where("height = ? And voter_index = ?", voteHeight, acc.Index).First(&ProposalVote{}).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
				vote := ProposalVote{
					Proposal:     newProposel.Id,
					VoterIndex:   acc.Index,
					VoterAddress: v.ValidatorAddress.String(),
					Height:       uint64(voteHeight),
					Vote:         uint64(v.VoteCode),
				}
				if err := c.db.Create(&vote).Error; err != nil {
					return err
				}
			}
		}
		return nil
	}
	// settle proposal
	settleProposel := Proposal{}
	if err := c.db.Where("settle_height = ?", voteHeight).First(&settleProposel).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}
	if settleProposel.Id != 0 {
		for _, v := range res.Commit.Signatures {
			acc, err := c.queryAccount(ctx, 0, v.ValidatorAddress.String())
			if err != nil {
				return err
			}
			if acc == nil {
				return fmt.Errorf("commit sig address not exist address:%s", v.ValidatorAddress.String())
			}
			if err := c.db.Where("height = ? And voter_index = ?", voteHeight, acc.Index).First(&ProposalVote{}).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
				vote := ProposalVote{
					Proposal:     settleProposel.Id,
					VoterIndex:   acc.Index,
					VoterAddress: v.ValidatorAddress.String(),
					Height:       uint64(voteHeight),
					Vote:         uint64(v.VoteCode),
				}
				if err := c.db.Create(&vote).Error; err != nil {
					return err
				}
			}
		}
		return nil
	}
	// grant grant
	grant := Grant{}
	if err := c.db.Where("height = ?", voteHeight).First(&grant).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	}
	if grant.Id != 0 {
		for _, v := range res.Commit.Signatures {
			acc, err := c.queryAccount(ctx, 0, v.ValidatorAddress.String())
			if err != nil {
				return err
			}
			if acc == nil {
				return fmt.Errorf("commit sig address not exist address:%s", v.ValidatorAddress.String())
			}
			if err := c.db.Where("height = ? And voter_index = ?", voteHeight, acc.Index).First(&GrantVote{}).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					return err
				}
				vote := GrantVote{
					ProposerIndex:   grant.Proposer,
					ProposerAddress: grant.ProposerAddress,
					AccountIndex:    grant.Id,
					AccountAddr:     grant.Address,
					VoterIndex:      acc.Index,
					VoterAddress:    acc.Address(),
					Height:          uint64(voteHeight),
					Vote:            uint64(v.VoteCode),
				}
				if err := c.db.Create(&vote).Error; err != nil {
					return err
				}
			}
		}
		return nil
	}
	return nil
}

func (c *ChainIndexer) Start(ctx context.Context) {
	var err error
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if c.cli == nil {
				c.cli, err = http.New(c.Url, "/websocket")
				if err != nil {
					c.logger.Error("connect fail", "err", err)
					continue
				}
			}
			b, err := c.cli.Status(context.TODO())
			if err != nil {
				c.logger.Error("get status fail", "err", err)
				if !c.cli.IsRunning() {
					c.cli.Stop()
					c.cli, err = http.New(c.Url, "/websocket")
					if err != nil {
						c.logger.Error("reconnect fail", "err", err)
						continue
					}
				}
			}
			for b.SyncInfo.LatestBlockHeight > c.Height {
				time.Sleep(time.Millisecond * 100)
				c.logger.Info("indexer syncing", "height", c.Height)
				events, err := c.cli.BlockResults(ctx, &c.Height)
				if err != nil {
					c.logger.Error("get status fail", "err", err)
					if !c.cli.IsRunning() {
						c.cli.Stop()
						c.cli, err = http.New(c.Url, "/websocket")
						if err != nil {
							c.logger.Error("reconnect fail", "err", err)
							continue
						}
					}
				}
				for _, res := range events.TxsResults {
					for _, event := range res.Events {
						c.handleEvent(ctx, event, c.Height)
					}
				}
				err = c.handleVote(ctx, c.Height)
				if err != nil {
					c.logger.Error("handleVote fail", "height", c.Height, "err", err)
					continue
				}
				if err := c.db.Save(Height{
					Id:     1,
					Height: uint64(c.Height),
				}).Error; err != nil {
					c.logger.Error("save height fail", "err", err)
					continue
				}
				c.Height++
			}
		}
	}
}

func (c *ChainIndexer) queryAccount(ctx context.Context, index uint64, address string) (*state.Account, error) {
	var err error
	var dat []byte
	if len(address) > 0 {
		dat, err = hex.DecodeString(address)
		if err != nil {
			return nil, err
		}
	} else {
		s := fmt.Sprintf("0%x", index)
		if len(s)&1 == 1 {
			s = s[1:]
		}
		dat, _ = hex.DecodeString(s)
	}
	res, err := c.cli.ABCIQuery(ctx, "/accounts/", dat)
	if err != nil {
		c.logger.Error("ABCIQuery fail", "err", err)
		if !c.cli.IsRunning() {
			c.cli.Stop()
			c.cli, err = http.New(c.Url, "/websocket")
			if err != nil {
				c.logger.Error("reconnect fail", "err", err)
				return nil, err
			}
		}
	}
	if res.Response.Code != 0 {
		fmt.Printf("%#v\n", res)
		return nil, errors.New("response code 0")
	}
	var act state.Account
	err = act.UnmarshalJSON(res.Response.Value)
	if err != nil {
		return nil, err
	}
	return &act, err
}
