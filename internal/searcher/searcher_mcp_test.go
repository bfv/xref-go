package searcher

import (
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/bfv/xref/internal/datafile"
)

func examplesXrefJSON() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "testcases", "xref.json")
}

func setupMCPSearcher(t *testing.T) *Searcher {
	t.Helper()
	xreffiles, err := datafile.Load(examplesXrefJSON())
	if err != nil {
		t.Fatalf("failed to load xref.json: %v", err)
	}
	return NewSearcher(xreffiles)
}

// --- SearchSources ---

func TestSearchSources(t *testing.T) {
	s := setupMCPSearcher(t)

	results := s.SearchSources("db/")
	if len(results) != 11 {
		t.Fatalf("expected 11 db/ sources, got %d", len(results))
	}
	for _, src := range results {
		if src[:3] != "db/" {
			t.Errorf("expected db/ prefix, got %s", src)
		}
	}
}

func TestSearchSourcesWithWildcard(t *testing.T) {
	s := setupMCPSearcher(t)

	results := s.SearchSources("oo/*")
	if len(results) != 12 {
		t.Fatalf("expected 12 oo/ sources, got %d", len(results))
	}
}

func TestSearchSourcesCaseInsensitive(t *testing.T) {
	s := setupMCPSearcher(t)

	results := s.SearchSources("DB/")
	if len(results) != 11 {
		t.Fatalf("expected 11 sources with case-insensitive match, got %d", len(results))
	}
}

func TestSearchSourcesNoMatch(t *testing.T) {
	s := setupMCPSearcher(t)

	results := s.SearchSources("nonexistent/")
	if len(results) != 0 {
		t.Fatalf("expected 0 results for nonexistent prefix, got %d", len(results))
	}
}

func TestSearchSourcesTT(t *testing.T) {
	s := setupMCPSearcher(t)

	results := s.SearchSources("tt/")
	if len(results) != 7 {
		t.Fatalf("expected 7 tt/ sources, got %d", len(results))
	}
}

func TestSearchSourcesInclude(t *testing.T) {
	s := setupMCPSearcher(t)

	results := s.SearchSources("include/")
	if len(results) != 4 {
		t.Fatalf("expected 4 include/ sources, got %d", len(results))
	}
}

func TestSearchSourcesRun(t *testing.T) {
	s := setupMCPSearcher(t)

	results := s.SearchSources("run/")
	if len(results) != 5 {
		t.Fatalf("expected 5 run/ sources, got %d", len(results))
	}
}

// --- GetSummary ---

func TestGetSummary(t *testing.T) {
	s := setupMCPSearcher(t)

	summary := s.GetSummary()

	if summary.SourceCount != 40 {
		t.Errorf("expected 40 sources, got %d", summary.SourceCount)
	}
	if summary.ClassCount != 11 {
		t.Errorf("expected 11 classes, got %d", summary.ClassCount)
	}
	if summary.InterfaceCount != 3 {
		t.Errorf("expected 3 interfaces, got %d", summary.InterfaceCount)
	}
	if summary.DatabaseCount != 1 {
		t.Errorf("expected 1 database, got %d", summary.DatabaseCount)
	}
	if summary.TableCount != 3 {
		t.Errorf("expected 3 tables, got %d", summary.TableCount)
	}
	if summary.IncludeCount != 8 {
		t.Errorf("expected 8 unique includes, got %d", summary.IncludeCount)
	}
	if summary.ProcedureCount != 6 {
		t.Errorf("expected 6 procedures, got %d", summary.ProcedureCount)
	}
}

// --- GetRunReferences ---

func TestGetRunReferences(t *testing.T) {
	s := setupMCPSearcher(t)

	results := s.GetRunReferences("xref/oo/proc1.p")
	if len(results) != 1 {
		t.Fatalf("expected 1 source running xref/oo/proc1.p, got %d", len(results))
	}
	if results[0] != "run/master.p" {
		t.Errorf("expected run/master.p, got %s", results[0])
	}
}

func TestGetRunReferencesMultiple(t *testing.T) {
	s := setupMCPSearcher(t)

	results := s.GetRunReferences("xref/oo/proc3.p")
	if len(results) != 2 {
		t.Fatalf("expected 2 sources, got %d: %v", len(results), results)
	}
	sort.Strings(results)
	if results[0] != "run/proc1.p" || results[1] != "run/proc2.p" {
		t.Errorf("expected [run/proc1.p, run/proc2.p], got %v", results)
	}
}

