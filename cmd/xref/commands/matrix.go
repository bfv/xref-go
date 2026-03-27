package commands

import (
	"fmt"
	"strings"

	"github.com/bfv/xref/internal/config"
	"github.com/bfv/xref/internal/models"
	"github.com/spf13/cobra"
)

// NewMatrixCmd returns the matrix subcommand.
func NewMatrixCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "matrix",
		Short: "Show a source/table matrix",
		RunE: func(cmd *cobra.Command, args []string) error {
			repoName, _ := cmd.Flags().GetString("name")
			dbFilter, _ := cmd.Flags().GetString("database")

			cfg, err := config.NewConfig()
			if err != nil {
				return err
			}

			if repoName == "" {
				repoName = cfg.Data.Current
			}
			if repoName == "" {
				return fmt.Errorf("no repo specified and no current repo set")
			}

			xrefdata, err := cfg.ReadRepoData(repoName)
			if err != nil {
				return err
			}

			// Collect unique table names
			tableSet := map[string]bool{}
			for _, xf := range xrefdata {
				for _, t := range xf.Tables {
					if dbFilter != "" && !strings.EqualFold(t.Database, dbFilter) {
						continue
					}
					tableSet[t.Name] = true
				}
			}

			var tableNames []string
			for t := range tableSet {
				tableNames = append(tableNames, t)
			}

			// Print header
			fmt.Printf("%-40s", "Source")
			for _, t := range tableNames {
				fmt.Printf(" %-12s", t)
			}
			fmt.Println()

			// Print rows
			for _, xf := range xrefdata {
				if !hasRelevantTables(xf, tableNames, dbFilter) {
					continue
				}
				fmt.Printf("%-40s", truncate(xf.SourceFile, 39))
				for _, tn := range tableNames {
					flag := tableFlag(xf, tn, dbFilter)
					fmt.Printf(" %-12s", flag)
				}
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Repository name (defaults to current)")
	cmd.Flags().StringP("database", "d", "", "Filter by database name")

	return cmd
}

func hasRelevantTables(xf *models.XrefFile, tableNames []string, dbFilter string) bool {
	for _, t := range xf.Tables {
		if dbFilter != "" && !strings.EqualFold(t.Database, dbFilter) {
			continue
		}
		for _, tn := range tableNames {
			if t.Name == tn {
				return true
			}
		}
	}
	return false
}

func tableFlag(xf *models.XrefFile, tableName, dbFilter string) string {
	for _, t := range xf.Tables {
		if t.Name != tableName {
			continue
		}
		if dbFilter != "" && !strings.EqualFold(t.Database, dbFilter) {
			continue
		}
		flags := ""
		if t.IsCreated {
			flags += "C"
		}
		if t.IsUpdated {
			flags += "U"
		}
		if t.IsDeleted {
			flags += "D"
		}
		if flags == "" {
			flags = "R"
		}
		return flags
	}
	return "-"
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
