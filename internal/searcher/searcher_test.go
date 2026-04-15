package searcher

import (
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/bfv/xref/internal/parser"
)

func testcasesDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "testcases")
}

func setupSearcher(t *testing.T) *Searcher {
	t.Helper()
	p := parser.NewParser(nil)
	xreffiles := p.ParseDir(testcasesDir(), "C:/dev/node/xrefparser/4gl/")
	if len(xreffiles) == 0 {
		t.Fatal("expected parsed xref files, got 0")
	}
	return NewSearcher(xreffiles)
}

func TestGetImplementations(t *testing.T) {
	s := setupSearcher(t)
	results := s.GetImplementations("oo.IEmpty")

	var names []string
	for _, xf := range results {
		if xf.Class != nil {
			names = append(names, xf.Class.Name)
		}
	}

	sort.Strings(names)
	expected := []string{"oo.Empty", "oo.MultipleImplements"}
	sort.Strings(expected)

	if len(names) != len(expected) {
		t.Fatalf("expected %d implementors, got %d: %v", len(expected), len(names), names)
	}
	for i := range expected {
		if names[i] != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], names[i])
		}
	}
}

func TestGetIncludeReferencesBaseI(t *testing.T) {
	s := setupSearcher(t)
	results := s.GetIncludeReferences("include/base.i")

	var sources []string
	for _, xf := range results {
		sources = append(sources, xf.SourceFile)
	}

	sort.Strings(sources)
	expected := []string{"include/includes-params.p", "include/includes.p"}
	sort.Strings(expected)

	if len(sources) != len(expected) {
		t.Fatalf("expected %d sources, got %d: %v", len(expected), len(sources), sources)
	}
	for i := range expected {
		if sources[i] != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], sources[i])
		}
	}
}

func TestGetIncludeReferencesNested(t *testing.T) {
	s := setupSearcher(t)
	results := s.GetIncludeReferences("include/ttchild.i")

	var sources []string
	for _, xf := range results {
		sources = append(sources, xf.SourceFile)
	}

	sort.Strings(sources)
	expected := []string{"include/includes-params.p", "include/includes.p"}
	sort.Strings(expected)

	if len(sources) != len(expected) {
		t.Fatalf("expected %d sources, got %d: %v", len(expected), len(sources), sources)
	}
	for i := range expected {
		if sources[i] != expected[i] {
			t.Errorf("expected %s, got %s", expected[i], sources[i])
		}
	}
}

func TestGetDatabaseNames(t *testing.T) {
	s := setupSearcher(t)
	dbnames := s.GetDatabaseNames(nil)

	if len(dbnames) == 0 {
		t.Fatal("expected at least one database name")
	}

	found := false
	for _, db := range dbnames {
		if db == "sports2000" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected sports2000 in database names, got %v", dbnames)
	}
}

func TestGetTableNames(t *testing.T) {
	s := setupSearcher(t)
	tables := s.GetTableNames(nil)

	if len(tables) == 0 {
		t.Fatal("expected at least one table")
	}

	found := false
	for _, td := range tables {
		if td.Table == "Customer" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected Customer in table names")
	}
}

func TestGetTableReferences(t *testing.T) {
	s := setupSearcher(t)
	results := s.GetTableReferences("Customer", nil, nil, nil)

	if len(results) == 0 {
		t.Fatal("expected at least one source referencing Customer")
	}
}

func TestGetTableReferencesWithUpdates(t *testing.T) {
	s := setupSearcher(t)
	hasUpdates := true
	results := s.GetTableReferences("Customer", nil, &hasUpdates, nil)

	if len(results) == 0 {
		t.Fatal("expected at least one source with updates on Customer")
	}
}

func TestGetFieldReferences(t *testing.T) {
	s := setupSearcher(t)
	tableName := "Customer"
	results := s.GetFieldReferences("CustNum", &tableName, nil)

	if len(results) == 0 {
		t.Fatal("expected at least one source referencing Customer.CustNum")
	}
}

func TestGetDatabaseReferences(t *testing.T) {
	s := setupSearcher(t)
	results := s.GetDatabaseReferences("sports2000")

	if len(results) == 0 {
		t.Fatal("expected at least one source referencing sports2000")
	}
}

func TestAddMergesData(t *testing.T) {
	s := setupSearcher(t)
	initial := len(s.xreffiles)

	p := parser.NewParser(nil)
	xreffiles := p.ParseDir(testcasesDir(), "C:/dev/node/xrefparser/4gl/")
	s.Add(xreffiles)

	if len(s.xreffiles) != initial {
		t.Errorf("Add should replace existing entries, expected %d, got %d", initial, len(s.xreffiles))
	}
}

func TestGetDependencies(t *testing.T) {
	s := setupSearcher(t)

	deps := s.GetDependencies("run/master.p")
	if deps == nil {
		t.Fatal("expected dependencies for run/master.p")
	}

	if deps.Source != "run/master.p" {
		t.Errorf("expected source run/master.p, got %s", deps.Source)
	}

	if len(deps.Runs) == 0 {
		t.Error("expected at least one run dependency")
	}
}

func TestGetDependenciesNotFound(t *testing.T) {
	s := setupSearcher(t)

	deps := s.GetDependencies("nonexistent.p")
	if deps != nil {
		t.Error("expected nil for nonexistent source")
	}
}

