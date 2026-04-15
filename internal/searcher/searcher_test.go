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

func TestGetMigrationScope(t *testing.T) {
	s := setupSearcher(t)

	scope := s.GetMigrationScope("db/customer.p")
	if scope == nil {
		t.Fatal("expected migration scope for db/customer.p")
	}

	if scope.StartSource != "db/customer.p" {
		t.Errorf("expected startSource db/customer.p, got %s", scope.StartSource)
	}

	// Should include at least the starting source
	if len(scope.Sources) == 0 {
		t.Fatal("expected at least one source in migration scope")
	}

	// Should include Customer table
	if len(scope.Tables) == 0 {
		t.Fatal("expected at least one table in migration scope")
	}

	foundCustomerTable := false
	for _, table := range scope.Tables {
		if table == "sports2000.customer" {
			foundCustomerTable = true
			break
		}
	}
	if !foundCustomerTable {
		t.Errorf("expected sports2000.customer in tables, got %v", scope.Tables)
	}
}

func TestGetMigrationScopeNotFound(t *testing.T) {
	s := setupSearcher(t)

	scope := s.GetMigrationScope("nonexistent.p")
	if scope != nil {
		t.Error("expected nil for nonexistent source")
	}
}

func TestGetMigrationScopeSharedTables(t *testing.T) {
	s := setupSearcher(t)

	// db/customer.p and db/custorderline.p both reference Customer
	// They should appear in each other's migration scope
	scope := s.GetMigrationScope("db/customer.p")
	if scope == nil {
		t.Fatal("expected migration scope")
	}

	foundCustOrderLine := false
	for _, src := range scope.Sources {
		if src == "db/custorderline.p" {
			foundCustOrderLine = true
			break
		}
	}
	if !foundCustOrderLine {
		t.Errorf("expected db/custorderline.p in scope (shares Customer table), got %v", scope.Sources)
	}
}

func TestGetMigrationScopeClassHierarchy(t *testing.T) {
	s := setupSearcher(t)

	// oo/Address.cls inherits from oo/AddressBase.cls, so they should be in scope together
	scope := s.GetMigrationScope("oo/Address.cls")
	if scope == nil {
		t.Fatal("expected migration scope for oo/Address.cls")
	}

	foundBase := false
	foundDeliver := false
	for _, src := range scope.Sources {
		if src == "oo/AddressBase.cls" {
			foundBase = true
		}
		if src == "oo/DeliverAddress.cls" {
			foundDeliver = true
		}
	}
	if !foundBase {
		t.Errorf("expected oo/AddressBase.cls in scope (parent class), got %v", scope.Sources)
	}
	if !foundDeliver {
		t.Errorf("expected oo/DeliverAddress.cls in scope (child class), got %v", scope.Sources)
	}
}

func TestGetCrudMatrix(t *testing.T) {
	s := setupSearcher(t)

	matrix := s.GetCrudMatrix(nil)

	if len(matrix.Sources) == 0 {
		t.Fatal("expected at least one source in CRUD matrix")
	}

	if len(matrix.Tables) == 0 {
		t.Fatal("expected at least one table in CRUD matrix")
	}

	if len(matrix.Entries) == 0 {
		t.Fatal("expected at least one entry in CRUD matrix")
	}
}

func TestGetCrudMatrixFiltered(t *testing.T) {
	s := setupSearcher(t)

	matrix := s.GetCrudMatrix([]string{"db/customer.p"})

	if len(matrix.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d: %v", len(matrix.Sources), matrix.Sources)
	}

	if matrix.Sources[0] != "db/customer.p" {
		t.Errorf("expected db/customer.p, got %s", matrix.Sources[0])
	}

	// Should have Customer table entry
	foundCustomer := false
	for _, entry := range matrix.Entries {
		if entry.Table == "sports2000.Customer" {
			foundCustomer = true
			if !entry.Reads {
				t.Error("expected reads=true for Customer")
			}
			break
		}
	}
	if !foundCustomer {
		t.Error("expected Customer table in CRUD matrix entries")
	}
}

func TestGetCrudMatrixCrudFlags(t *testing.T) {
	s := setupSearcher(t)

	// create-customer.p should have creates=true for Customer
	matrix := s.GetCrudMatrix([]string{"db/create-customer.p"})

	for _, entry := range matrix.Entries {
		if entry.Table == "sports2000.Customer" {
			if !entry.Creates {
				t.Error("expected creates=true for create-customer.p on Customer")
			}
			return
		}
	}
	t.Error("expected Customer table entry for create-customer.p")
}

func TestGetCrudMatrixDeleteFlags(t *testing.T) {
	s := setupSearcher(t)

	// delete-customer.p should have deletes=true for Customer
	matrix := s.GetCrudMatrix([]string{"db/delete-customer.p"})

	for _, entry := range matrix.Entries {
		if entry.Table == "sports2000.Customer" {
			if !entry.Deletes {
				t.Error("expected deletes=true for delete-customer.p on Customer")
			}
			return
		}
	}
	t.Error("expected Customer table entry for delete-customer.p")
}
