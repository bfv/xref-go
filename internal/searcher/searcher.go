package searcher

import (
	"sort"
	"strings"

	"github.com/bfv/xref/internal/models"
)

// Searcher provides query methods over parsed XrefFile data.
type Searcher struct {
	xreffiles []*models.XrefFile
}

// NewSearcher creates a Searcher with the given xref data.
func NewSearcher(xreffiles []*models.XrefFile) *Searcher {
	return &Searcher{xreffiles: xreffiles}
}

// GetDatabaseNames returns unique database names found across all (or filtered) sources.
func (s *Searcher) GetDatabaseNames(sources []string) []string {
	seen := map[string]bool{}
	var dbnames []string

	for _, xf := range s.xreffiles {
		if sources != nil && !containsSource(sources, xf.SourceFile) {
			continue
		}
		for _, table := range xf.Tables {
			lower := strings.ToLower(table.Database)
			if !seen[lower] {
				seen[lower] = true
				dbnames = append(dbnames, lower)
			}
		}
	}

	sort.Slice(dbnames, func(i, j int) bool {
		return strings.ToLower(dbnames[i]) < strings.ToLower(dbnames[j])
	})
	return dbnames
}

// GetTableNames returns unique table definitions found across all (or filtered) sources.
func (s *Searcher) GetTableNames(sources []string) []models.TableDefinition {
	var tables []models.TableDefinition

	for _, xf := range s.xreffiles {
		if sources != nil && !containsSource(sources, xf.SourceFile) {
			continue
		}
		for _, table := range xf.Tables {
			db := strings.ToLower(table.Database)
			found := false
			for _, t := range tables {
				if t.Database == db && t.Table == table.Name {
					found = true
					break
				}
			}
			if !found {
				tables = append(tables, models.TableDefinition{
					Table:    table.Name,
					Database: db,
				})
			}
		}
	}

	sort.Slice(tables, func(i, j int) bool {
		ai := strings.ToLower(tables[i].Database + "." + tables[i].Table)
		bj := strings.ToLower(tables[j].Database + "." + tables[j].Table)
		return ai < bj
	})
	return tables
}

// GetSourceNames returns the sorted list of source file names across all xref files.
func (s *Searcher) GetSourceNames() []string {
	seen := map[string]bool{}
	var sources []string

	for _, xf := range s.xreffiles {
		if xf.SourceFile != "" && !seen[xf.SourceFile] {
			seen[xf.SourceFile] = true
			sources = append(sources, xf.SourceFile)
		}
	}

	sort.Slice(sources, func(i, j int) bool {
		return strings.ToLower(sources[i]) < strings.ToLower(sources[j])
	})
	return sources
}

// GetSourceByName returns the XrefFile for the given source file path, or nil if not found.
func (s *Searcher) GetSourceByName(sourcefile string) *models.XrefFile {
	for _, xf := range s.xreffiles {
		if strings.EqualFold(xf.SourceFile, sourcefile) {
			return xf
		}
	}
	return nil
}

// GetTableReferences returns XrefFiles that reference the given table with optional CRUD filters.
// Pass nil for a has* parameter to ignore that criterion.
func (s *Searcher) GetTableReferences(tablename string, hasCreates, hasUpdates, hasDeletes *bool) []*models.XrefFile {
	noCriteria := hasCreates == nil && hasUpdates == nil && hasDeletes == nil

	var result []*models.XrefFile
	for _, xf := range s.xreffiles {
		for _, table := range xf.Tables {
			if !strings.EqualFold(table.Name, tablename) {
				continue
			}
			if noCriteria ||
				(hasCreates != nil && table.IsCreated == *hasCreates) ||
				(hasUpdates != nil && table.IsUpdated == *hasUpdates) ||
				(hasDeletes != nil && table.IsDeleted == *hasDeletes) {
				result = append(result, xf)
				break
			}
		}
	}

	return result
}

// GetFieldReferences returns XrefFiles that reference the given field, optionally within a table.
func (s *Searcher) GetFieldReferences(fieldname string, tablename *string, hasUpdates *bool) []*models.XrefFile {
	var xreffiles []*models.XrefFile
	if tablename != nil {
		xreffiles = s.GetTableReferences(*tablename, nil, nil, nil)
	} else {
		xreffiles = s.xreffiles
	}

	noCriteria := hasUpdates == nil

	var result []*models.XrefFile
	for _, xf := range xreffiles {
		found := false
		for _, table := range xf.Tables {
			if tablename != nil && !strings.EqualFold(table.Name, *tablename) {
				continue
			}
			for _, field := range table.Fields {
				if strings.EqualFold(field.Name, fieldname) &&
					(noCriteria || field.IsUpdated == *hasUpdates) {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if found {
			result = append(result, xf)
		}
	}

	return result
}

// GetDatabaseReferences returns XrefFiles that reference the given database.
func (s *Searcher) GetDatabaseReferences(databaseName string) []*models.XrefFile {
	var result []*models.XrefFile

	for _, xf := range s.xreffiles {
		for _, table := range xf.Tables {
			if strings.EqualFold(table.Database, databaseName) {
				result = append(result, xf)
				break
			}
		}
	}

	return result
}

// GetImplementations returns XrefFiles whose class implements the given interface.
func (s *Searcher) GetImplementations(interfaceName string) []*models.XrefFile {
	var result []*models.XrefFile

	for _, xf := range s.xreffiles {
		if xf.Class != nil {
			for _, impl := range xf.Class.Implements {
				if impl == interfaceName {
					result = append(result, xf)
					break
				}
			}
		}
	}

	return result
}

// GetIncludeReferences returns XrefFiles that include the given file.
func (s *Searcher) GetIncludeReferences(includeName string) []*models.XrefFile {
	var result []*models.XrefFile

	for _, xf := range s.xreffiles {
		for _, inc := range xf.Includes {
			if inc == includeName {
				result = append(result, xf)
				break
			}
		}
	}

	return result
}

// Add merges xreffiles into the searcher, replacing existing entries by sourcefile.
func (s *Searcher) Add(xreffiles []*models.XrefFile) {
	for _, xf := range xreffiles {
		found := false
		for i, existing := range s.xreffiles {
			if existing.SourceFile == xf.SourceFile {
				s.xreffiles[i] = xf
				found = true
				break
			}
		}
		if !found {
			s.xreffiles = append(s.xreffiles, xf)
		}
	}
}

func containsSource(sources []string, source string) bool {
	for _, s := range sources {
		if s == source {
			return true
		}
	}
	return false
}
