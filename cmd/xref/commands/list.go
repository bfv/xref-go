package commands

import (
	"fmt"

	"github.com/bfv/xref/internal/datafile"
	"github.com/bfv/xref/internal/searcher"
	"github.com/spf13/cobra"
)

// NewListCmd returns the list subcommand.
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List databases, tables or sources",
		RunE: func(cmd *cobra.Command, args []string) error {
			input, _ := cmd.Flags().GetString("input")
			listDbs, _ := cmd.Flags().GetBool("databases")
			listTables, _ := cmd.Flags().GetBool("tables")
			listSources, _ := cmd.Flags().GetBool("sources")

			xrefdata, err := datafile.Load(input)
			if err != nil {
				return err
			}

			s := searcher.NewSearcher(xrefdata)

			if listDbs {
				dbnames := s.GetDatabaseNames(nil)
				fmt.Printf("Databases (%d):\n", len(dbnames))
				for _, db := range dbnames {
					fmt.Println(" ", db)
				}
				return nil
			}

			if listTables {
				tables := s.GetTableNames(nil)
				fmt.Printf("Tables (%d):\n", len(tables))
				for _, t := range tables {
					fmt.Printf("  %s.%s\n", t.Database, t.Table)
				}
				return nil
			}

			if listSources {
				sources := s.GetSourceNames()
				fmt.Printf("Sources (%d):\n", len(sources))
				for _, src := range sources {
					fmt.Println(" ", src)
				}
				return nil
			}

			return fmt.Errorf("specify --databases, --tables or --sources")
		},
	}

	cmd.Flags().StringP("input", "i", datafile.DefaultDataFile, "Input JSON data file")
	cmd.Flags().BoolP("databases", "d", false, "List database names")
	cmd.Flags().BoolP("tables", "t", false, "List table names")
	cmd.Flags().BoolP("sources", "s", false, "List source files")

	return cmd
}
