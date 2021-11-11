package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

const envVarPrefix = "AB_"

var RootCmd = &cobra.Command{
	Use:   "bot",
	Short: "Alert bot",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			optName := strings.ToUpper(f.Name)
			optName = strings.ReplaceAll(optName, "-", "_")
			varName := envVarPrefix + optName

			if val, ok := os.LookupEnv(varName); !f.Changed && ok {
				if err := f.Value.Set(val); err != nil {
					panic(fmt.Errorf("invalid environment variable %s: %w", varName, err))
				}
			}
		})
	},
}