func TestGetRunReferencesProc2(t *testing.T) {
	s := setupMCPSearcher(t)

	results := s.GetRunReferences("xref/oo/proc2.p")
	if len(results) != 2 {
		t.Fatalf("expected 2 sources running xref/oo/proc2.p, got %d: %v", len(results), results)
	}
	sort.Strings(results)
	if results[0] != "run/master.p" || results[1] != "run/proc1.p" {
		t.Errorf("expected [run/master.p, run/proc1.p], got %v", results)
	}
}

func TestGetRunReferencesNoMatch(t *testing.T) {
	s := setupMCPSearcher(t)

	results := s.GetRunReferences("nonexistent.p")
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestGetRunReferencesIncludeCode1(t *testing.T) {
	s := setupMCPSearcher(t)

	results := s.GetRunReferences("include/code1.p")
	if len(results) != 1 {
		t.Fatalf("expected 1 source running include/code1.p, got %d", len(results))
	}
	if results[0] != "include/code2.p" {
		t.Errorf("expected include/code2.p, got %s", results[0])
	}
}

// --- GetClassReferences ---

func TestGetClassReferencesAddress(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetClassReferences("oo.Address")
	if len(refs) != 2 {
		t.Fatalf("expected 2 class references for oo.Address, got %d: %v", len(refs), refs)
	}

	found := map[string]bool{}
	for _, ref := range refs {
		found[ref.Source] = true
		switch ref.Source {
		case "oo/Person.cls":
			if !ref.Instantiates {
				t.Error("expected oo/Person.cls instantiates=true")
			}
			if !ref.Invokes {
				t.Error("expected oo/Person.cls invokes=true")
			}
		case "oo/DeliverAddress.cls":
			if !ref.Instantiates {
				t.Error("expected oo/DeliverAddress.cls instantiates=true")
			}
			if !ref.Inherits {
				t.Error("expected oo/DeliverAddress.cls inherits=true")
			}
			if !ref.Invokes {
				t.Error("expected oo/DeliverAddress.cls invokes=true")
			}
		}
	}
	if !found["oo/Person.cls"] {
		t.Error("expected oo/Person.cls in references")
	}
	if !found["oo/DeliverAddress.cls"] {
		t.Error("expected oo/DeliverAddress.cls in references")
	}
}

func TestGetClassReferencesAddressBase(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetClassReferences("oo.AddressBase")
	if len(refs) != 2 {
		t.Fatalf("expected 2 references for oo.AddressBase, got %d: %v", len(refs), refs)
	}

	found := map[string]bool{}
	for _, ref := range refs {
		found[ref.Source] = true
		if ref.Source == "oo/Address.cls" {
			if !ref.Inherits {
				t.Error("expected oo/Address.cls inherits=true for oo.AddressBase")
			}
			if !ref.Invokes {
				t.Error("expected oo/Address.cls invokes=true for oo.AddressBase")
			}
		}
		if ref.Source == "oo/DeliverAddress.cls" {
			if !ref.Inherits {
				t.Error("expected oo/DeliverAddress.cls inherits=true for oo.AddressBase")
			}
		}
	}
	if !found["oo/Address.cls"] || !found["oo/DeliverAddress.cls"] {
		t.Errorf("expected both oo/Address.cls and oo/DeliverAddress.cls, got %v", refs)
	}
}

func TestGetClassReferencesPerson(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetClassReferences("oo.Person")
	if len(refs) != 1 {
		t.Fatalf("expected 1 reference for oo.Person, got %d: %v", len(refs), refs)
	}
	if refs[0].Source != "oo/Employee.cls" {
		t.Errorf("expected oo/Employee.cls, got %s", refs[0].Source)
	}
	if !refs[0].Inherits {
		t.Error("expected inherits=true")
	}
	if !refs[0].Invokes {
		t.Error("expected invokes=true")
	}
}

func TestGetClassReferencesNoMatch(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetClassReferences("nonexistent.class")
	if len(refs) != 0 {
		t.Fatalf("expected 0 results, got %d", len(refs))
	}
}

// --- GetInterfaceNames ---

func TestGetInterfaceNames(t *testing.T) {
	s := setupMCPSearcher(t)

	interfaces := s.GetInterfaceNames()
	if len(interfaces) != 3 {
		t.Fatalf("expected 3 interfaces, got %d: %v", len(interfaces), interfaces)
	}

	expected := map[string]bool{"oo.IBase": true, "oo.IDisposable": true, "oo.IEmpty": true}
	for _, iface := range interfaces {
		if !expected[iface] {
			t.Errorf("unexpected interface: %s", iface)
		}
	}
}

// --- GetSourceNames / GetSourceByName ---

func TestGetSourceNames(t *testing.T) {
	s := setupMCPSearcher(t)

	sources := s.GetSourceNames()
	if len(sources) != 40 {
		t.Fatalf("expected 40 sources, got %d", len(sources))
	}
	// Verify sorted
	for i := 1; i < len(sources); i++ {
		if sources[i-1] > sources[i] {
			t.Errorf("sources not sorted: %s > %s", sources[i-1], sources[i])
		}
	}
}

func TestGetSourceByName(t *testing.T) {
	s := setupMCPSearcher(t)

	xf := s.GetSourceByName("oo/Address.cls")
	if xf == nil {
		t.Fatal("expected to find oo/Address.cls")
	}
	if xf.Class == nil {
		t.Fatal("expected class info for oo/Address.cls")
	}
	if xf.Class.Name != "oo.Address" {
		t.Errorf("expected class name oo.Address, got %s", xf.Class.Name)
	}
}

func TestGetSourceByNameCaseInsensitive(t *testing.T) {
	s := setupMCPSearcher(t)

	xf := s.GetSourceByName("OO/ADDRESS.CLS")
	if xf == nil {
		t.Fatal("expected case-insensitive match for OO/ADDRESS.CLS")
	}
}

func TestGetSourceByNameNotFound(t *testing.T) {
	s := setupMCPSearcher(t)

	xf := s.GetSourceByName("nonexistent.p")
	if xf != nil {
		t.Error("expected nil for nonexistent source")
	}
}

// --- GetDatabaseNames / GetDatabaseReferences ---

func TestGetDatabaseNamesMCP(t *testing.T) {
	s := setupMCPSearcher(t)

	dbs := s.GetDatabaseNames(nil)
	if len(dbs) != 1 {
		t.Fatalf("expected 1 database, got %d: %v", len(dbs), dbs)
	}
	if dbs[0] != "sports2000" {
		t.Errorf("expected sports2000, got %s", dbs[0])
	}
}

func TestGetDatabaseReferencesMCP(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetDatabaseReferences("sports2000")
	if len(refs) != 14 {
		t.Fatalf("expected 14 sources referencing sports2000, got %d", len(refs))
	}
}

func TestGetDatabaseReferencesNoMatch(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetDatabaseReferences("nonexistent")
	if len(refs) != 0 {
		t.Fatalf("expected 0 results, got %d", len(refs))
	}
}

// --- GetTableNames / GetTableReferences ---

func TestGetTableNamesMCP(t *testing.T) {
	s := setupMCPSearcher(t)

	tables := s.GetTableNames(nil)
	if len(tables) != 3 {
		t.Fatalf("expected 3 tables, got %d: %v", len(tables), tables)
	}

	expected := map[string]bool{
		"sports2000.Customer":  true,
		"sports2000.Order":     true,
		"sports2000.OrderLine": true,
	}
	for _, td := range tables {
		key := td.Database + "." + td.Table
		if !expected[key] {
			t.Errorf("unexpected table: %s", key)
		}
	}
}

func TestGetTableReferencesCustomer(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetTableReferences("Customer", nil, nil, nil)
	if len(refs) != 12 {
		t.Fatalf("expected 12 sources referencing Customer, got %d", len(refs))
	}
}

func TestGetTableReferencesOrder(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetTableReferences("Order", nil, nil, nil)
	if len(refs) != 5 {
		t.Fatalf("expected 5 sources referencing Order, got %d", len(refs))
	}
}

func TestGetTableReferencesOrderLine(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetTableReferences("OrderLine", nil, nil, nil)
	if len(refs) != 2 {
		t.Fatalf("expected 2 sources referencing OrderLine, got %d", len(refs))
	}
}

func TestGetTableReferencesWithCreates(t *testing.T) {
	s := setupMCPSearcher(t)

	creates := true
	refs := s.GetTableReferences("Customer", &creates, nil, nil)
	if len(refs) != 1 {
		t.Fatalf("expected 1 source creating Customer, got %d", len(refs))
	}
	if refs[0].SourceFile != "db/create-customer.p" {
		t.Errorf("expected db/create-customer.p, got %s", refs[0].SourceFile)
	}
}

func TestGetTableReferencesWithDeletes(t *testing.T) {
	s := setupMCPSearcher(t)

	deletes := true
	refs := s.GetTableReferences("Customer", nil, nil, &deletes)
	if len(refs) != 1 {
		t.Fatalf("expected 1 source deleting Customer, got %d", len(refs))
	}
	if refs[0].SourceFile != "db/delete-customer.p" {
		t.Errorf("expected db/delete-customer.p, got %s", refs[0].SourceFile)
	}
}

func TestGetTableReferencesWithUpdatesMCP(t *testing.T) {
	s := setupMCPSearcher(t)

	updates := true
	refs := s.GetTableReferences("Customer", nil, &updates, nil)
	// create-customer.p, test_reference_update.p, update_credlimit.p
	if len(refs) != 3 {
		sources := make([]string, len(refs))
		for i, r := range refs {
			sources[i] = r.SourceFile
		}
		t.Fatalf("expected 3 sources updating Customer, got %d: %v", len(refs), sources)
	}
}

func TestGetTableReferencesNoMatch(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetTableReferences("NonExistentTable", nil, nil, nil)
	if len(refs) != 0 {
		t.Fatalf("expected 0 results, got %d", len(refs))
	}
}

// --- GetFieldReferences ---

func TestGetFieldReferencesCustNum(t *testing.T) {
	s := setupMCPSearcher(t)

	table := "Customer"
	refs := s.GetFieldReferences("CustNum", &table, nil)
	// CustNum in Customer appears in: create-customer.p, customer.p, custorderline.p,
	// delete-customer.p, tables_tt.p, test_reference.p, test_reference_only.p,
	// test_reference_update.p
	if len(refs) < 8 {
		t.Fatalf("expected at least 8 sources referencing Customer.CustNum, got %d", len(refs))
	}
}

func TestGetFieldReferencesCreditLimitUpdated(t *testing.T) {
	s := setupMCPSearcher(t)

	table := "Customer"
	updates := true
	refs := s.GetFieldReferences("CreditLimit", &table, &updates)
	// test_reference_update.p and update_credlimit.p
	if len(refs) != 2 {
		sources := make([]string, len(refs))
		for i, r := range refs {
			sources[i] = r.SourceFile
		}
		t.Fatalf("expected 2 sources updating Customer.CreditLimit, got %d: %v", len(refs), sources)
	}
}

func TestGetFieldReferencesNoTable(t *testing.T) {
	s := setupMCPSearcher(t)

	// Search CustNum across all tables
	refs := s.GetFieldReferences("CustNum", nil, nil)
	if len(refs) < 8 {
		t.Fatalf("expected at least 8 sources referencing CustNum (any table), got %d", len(refs))
	}
}

// --- GetIncludeReferences ---

func TestGetIncludeReferencesTTSports(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetIncludeReferences("tt/ttsports.i")
	if len(refs) != 7 {
		sources := make([]string, len(refs))
		for i, r := range refs {
			sources[i] = r.SourceFile
		}
		t.Fatalf("expected 7 sources including tt/ttsports.i, got %d: %v", len(refs), sources)
	}
}

func TestGetIncludeReferencesBaseIMCP(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetIncludeReferences("include/base.i")
	if len(refs) != 2 {
		t.Fatalf("expected 2 sources including include/base.i, got %d", len(refs))
	}

	sources := make([]string, len(refs))
	for i, r := range refs {
		sources[i] = r.SourceFile
	}
	sort.Strings(sources)
	if sources[0] != "include/includes-params.p" || sources[1] != "include/includes.p" {
		t.Errorf("expected [include/includes-params.p, include/includes.p], got %v", sources)
	}
}

func TestGetIncludeReferencesCode2I(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetIncludeReferences("include/code2.i")
	if len(refs) != 2 {
		t.Fatalf("expected 2 sources including include/code2.i, got %d", len(refs))
	}
}

func TestGetIncludeReferencesTTTest(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetIncludeReferences("tt/tttest.i")
	if len(refs) != 1 {
		t.Fatalf("expected 1 source including tt/tttest.i, got %d", len(refs))
	}
	if refs[0].SourceFile != "tt/TTTestClass.cls" {
		t.Errorf("expected tt/TTTestClass.cls, got %s", refs[0].SourceFile)
	}
}

func TestGetIncludeReferencesNoMatch(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetIncludeReferences("nonexistent.i")
	if len(refs) != 0 {
		t.Fatalf("expected 0 results, got %d", len(refs))
	}
}

// --- GetImplementations ---

func TestGetImplementationsIEmpty(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetImplementations("oo.IEmpty")
	if len(refs) != 2 {
		t.Fatalf("expected 2 implementors of oo.IEmpty, got %d", len(refs))
	}

	names := make([]string, len(refs))
	for i, r := range refs {
		if r.Class != nil {
			names[i] = r.Class.Name
		}
	}
	sort.Strings(names)
	if names[0] != "oo.Empty" || names[1] != "oo.MultipleImplements" {
		t.Errorf("expected [oo.Empty, oo.MultipleImplements], got %v", names)
	}
}

func TestGetImplementationsIBase(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetImplementations("oo.IBase")
	if len(refs) != 1 {
		t.Fatalf("expected 1 implementor of oo.IBase, got %d", len(refs))
	}
	if refs[0].Class.Name != "oo.AddressBase" {
		t.Errorf("expected oo.AddressBase, got %s", refs[0].Class.Name)
	}
}

func TestGetImplementationsIDisposable(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetImplementations("oo.IDisposable")
	if len(refs) != 1 {
		t.Fatalf("expected 1 implementor of oo.IDisposable, got %d", len(refs))
	}
	if refs[0].Class.Name != "oo.MultipleImplements" {
		t.Errorf("expected oo.MultipleImplements, got %s", refs[0].Class.Name)
	}
}

func TestGetImplementationsNoMatch(t *testing.T) {
	s := setupMCPSearcher(t)

	refs := s.GetImplementations("nonexistent.IFace")
	if len(refs) != 0 {
		t.Fatalf("expected 0 results, got %d", len(refs))
	}
}

// --- GetDependencies ---

func TestGetDependenciesDeliverAddress(t *testing.T) {
	s := setupMCPSearcher(t)

	deps := s.GetDependencies("oo/DeliverAddress.cls")
	if deps == nil {
		t.Fatal("expected dependencies for oo/DeliverAddress.cls")
	}
	if deps.Source != "oo/DeliverAddress.cls" {
		t.Errorf("expected source oo/DeliverAddress.cls, got %s", deps.Source)
	}
	if len(deps.Instantiates) != 1 || deps.Instantiates[0] != "oo.Address" {
		t.Errorf("expected instantiates [oo.Address], got %v", deps.Instantiates)
	}
	if len(deps.Invokes) != 1 || deps.Invokes[0].Class != "oo.Address" {
		t.Errorf("expected invokes oo.Address, got %v", deps.Invokes)
	}
	if deps.Class == nil || deps.Class.Name != "oo.DeliverAddress" {
		t.Error("expected class oo.DeliverAddress")
	}
	if len(deps.Tables) != 0 {
		t.Errorf("expected 0 tables, got %d", len(deps.Tables))
	}
}

func TestGetDependenciesCustomer(t *testing.T) {
	s := setupMCPSearcher(t)

	deps := s.GetDependencies("db/customer.p")
	if deps == nil {
		t.Fatal("expected dependencies for db/customer.p")
	}
	if len(deps.Tables) != 1 {
		t.Fatalf("expected 1 table, got %d", len(deps.Tables))
	}
	if deps.Tables[0].Name != "Customer" {
		t.Errorf("expected Customer table, got %s", deps.Tables[0].Name)
	}
}

func TestGetDependenciesMaster(t *testing.T) {
	s := setupMCPSearcher(t)

	deps := s.GetDependencies("run/master.p")
	if deps == nil {
		t.Fatal("expected dependencies for run/master.p")
	}
	if len(deps.Runs) != 5 {
		t.Fatalf("expected 5 runs, got %d", len(deps.Runs))
	}
}

func TestGetDependenciesNotFoundMCP(t *testing.T) {
	s := setupMCPSearcher(t)

	deps := s.GetDependencies("nonexistent.p")
	if deps != nil {
		t.Error("expected nil for nonexistent source")
	}
}

// --- GetClassHierarchy ---

func TestGetClassHierarchyDeliverAddress(t *testing.T) {
	s := setupMCPSearcher(t)

	hierarchy := s.GetClassHierarchy("oo.DeliverAddress")
	// oo.DeliverAddress → oo.Address → oo.AddressBase → oo.IBase = 4 entries
	if len(hierarchy) != 4 {
		names := make([]string, len(hierarchy))
		for i, h := range hierarchy {
			names[i] = h.Name
		}
		t.Fatalf("expected 4 hierarchy entries, got %d: %v", len(hierarchy), names)
	}
	if hierarchy[0].Name != "oo.DeliverAddress" {
		t.Errorf("expected root oo.DeliverAddress, got %s", hierarchy[0].Name)
	}
	if hierarchy[0].Type != "class" {
		t.Errorf("expected type class, got %s", hierarchy[0].Type)
	}
}

func TestGetClassHierarchyIDisposable(t *testing.T) {
	s := setupMCPSearcher(t)

	hierarchy := s.GetClassHierarchy("oo.IDisposable")
	// oo.IDisposable → oo.IEmpty = 2 entries
	if len(hierarchy) != 2 {
		t.Fatalf("expected 2 hierarchy entries, got %d", len(hierarchy))
	}
	if hierarchy[0].Name != "oo.IDisposable" || hierarchy[0].Type != "interface" {
		t.Errorf("expected root oo.IDisposable (interface), got %s (%s)", hierarchy[0].Name, hierarchy[0].Type)
	}
	if hierarchy[1].Name != "oo.IEmpty" || hierarchy[1].Type != "interface" {
		t.Errorf("expected oo.IEmpty (interface), got %s (%s)", hierarchy[1].Name, hierarchy[1].Type)
	}
}

func TestGetClassHierarchyFlat(t *testing.T) {
	s := setupMCPSearcher(t)

	// oo.Person has no parents
	hierarchy := s.GetClassHierarchy("oo.Person")
	if len(hierarchy) != 1 {
		t.Fatalf("expected 1 hierarchy entry, got %d", len(hierarchy))
	}
	if hierarchy[0].Name != "oo.Person" {
		t.Errorf("expected oo.Person, got %s", hierarchy[0].Name)
	}
}

// --- GetReverseDependencies ---

func TestGetReverseDependenciesAddressBase(t *testing.T) {
	s := setupMCPSearcher(t)

	rd := s.GetReverseDependencies("oo/AddressBase.cls")
	if rd == nil {
		t.Fatal("expected reverse dependencies for oo/AddressBase.cls")
	}

	sort.Strings(rd.InheritedBy)
	if len(rd.InheritedBy) != 2 {
		t.Fatalf("expected 2 inheritedBy, got %d: %v", len(rd.InheritedBy), rd.InheritedBy)
	}
	if rd.InheritedBy[0] != "oo/Address.cls" || rd.InheritedBy[1] != "oo/DeliverAddress.cls" {
		t.Errorf("expected [oo/Address.cls, oo/DeliverAddress.cls], got %v", rd.InheritedBy)
	}

	if len(rd.InvokedBy) != 1 || rd.InvokedBy[0] != "oo/Address.cls" {
		t.Errorf("expected invokedBy [oo/Address.cls], got %v", rd.InvokedBy)
	}

	if len(rd.IncludedBy) != 0 {
		t.Errorf("expected 0 includedBy, got %d", len(rd.IncludedBy))
	}
}

func TestGetReverseDependenciesIEmpty(t *testing.T) {
	s := setupMCPSearcher(t)

	rd := s.GetReverseDependencies("oo/IEmpty.cls")
	if rd == nil {
		t.Fatal("expected reverse dependencies for oo/IEmpty.cls")
	}

	sort.Strings(rd.InheritedBy)
	// oo/Empty.cls (implements), oo/IDisposable.cls (inherits), oo/MultipleImplements.cls (implements)
	if len(rd.InheritedBy) != 3 {
		t.Fatalf("expected 3 inheritedBy for oo/IEmpty.cls, got %d: %v", len(rd.InheritedBy), rd.InheritedBy)
	}
}

func TestGetReverseDependenciesAddress(t *testing.T) {
	s := setupMCPSearcher(t)

	rd := s.GetReverseDependencies("oo/Address.cls")
	if rd == nil {
		t.Fatal("expected reverse dependencies for oo/Address.cls")
	}

	if len(rd.InstantiatedBy) != 2 {
		t.Fatalf("expected 2 instantiatedBy, got %d: %v", len(rd.InstantiatedBy), rd.InstantiatedBy)
	}
	if len(rd.InheritedBy) != 1 || rd.InheritedBy[0] != "oo/DeliverAddress.cls" {
		t.Errorf("expected inheritedBy [oo/DeliverAddress.cls], got %v", rd.InheritedBy)
	}
}

func TestGetReverseDependenciesNotFoundMCP(t *testing.T) {
	s := setupMCPSearcher(t)

	rd := s.GetReverseDependencies("nonexistent.p")
	if rd != nil {
		t.Error("expected nil for nonexistent source")
	}
}

// --- GetMigrationScope ---

func TestGetMigrationScopeFromCustomer(t *testing.T) {
	s := setupMCPSearcher(t)

	scope := s.GetMigrationScope("db/customer.p")
	if scope == nil {
		t.Fatal("expected migration scope for db/customer.p")
	}
	if scope.StartSource != "db/customer.p" {
		t.Errorf("expected startSource db/customer.p, got %s", scope.StartSource)
	}
	// Should pull in all Customer/Order/OrderLine sources + include chains + tt sources
	if len(scope.Sources) < 15 {
		t.Errorf("expected at least 15 sources in migration scope, got %d: %v", len(scope.Sources), scope.Sources)
	}
	if len(scope.Tables) != 3 {
		t.Errorf("expected 3 tables in migration scope, got %d: %v", len(scope.Tables), scope.Tables)
	}
}

func TestGetMigrationScopeIsolated(t *testing.T) {
	s := setupMCPSearcher(t)

	// hello/helloworld.p has no tables, no includes, no class → only itself
	scope := s.GetMigrationScope("hello/helloworld.p")
	if scope == nil {
		t.Fatal("expected migration scope for hello/helloworld.p")
	}
	if len(scope.Sources) != 1 {
		t.Errorf("expected 1 source in isolated scope, got %d: %v", len(scope.Sources), scope.Sources)
	}
	if scope.Sources[0] != "hello/helloworld.p" {
		t.Errorf("expected hello/helloworld.p, got %s", scope.Sources[0])
	}
}

func TestGetMigrationScopeNotFoundMCP(t *testing.T) {
	s := setupMCPSearcher(t)

	scope := s.GetMigrationScope("nonexistent.p")
	if scope != nil {
		t.Error("expected nil for nonexistent source")
	}
}

// --- GetCrudMatrix ---

func TestGetCrudMatrixAll(t *testing.T) {
	s := setupMCPSearcher(t)

	matrix := s.GetCrudMatrix(nil)
	if len(matrix.Tables) != 3 {
		t.Errorf("expected 3 tables in CRUD matrix, got %d", len(matrix.Tables))
	}
	if len(matrix.Sources) != 14 {
		t.Errorf("expected 14 sources in CRUD matrix, got %d", len(matrix.Sources))
	}
	if len(matrix.Entries) == 0 {
		t.Error("expected non-empty entries")
	}
}

func TestGetCrudMatrixFilteredMCP(t *testing.T) {
	s := setupMCPSearcher(t)

	matrix := s.GetCrudMatrix([]string{"db/create-customer.p", "db/delete-customer.p"})
	if len(matrix.Sources) != 2 {
		t.Errorf("expected 2 sources, got %d", len(matrix.Sources))
	}
	if len(matrix.Tables) != 1 {
		t.Errorf("expected 1 table (Customer), got %d: %v", len(matrix.Tables), matrix.Tables)
	}

	for _, entry := range matrix.Entries {
		if entry.Source == "db/create-customer.p" {
			if !entry.Creates || !entry.Updates {
				t.Errorf("expected creates=true, updates=true for create-customer.p, got c=%v u=%v", entry.Creates, entry.Updates)
			}
		}
		if entry.Source == "db/delete-customer.p" {
			if !entry.Deletes {
				t.Errorf("expected deletes=true for delete-customer.p")
			}
		}
	}
}

func TestGetCrudMatrixNoTables(t *testing.T) {
	s := setupMCPSearcher(t)

	// hello/helloworld.p has no tables
	matrix := s.GetCrudMatrix([]string{"hello/helloworld.p"})
	if len(matrix.Sources) != 0 {
		t.Errorf("expected 0 sources (no tables), got %d", len(matrix.Sources))
	}
	if len(matrix.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(matrix.Entries))
	}
}
