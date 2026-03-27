package commands

import (
	"fmt"

	"github.com/bfv/xref/internal/datafile"
	"github.com/bfv/xref/internal/logging"
	"github.com/bfv/xref/internal/searcher"
	"github.com/spf13/cobra"
)

// NewSearchCmd returns the search subcommand.
func NewSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for table, field, or database references",
		RunE: func(cmd *cobra.Command, args []string) error {
			input, _ := cmd.Flags().GetString("input")
			tableName, _ := cmd.Flags().GetString("table")
			fieldName, _ := cmd.Flags().GetString("field")
			dbName, _ := cmd.Flags().GetString("database")
			hasCreates, _ := cmd.Flags().GetBool("creates")
			hasUpdates, _ := cmd.Flags().GetBool("updates")
			hasDeletes, _ := cmd.Flags().GetBool("deletes")
			createsSet := cmd.Flags().Changed("creates")
			updatesSet := cmd.Flags().Changed("updates")
			deletesSet := cmd.Flags().Changed("deletes")

			xrefdata, err := datafile.Load(input)
			if err != nil {
				return err
			}

			logging.Logger.Debug().Int("files", len(xrefdata)).Msg("loaded repo data")

			s := searcher.NewSearcher(xrefdata)

			if dbName != "" {
				refs := s.GetDatabaseReferences(dbName)
				fmt.Printf("Sources referencing database '%s': %d\n", dbName, len(refs))
				for _, xf := range refs {
					fmt.Println(" ", xf.SourceFile)
				}
				return nil
			}

			if tableName != "" && fieldName != "" {
				var hasUpdatesPtr *bool
				if updatesSet {
					hasUpdatesPtr = &hasUpdates
				}
				refs := s.GetFieldReferences(fieldName, &tableName, hasUpdatesPtr)
				fmt.Printf("Sources referencing field '%s' in table '%s': %d\n", fieldName, tableName, len(refs))
				for _, xf := range refs {
					fmt.Println(" ", xf.SourceFile)
				}
				return nil
			}

			if fieldName != "" {
				var hasUpdatesPtr *bool
				if updatesSet {
					hasUpdatesPtr = &hasUpdates
				}
				refs := s.GetFieldReferences(fieldName, nil, hasUpdatesPtr)
				fmt.Printf("Sources referencing field '%s': %d\n", fieldName, len(refs))
				for _, xf := range refs {
					fmt.Println(" ", xf.SourceFile)
				}
				return nil
			}

			if tableName != "" {
				var createsPtr, updatesPtr, deletesPtr *bool
				if createsSet {
					createsPtr = &hasCreates
				}
				if updatesSet {
					updatesPtr = &hasUpdates
				}
				if deletesSet {
					deletesPtr = &hasDeletes
				}
				refs := s.GetTableReferences(tableName, createsPtr, updatesPtr, deletesPtr)
				fmt.Printf("Sources referencing table '%s': %d\n", tableName, len(refs))
				for _, xf := range refs {
					fmt.Println(" ", xf.SourceFile)
				}
				return nil
			}

			return fmt.Errorf("specify at least --table, --field, or --database")
		},
	}

	cmd.Flags().StringP("input", "i", datafile.DefaultDataFile, "Input JSON data file")
	cmd.Flags().StringP("table", "t", "", "Table name to search for")
	cmd.Flags().StringP("field", "f", "", "Field name to search for")
	cmd.Flags().StringP("database", "d", "", "Database name to search for")
	cmd.Flags().Bool("creates", false, "Filter on creates")
	cmd.Flags().Bool("updates", false, "Filter on updates")
	cmd.Flags().Bool("deletes", false, "Filter on deletes")

	return cmd
}
