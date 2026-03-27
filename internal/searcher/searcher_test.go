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
