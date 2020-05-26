package cmd

import (
	"errors"
	"fmt"
	"os"
)

type Runner interface {
	Init([]string) error
	ValidateFlags() error
	Run() (err error)
	Name() string
}

// Execute command
func Execute(args []string) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub-command")
	}

	cmds := []Runner{
		NewAdjustCmd(),
		NewShowCmd(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			if err := cmd.Init(os.Args[2:]); err != nil {
				return err
			}
			if err := cmd.ValidateFlags(); err != nil {
				return err
			}

			return cmd.Run()
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}
