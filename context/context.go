package context

import (
	"context"

	"github.com/AP-Hunt/what-next/m/todo"
)

type ContextKey = string

const (
	CtxTodoRepo ContextKey = "TodoRepo"
)

type CommandContext struct {
	context.Context
}

func NewCommandContext(parentContext context.Context) CommandContext {
	return CommandContext{
		parentContext,
	}
}

func (ctx *CommandContext) WithTodoRepository(repo todo.TodoRepositoryInterface) CommandContext {
	return CommandContext{context.WithValue(ctx, CtxTodoRepo, repo)}
}

func (ctx *CommandContext) TodoRepository() todo.TodoRepositoryInterface {
	return ctx.Value(CtxTodoRepo).(todo.TodoRepositoryInterface)
}
