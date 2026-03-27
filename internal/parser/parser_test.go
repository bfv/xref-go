package parser

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/bfv/xref/internal/models"
)

func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "testcases")
}

func TestParseDirNonExistent(t *testing.T) {
	p := NewParser(nil)
	result := p.ParseDir("C:\\nonexistent_xyz", "")
	if len(result) != 0 {
		t.Errorf("expected 0 XrefFiles for non-existent dir, got %d", len(result))
	}
}

func TestParseDirReturnsFiles(t *testing.T) {
	p := NewParser(nil)
	result := p.ParseDir(testdataDir(), "")
	if len(result) == 0 {
		t.Fatal("expected XrefFiles from testdata, got 0")
	}
}

func TestParseFileCustomerTables(t *testing.T) {
	p := NewParser(nil)
	result := p.ParseDir(testdataDir(), "")

	xf := findXrefFile(result, "db/customer.p.xref")
	if xf == nil {
		t.Fatal("customer.p.xref not found in results")
	}

	if len(xf.TableNames) != 1 {
		t.Errorf("expected 1 tablename, got %d", len(xf.TableNames))
	}

	found := false
	for _, tn := range xf.TableNames {
		if tn == "Customer" {
			found = true
		}
	}
	if !found {
		t.Error("expected Customer in tablenames")
	}

	if len(xf.Tables) != 1 {
		t.Errorf("expected 1 table, got %d", len(xf.Tables))
	}

	if strings.ToLower(xf.Tables[0].Name) != "customer" {
		t.Errorf("expected Customer table, got %s", xf.Tables[0].Name)
	}
}

func TestParseFileHelloWorld(t *testing.T) {
	p := NewParser(nil)
	result := p.ParseDir(testdataDir(), "")

	xf := findXrefFile(result, "hello/helloworld.p.xref")
	if xf == nil {
		t.Fatal("helloworld.p.xref not found in results")
	}

	if xf.CpInternal != "UTF-8" {
		t.Errorf("expected cpInternal=UTF-8, got %s", xf.CpInternal)
	}

	if xf.CpStream != "UTF-8" {
		t.Errorf("expected cpStream=UTF-8, got %s", xf.CpStream)
	}
}

func TestParseFileClassParsing(t *testing.T) {
	p := NewParser(nil)
	result := p.ParseDir(testdataDir(), "")

	xf := findXrefFile(result, "oo/CustomerBE.cls.xref")
	if xf == nil {
		t.Fatal("CustomerBE.cls.xref not found in results")
	}

	if xf.Class == nil {
		t.Fatal("expected class to be parsed")
	}

	if xf.Class.Name != "oo.CustomerBE" {
		t.Errorf("expected class name oo.CustomerBE, got %s", xf.Class.Name)
	}

	if len(xf.Class.Methods) != 2 {
		t.Errorf("expected 2 methods, got %d", len(xf.Class.Methods))
	}

	if len(xf.Class.Constructors) != 1 {
		t.Errorf("expected 1 constructor, got %d", len(xf.Class.Constructors))
	}

	if len(xf.Tables) != 1 {
		t.Errorf("expected 1 table, got %d", len(xf.Tables))
	}
	if len(xf.Tables) > 0 && xf.Tables[0].Name != "Customer" {
		t.Errorf("expected Customer table, got %s", xf.Tables[0].Name)
	}
}

func TestXrefLineParse(t *testing.T) {
	line := `C:\test\file.p C:\test\file.p 42 ACCESS sports2000.Customer CustNum`
	xl := &XrefLine{}
	xl.ParseLine(line)

	if xl.Type != "ACCESS" {
		t.Errorf("expected type ACCESS, got %s", xl.Type)
	}
	if xl.LineNumber != 42 {
		t.Errorf("expected line 42, got %d", xl.LineNumber)
	}
	if xl.Info != "sports2000.Customer CustNum" {
		t.Errorf("expected info 'sports2000.Customer CustNum', got '%s'", xl.Info)
	}
}

func findXrefFile(files []*models.XrefFile, name string) *models.XrefFile {
	for _, f := range files {
		if strings.HasSuffix(f.XrefFilePath, name) {
			return f
		}
	}
	return nil
}
