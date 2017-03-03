package cli

import (
	"bytes"
	"fmt"
)

type ErrMissingRequired struct {
	flags []string
}

func (e *ErrMissingRequired) Error() string {
	buf := &bytes.Buffer{}
	prefixed := func (n string) string {
		return prefixFor(n) + n
	}

	if len(e.flags) > 1 {
		fmt.Fprintf(buf, "Flags ")
		for i, name := range e.flags {
			if i > 0 {
				fmt.Fprintf(buf, ", ")
			}
			if i == len(e.flags)-1 {
				fmt.Fprintf(buf, "and ")
			}
			fmt.Fprintf(buf, "`%s`", prefixed(name))
		}
		fmt.Fprintf(buf, " are required")
	} else {
		fmt.Fprintf(buf, "Flag `%s` is required",
			prefixed(e.flags[0]))
	}
	return buf.String()
}

func checkCommandMissingRequiredFlags(c *Context, cmd *Command) error {
	var oops []string

	for _, f := range c.App.Flags {
		if req, ok := f.(IsRequirable); ok && req.IsRequired() {
			flagName := f.Names()[0]
			if !c.IsSet(flagName) {
				oops = append(oops, flagName)
			}
		}
	}
	if cmd != nil {
	for _, f := range cmd.Flags {
		if req, ok := f.(IsRequirable); ok && req.IsRequired() {
			flagName := f.Names()[0]
			if !c.IsSet(flagName) {
				oops = append(oops, flagName)
			}
		}
	}
	}

	if oops == nil {
		return nil
	}
	return &ErrMissingRequired{oops}
}
