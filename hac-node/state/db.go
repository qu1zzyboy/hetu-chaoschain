package state

import (
	"sync"

	cmtlog "github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/iavl"
	dbm "github.com/cosmos/iavl/db"
	"github.com/ethereum/go-ethereum/common"
)

type StateDB struct {
	mtx sync.RWMutex

	dir    string
	logger cmtlog.Logger
	db     *iavl.MutableTree

	state *State
}

func NewStateDB(dir string, logger cmtlog.Logger) (db *StateDB, err error) {
	logger = logger.With("module", "hacdb")
	ldb, err := dbm.NewDB("hac", "goleveldb", dir)
	if err != nil {
		return nil, err
	}
	tdb := iavl.NewMutableTree(ldb, 128, true, Cometbft2CosmosLogger(logger))
	version, err := tdb.Load()
	if err != nil {
		return nil, err
	}
	logger.Info("load db success", "version", version)
	st := newState(tdb, logger)
	err = st.load()
	if err != nil {
		logger.Error("from hacdb load fail", "err", err)
		return nil, err
	}
	db = &StateDB{
		dir:    dir,
		logger: logger,
		db:     tdb,
		state:  st,
	}
	return
}

func (db *StateDB) Close() (err error) {
	err = db.db.Close()
	return
}

func (db *StateDB) Header() (header *StateHeader) {
	db.mtx.RLock()
	defer db.mtx.RUnlock()
	header = db.state.Header()
	return
}

func (db *StateDB) State() *State {
	db.mtx.RLock()
	defer db.mtx.RUnlock()
	return db.state
}

func (db *StateDB) NewState() (st *State) {
	db.mtx.RLock()
	defer db.mtx.RUnlock()
	st = db.state.nextState()
	return
}

func (db *StateDB) SetState(st *State) (hash common.Hash, err error) {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	hash, err = st.save()
	if err != nil {
		return
	}
	db.state = st
	return
}

func (db *StateDB) GetAccountByIndex(idx uint64) (acnt *Account, height uint64, err error) {
	db.mtx.RLock()
	defer db.mtx.RUnlock()
	acnt, err = db.state.GetAccount(idx)
	if err != nil {
		return
	}
	if acnt != nil {
		acnt = acnt.Clone()
	}
	height = db.state.header.Height

	return

}

func (db *StateDB) GetAccountByAddress(addr []byte) (acnt *Account, height uint64, err error) {
	db.mtx.RLock()
	defer db.mtx.RUnlock()
	acnt, err = db.state.FindAccount(addr)
	if err != nil {
		return
	}
	if acnt != nil {
		acnt = acnt.Clone()
	}
	height = db.state.header.Height

	return
}
