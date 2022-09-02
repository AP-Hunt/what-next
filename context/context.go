package context

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ContextKey = string

const (
	CtxDatabaseConnection ContextKey = "DatabaseConnection"
)

type CommandContext struct {
	context.Context
}

func NewCommandContext(parentContext context.Context) CommandContext {
	return CommandContext{
		parentContext,
	}
}

func (ctx *CommandContext) WithDatabaseConnection(conn *sqlx.DB) CommandContext {
	return CommandContext{context.WithValue(ctx, CtxDatabaseConnection, conn)}
}

func (ctx *CommandContext) DatabaseConnection() *sqlx.DB {
	return ctx.Value(CtxDatabaseConnection).(*sqlx.DB)
}
