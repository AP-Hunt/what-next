package context

import (
	"context"

	"github.com/AP-Hunt/what-next/m/todo"
	"github.com/AP-Hunt/what-next/m/views"
)

type ContextKey = string

const (
	CtxTodoRepo   ContextKey = "TodoRepo"
	CtxViewEngine ContextKey = "ViewEngine"
)

type CommandContext struct {
	context.Context
}

func NewCommandContext(parentContext context.Context) CommandContext {
	return CommandContext{
		parentContext,
	}
}

func (ctx CommandContext) WithTodoRepository(repo todo.TodoRepositoryInterface) CommandContext {
	return CommandContext{context.WithValue(ctx, CtxTodoRepo, repo)}
}

func (ctx CommandContext) TodoRepository() todo.TodoRepositoryInterface {
	return ctx.Value(CtxTodoRepo).(todo.TodoRepositoryInterface)
}

func (ctx CommandContext) WithViewEngine(engine views.ViewEngineInterface) CommandContext {
	return CommandContext{context.WithValue(ctx, CtxViewEngine, engine)}
}

func (ctx CommandContext) ViewEngine() views.ViewEngineInterface {
	return ctx.Value(CtxViewEngine).(views.ViewEngineInterface)
}
