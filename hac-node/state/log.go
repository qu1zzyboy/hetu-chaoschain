package state

import (
	cosmoslog "cosmossdk.io/log"
	cmtlog "github.com/cometbft/cometbft/libs/log"
)

type Logger struct {
	logger cmtlog.Logger
}

func Cometbft2CosmosLogger(lg cmtlog.Logger) cosmoslog.Logger {
	return Logger{logger: lg}
}

func (l Logger) Info(msg string, keyVals ...any) {
	l.logger.Info(msg, keyVals...)
}

func (l Logger) Error(msg string, keyVals ...any) {
	l.logger.Error(msg, keyVals...)
}

func (l Logger) Debug(msg string, keyVals ...any) {
	l.logger.Debug(msg, keyVals...)
}

func (l Logger) With(keyVals ...any) cosmoslog.Logger {
	return Logger{l.logger.With(keyVals...)}
}

func (l Logger) Impl() any {
	return l.logger
}
