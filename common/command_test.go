//go:build spec

package common

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestHandler struct{}

type TestCommand struct {
	test string
}

func (handler *TestHandler) Handle(ctx context.Context, command *TestCommand) error {
	return fmt.Errorf(command.test)
}

func TestMemoryCommandBus(t *testing.T) {
	assert := assert.New(t)
	commandBus := NewMemoryCommandBus()

	t.Run("it should register a handler", func(t *testing.T) {
		err := commandBus.Register(&TestCommand{}, &TestHandler{})

		assert.NoError(err)
	})

	t.Run("it should not register a handler twice", func(t *testing.T) {
		commandBus.Register(&TestCommand{}, &TestHandler{})

		err = commandBus.Register(&TestCommand{}, &TestHandler{})

		assert.Error(err)
	})

	t.Run("it should execute a command", func(t *testing.T) {
		commandBus.Register(&TestCommand{}, &TestHandler{})

		err := commandBus.Execute(&TestCommand{"yolala"})

		assert.ErrorContains(err, "yolala")
	})
}
