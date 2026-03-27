package commands

import (
	"fmt"

	"github.com/bfv/xref/internal/config"
	"github.com/spf13/cobra"
)

// NewSwitchCmd returns the switch subcommand.
func NewSwitchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch",
		Short: "Switch the current active repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")

			if name == "" && len(args) > 0 {
				name = args[0]
			}
			if name == "" {
				return fmt.Errorf("specify a repo name with --name or as argument")
			}

			cfg, err := config.NewConfig()
			if err != nil {
				return err
			}

			if err := cfg.SetCurrent(name); err != nil {
				return err
			}

			if err := cfg.Save(); err != nil {
				return err
			}

			fmt.Printf("Switched to repo '%s'\n", name)
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Repository name to switch to")

	return cmd
}