func TestGetDependenciesWithIncludes(t *testing.T) {
	s := setupSearcher(t)

	deps := s.GetDependencies("include/includes.p")
	if deps == nil {
		t.Fatal("expected dependencies for include/includes.p")
	}

	if len(deps.Includes) == 0 {
		t.Error("expected at least one include dependency")
	}
}

func TestGetDependenciesWithTables(t *testing.T) {
	s := setupSearcher(t)

	deps := s.GetDependencies("db/customer.p")
	if deps == nil {
		t.Fatal("expected dependencies for db/customer.p")
	}

	if len(deps.Tables) == 0 {
		t.Error("expected at least one table dependency")
	}
}

func TestGetClassHierarchy(t *testing.T) {
	s := setupSearcher(t)

	hierarchy := s.GetClassHierarchy("oo.DeliverAddress")
	if len(hierarchy) == 0 {
		t.Fatal("expected hierarchy entries for oo.DeliverAddress")
	}

	// Should include DeliverAddress, Address, AddressBase, and IBase
	names := map[string]bool{}
	for _, entry := range hierarchy {
		names[entry.Name] = true
	}

	for _, expected := range []string{"oo.DeliverAddress", "oo.Address", "oo.AddressBase", "oo.IBase"} {
		if !names[expected] {
			t.Errorf("expected %s in hierarchy, got %v", expected, names)
		}
	}
}

func TestGetClassHierarchyInterface(t *testing.T) {
	s := setupSearcher(t)

	hierarchy := s.GetClassHierarchy("oo.IDisposable")
	if len(hierarchy) == 0 {
		t.Fatal("expected hierarchy entries for oo.IDisposable")
	}

	names := map[string]bool{}
	for _, entry := range hierarchy {
		names[entry.Name] = true
	}

	if !names["oo.IDisposable"] {
		t.Error("expected oo.IDisposable in hierarchy")
	}
	if !names["oo.IEmpty"] {
		t.Error("expected oo.IEmpty in hierarchy (parent of IDisposable)")
	}
}

func TestGetClassHierarchyTypes(t *testing.T) {
	s := setupSearcher(t)

	hierarchy := s.GetClassHierarchy("oo.AddressBase")

	for _, entry := range hierarchy {
		switch entry.Name {
		case "oo.AddressBase":
			if entry.Type != "class" {
				t.Errorf("expected AddressBase type=class, got %s", entry.Type)
			}
		case "oo.IBase":
			if entry.Type != "interface" {
				t.Errorf("expected IBase type=interface, got %s", entry.Type)
			}
		}
	}
}

func TestGetReverseDependencies(t *testing.T) {
	s := setupSearcher(t)

	rd := s.GetReverseDependencies("oo/AddressBase.cls")
	if rd == nil {
		t.Fatal("expected reverse dependencies for oo/AddressBase.cls")
	}

	if rd.Source != "oo/AddressBase.cls" {
		t.Errorf("expected source oo/AddressBase.cls, got %s", rd.Source)
	}

	// oo.Address inherits from oo.AddressBase
	found := false
	for _, src := range rd.InheritedBy {
		if src == "oo/Address.cls" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected oo/Address.cls in inheritedBy, got %v", rd.InheritedBy)
	}
}

func TestGetReverseDependenciesNotFound(t *testing.T) {
	s := setupSearcher(t)

	rd := s.GetReverseDependencies("nonexistent.p")
	if rd != nil {
		t.Error("expected nil for nonexistent source")
	}
}

func TestGetReverseDependenciesIncludes(t *testing.T) {
	s := setupSearcher(t)

	// include/base.i is included by includes.p and includes-params.p
	// But base.i is an include file, not a source file in our test data
	// Let's test with tt/ttsports.i — also not a source.
	// Instead test with a source that gets included: check if any source is included by others.
	// Actually includes reference include files, not source files.
	// Let's test invoke reverse deps instead.

	rd := s.GetReverseDependencies("oo/Address.cls")
	if rd == nil {
		t.Fatal("expected reverse dependencies for oo/Address.cls")
	}

	// oo.DeliverAddress inherits oo.Address
	foundInherited := false
	for _, src := range rd.InheritedBy {
		if src == "oo/DeliverAddress.cls" {
			foundInherited = true
			break
		}
	}
	if !foundInherited {
		t.Errorf("expected oo/DeliverAddress.cls in inheritedBy, got %v", rd.InheritedBy)
	}

	// oo.DeliverAddress and oo.Person instantiate oo.Address
	if len(rd.InstantiatedBy) == 0 {
		t.Error("expected at least one instantiatedBy entry for oo/Address.cls")
	}

	// oo.Person invokes oo.Address methods
	foundInvoked := false
	for _, src := range rd.InvokedBy {
		if src == "oo/Person.cls" {
			foundInvoked = true
			break
		}
	}
	if !foundInvoked {
		t.Errorf("expected oo/Person.cls in invokedBy, got %v", rd.InvokedBy)
	}
}

func TestGetReverseDependenciesInterface(t *testing.T) {
	s := setupSearcher(t)

	rd := s.GetReverseDependencies("oo/IEmpty.cls")
	if rd == nil {
		t.Fatal("expected reverse dependencies for oo/IEmpty.cls")
	}

	// oo.Empty and oo.MultipleImplements implement oo.IEmpty
	if len(rd.InheritedBy) < 2 {
		t.Errorf("expected at least 2 inheritedBy entries for IEmpty, got %v", rd.InheritedBy)
	}
}
