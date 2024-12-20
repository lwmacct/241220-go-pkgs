package cmd

import (
	"reflect"

	"github.com/spf13/cobra"
)

type Ts struct {
	cc    *cobra.Command
	flags any
}

func New(flags any) *Ts {
	ref := &Ts{
		cc: &cobra.Command{
			CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
		},
	}

	ref.flags = flags
	return ref
}

func (t *Ts) CmdGet() *cobra.Command {

	return t.cc
}

func (t *Ts) CmdSet(name, short string) *cobra.Command {
	t.cc.Use = name
	t.cc.Short = short
	return t.cc
}

func (t *Ts) CmdSetPkgName() *cobra.Command {
	t.CmdSet(getPackageName(2), "")
	return t.cc
}

func (t *Ts) CmdAddCommand(cc *cobra.Command) {
	t.cc.AddCommand(cc)
}

func (t *Ts) AddCmd(
	runFunc func(cmd *cobra.Command, args []string),
	name, short string,
	group ...string,
) {
	cmd := &cobra.Command{
		Use:   name,
		Short: short,
		Run:   runFunc,
	}
	if len(group) > 0 {
		bind(cmd, reflect.ValueOf(t.flags).Elem(), "", group)
	}

	t.cc.AddCommand(cmd)

}

func (c *Ts) Execute() error {
	return c.cc.Execute()
}
