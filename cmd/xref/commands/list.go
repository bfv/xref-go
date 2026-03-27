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
		Short: "List databases or tables",
		RunE: func(cmd *cobra.Command, args []string) error {
			input, _ := cmd.Flags().GetString("input")
			listDbs, _ := cmd.Flags().GetBool("databases")
			listTables, _ := cmd.Flags().GetBool("tables")

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

			return fmt.Errorf("specify --databases or --tables")
		},
	}

	cmd.Flags().StringP("input", "i", datafile.DefaultDataFile, "Input JSON data file")
	cmd.Flags().Bool("databases", false, "List database names")
	cmd.Flags().Bool("tables", false, "List table names")

	return cmd
}
