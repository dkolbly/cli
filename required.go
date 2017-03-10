package cli

import (
	"os"
	"bytes"
	"fmt"
)

type ErrMissingRequired struct {
	flags []string
}

func (e *ErrMissingRequired) Error() string {
	buf := &bytes.Buffer{}
	prefixed := func(n string) string {
		return prefixFor(n) + n
	}

	if len(e.flags) > 1 {
		fmt.Fprintf(buf, "Flags ")
		for i, name := range e.flags {
			if i > 0 {
				if len(e.flags) > 2 {
					fmt.Fprintf(buf, ", ")
				} else {
					fmt.Fprintf(buf, " ")
				}
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

func isEnvVarSet(envVars []string) bool {
	for _, envVar := range envVars {
		if envVal := os.Getenv(envVar); envVal != "" {
			// TODO: Can't use this for bools as
			// set means that it was true or false based on
			// Bool flag type, should work for other types
			if len(envVal) > 0 {
				return true
			}
		}
	}

	return false
}

func checkCommandMissingRequiredFlags(c *Context, cmd *Command) error {
	var oops []string

	for _, f := range c.App.Flags {
		if req, ok := f.(IsRequirable); ok && req.IsRequired() {
			if !req.IsSetIn(c) {
				oops = append(oops, f.Names()[0])
			}
		}
	}
	if cmd != nil {
		for _, f := range cmd.Flags {
			if req, ok := f.(IsRequirable); ok && req.IsRequired() {
				if !req.IsSetIn(c) {
					oops = append(oops, f.Names()[0])
				}
			}
		}
	}

	if oops == nil {
		return nil
	}
	return &ErrMissingRequired{oops}
}
