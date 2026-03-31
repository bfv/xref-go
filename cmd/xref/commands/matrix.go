package commands

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/bfv/xref/internal/datafile"
	"github.com/bfv/xref/internal/models"
	"github.com/spf13/cobra"
)

// NewMatrixCmd returns the matrix subcommand.
func NewMatrixCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "matrix",
		Short: "Show a source/table matrix",
		RunE: func(cmd *cobra.Command, args []string) error {
			input, _ := cmd.Flags().GetString("input")
			dbFilter, _ := cmd.Flags().GetString("database")
			tablesFlag, _ := cmd.Flags().GetString("tables")
			tablesFileFlag, _ := cmd.Flags().GetString("tablesfile")
			noReads, _ := cmd.Flags().GetBool("noreads")

			xrefdata, err := datafile.Load(input)
			if err != nil {
				return err
			}

			// Build optional table filter
			tableFilter, err := buildTableFilter(tablesFlag, tablesFileFlag)
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
					if tableFilter != nil && !tableFilter[strings.ToLower(t.Name)] {
						continue
					}
					if noReads && !t.IsCreated && !t.IsUpdated && !t.IsDeleted {
						continue
					}
					tableSet[t.Name] = true
				}
			}

			var tableNames []string
			for t := range tableSet {
				tableNames = append(tableNames, t)
			}
			sort.Slice(tableNames, func(i, j int) bool {
				return strings.ToLower(tableNames[i]) < strings.ToLower(tableNames[j])
			})

			// Sort sources by name
			sort.Slice(xrefdata, func(i, j int) bool {
				return strings.ToLower(xrefdata[i].SourceFile) < strings.ToLower(xrefdata[j].SourceFile)
			})

			// Print header
			fmt.Printf("%-40s", "Source")
			for _, t := range tableNames {
				fmt.Printf(" %-12s", t)
			}
			fmt.Println()

			// Print rows
			for _, xf := range xrefdata {
				if !hasRelevantTables(xf, tableNames, dbFilter, noReads) {
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

	cmd.Flags().StringP("input", "i", datafile.DefaultDataFile, "Input JSON data file")
	cmd.Flags().StringP("database", "d", "", "Filter by database name")
	cmd.Flags().StringP("tables", "t", "", "Comma-separated list of table names to include")
	cmd.Flags().StringP("tablesfile", "f", "", "File with table names (one or more per line, comma-separated)")
	cmd.Flags().BoolP("noreads", "n", false, "Only show tables that have creates, updates or deletes")

	return cmd
}

func hasRelevantTables(xf *models.XrefFile, tableNames []string, dbFilter string, noReads bool) bool {
	for _, t := range xf.Tables {
		if dbFilter != "" && !strings.EqualFold(t.Database, dbFilter) {
			continue
		}
		if noReads && !t.IsCreated && !t.IsUpdated && !t.IsDeleted {
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

// buildTableFilter returns a set of lowercase table names to include, or nil if no filter is specified.
func buildTableFilter(tablesFlag, tablesFileFlag string) (map[string]bool, error) {
	var names []string

	if tablesFlag != "" {
		for _, t := range strings.Split(tablesFlag, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				names = append(names, strings.ToLower(t))
			}
		}
	}

	if tablesFileFlag != "" {
		f, err := os.Open(tablesFileFlag)
		if err != nil {
			return nil, fmt.Errorf("cannot open tables file: %w", err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			for _, t := range strings.Split(line, ",") {
				t = strings.TrimSpace(t)
				if t != "" {
					names = append(names, strings.ToLower(t))
				}
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading tables file: %w", err)
		}
	}

	if len(names) == 0 {
		return nil, nil
	}

	filter := map[string]bool{}
	for _, n := range names {
		filter[n] = true
	}
	return filter, nil
}
