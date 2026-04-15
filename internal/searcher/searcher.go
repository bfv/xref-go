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

// Dependencies holds the aggregated dependency data for a single source.
type Dependencies struct {
	Source       string                    `json:"source"`
	Tables       []models.Table            `json:"tables"`
	TempTables   []models.TempTable        `json:"temptables"`
	Includes     []string                  `json:"includes"`
	Runs         []models.Run              `json:"runs"`
	Instantiates []string                  `json:"instantiates"`
	Invokes      []models.MethodInvocation `json:"invokes"`
	Class        *models.Class             `json:"class,omitempty"`
	Interface    *models.Interface         `json:"interface,omitempty"`
}

// GetDependencies aggregates tables, includes, runs, instantiates, and invokes for a source.
func (s *Searcher) GetDependencies(source string) *Dependencies {
	xf := s.GetSourceByName(source)
	if xf == nil {
		return nil
	}
	return &Dependencies{
		Source:       xf.SourceFile,
		Tables:       xf.Tables,
		TempTables:   xf.TempTables,
		Includes:     xf.Includes,
		Runs:         xf.Runs,
		Instantiates: xf.Instantiates,
		Invokes:      xf.Invokes,
		Class:        xf.Class,
		Interface:    xf.Interface,
	}
}

// ClassHierarchyEntry represents a single class/interface in the hierarchy.
type ClassHierarchyEntry struct {
	Name       string   `json:"name"`
	Source     string   `json:"source,omitempty"`
	Inherits   []string `json:"inherits,omitempty"`
	Implements []string `json:"implements,omitempty"`
	Type       string   `json:"type"` // "class" or "interface"
}

// GetClassHierarchy resolves the full inheritance chain for a given class or interface name.
// It walks up the inheritance tree, collecting all ancestors.
func (s *Searcher) GetClassHierarchy(name string) []ClassHierarchyEntry {
	visited := map[string]bool{}
	var result []ClassHierarchyEntry
	s.collectHierarchy(name, visited, &result)
	return result
}

func (s *Searcher) collectHierarchy(name string, visited map[string]bool, result *[]ClassHierarchyEntry) {
	lower := strings.ToLower(name)
	if visited[lower] {
		return
	}
	visited[lower] = true

	for _, xf := range s.xreffiles {
		if xf.Class != nil && strings.EqualFold(xf.Class.Name, name) {
			entry := ClassHierarchyEntry{
				Name:       xf.Class.Name,
				Source:     xf.SourceFile,
				Inherits:   xf.Class.Inherits,
				Implements: xf.Class.Implements,
				Type:       "class",
			}
			*result = append(*result, entry)
			for _, parent := range xf.Class.Inherits {
				s.collectHierarchy(parent, visited, result)
			}
			for _, iface := range xf.Class.Implements {
				s.collectHierarchy(iface, visited, result)
			}
			return
		}
		if xf.Interface != nil && strings.EqualFold(xf.Interface.Name, name) {
			entry := ClassHierarchyEntry{
				Name:     xf.Interface.Name,
				Source:   xf.SourceFile,
				Inherits: xf.Interface.Inherits,
				Type:     "interface",
			}
			*result = append(*result, entry)
			for _, parent := range xf.Interface.Inherits {
				s.collectHierarchy(parent, visited, result)
			}
			return
		}
	}

	// Not found in our data — record it as an external/unknown entry
	*result = append(*result, ClassHierarchyEntry{
		Name: name,
		Type: "class",
	})
}

// ReverseDependencies holds sources that reference a given source.
type ReverseDependencies struct {
	Source         string   `json:"source"`
	IncludedBy     []string `json:"includedBy"`
	RunBy          []string `json:"runBy"`
	InheritedBy    []string `json:"inheritedBy"`
	InvokedBy      []string `json:"invokedBy"`
	InstantiatedBy []string `json:"instantiatedBy"`
}

// GetReverseDependencies finds sources that reference the given source via includes, RUN, inheritance, invokes, or instantiation.
func (s *Searcher) GetReverseDependencies(source string) *ReverseDependencies {
	xf := s.GetSourceByName(source)
	if xf == nil {
		return nil
	}

	rd := &ReverseDependencies{
		Source:         xf.SourceFile,
		IncludedBy:     []string{},
		RunBy:          []string{},
		InheritedBy:    []string{},
		InvokedBy:      []string{},
		InstantiatedBy: []string{},
	}

	// Determine class/interface name for this source
	var className string
	if xf.Class != nil {
		className = xf.Class.Name
	} else if xf.Interface != nil {
		className = xf.Interface.Name
	}

	for _, other := range s.xreffiles {
		if strings.EqualFold(other.SourceFile, xf.SourceFile) {
			continue
		}

		// Check includes
		for _, inc := range other.Includes {
			if strings.EqualFold(inc, xf.SourceFile) {
				rd.IncludedBy = append(rd.IncludedBy, other.SourceFile)
				break
			}
		}

		// Check RUN references
		for _, run := range other.Runs {
			if strings.EqualFold(run.Name, xf.SourceFile) {
				rd.RunBy = append(rd.RunBy, other.SourceFile)
				break
			}
		}

		if className != "" {
			// Check inheritance (class inherits or interface inherits)
			if other.Class != nil {
				for _, parent := range other.Class.Inherits {
					if strings.EqualFold(parent, className) {
						rd.InheritedBy = append(rd.InheritedBy, other.SourceFile)
						break
					}
				}
				for _, iface := range other.Class.Implements {
					if strings.EqualFold(iface, className) {
						rd.InheritedBy = append(rd.InheritedBy, other.SourceFile)
						break
					}
				}
			}
			if other.Interface != nil {
				for _, parent := range other.Interface.Inherits {
					if strings.EqualFold(parent, className) {
						rd.InheritedBy = append(rd.InheritedBy, other.SourceFile)
						break
					}
				}
			}

			// Check invokes
			for _, inv := range other.Invokes {
				if strings.EqualFold(inv.Class, className) {
					rd.InvokedBy = append(rd.InvokedBy, other.SourceFile)
					break
				}
			}

			// Check instantiates
			for _, inst := range other.Instantiates {
				if strings.EqualFold(inst, className) {
					rd.InstantiatedBy = append(rd.InstantiatedBy, other.SourceFile)
					break
				}
			}
		}
	}

	sort.Strings(rd.IncludedBy)
	sort.Strings(rd.RunBy)
	sort.Strings(rd.InheritedBy)
	sort.Strings(rd.InvokedBy)
	sort.Strings(rd.InstantiatedBy)

	return rd
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
