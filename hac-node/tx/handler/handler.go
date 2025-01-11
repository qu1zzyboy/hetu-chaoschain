package handler

import (
	"context"

	"github.com/calehh/hac-app/state"
	"github.com/calehh/hac-app/tx"
	abcitypes "github.com/cometbft/cometbft/abci/types"
)

type TxHandler interface {
	Check(ctx context.Context, st *state.State, btx *tx.HACTx) (res *abcitypes.ResponseCheckTx, err error)
	NewContext(ctx context.Context)
	Prepare(ctx context.Context, st *state.State, btx *tx.HACTx, code tx.VoteCode) (res *abcitypes.ExecTxResult, err error)
	Process(ctx context.Context, st *state.State, btx *tx.HACTx, code tx.VoteCode) (res *abcitypes.ExecTxResult, err error)
}
