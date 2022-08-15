package common

import (
	"context"
	"fmt"
	"reflect"
)

type Command interface{}

type CommandBus interface {
	Execute(command Command) error
	Register(command Command, commandHandler CommandHandler)
}

type CommandHandler interface {
	Handle(ctx context.Context, command Command) error
}

type MemoryCommandBus struct {
	commandHandlers map[string]CommandHandler
}

func NewMemoryCommandBus() MemoryCommandBus {
	return MemoryCommandBus{
		commandHandlers: make(map[string]CommandHandler),
	}
}

func (commandBus *MemoryCommandBus) Register(command Command, handler CommandHandler) error {
	name := reflect.TypeOf(command).String()
	_, alreadyRegistered := commandBus.commandHandlers[name]
	if alreadyRegistered {
		return ErrHandlerAlreadyRegistered{command: name, handler: reflect.TypeOf(handler).String()}
	}
	commandBus.commandHandlers[name] = handler
	return nil
}

type ErrHandlerAlreadyRegistered struct {
	command string
	handler string
}

func (err ErrHandlerAlreadyRegistered) Error() string {
	return fmt.Sprintf("handler '%s' already registered for command '%s'", err.handler, err.command)
}
