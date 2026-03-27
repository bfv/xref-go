package models

import (
	"sort"
	"strings"
)

// Accessor represents visibility of a class member.
type Accessor string

const (
	AccessorPublic    Accessor = "public"
	AccessorProtected Accessor = "protected"
	AccessorPrivate   Accessor = "private"
)

// ParameterMode represents the direction of a parameter.
type ParameterMode string

const (
	ParameterModeInput       ParameterMode = "input"
	ParameterModeOutput      ParameterMode = "output"
	ParameterModeInputOutput ParameterMode = "input-output"
	ParameterModeBuffer      ParameterMode = "buffer"
)

// Field represents a database field reference.
type Field struct {
	Name      string `json:"name"`
	IsUpdated bool   `json:"isUpdated"`
}

// Table represents a database table reference with CRUD flags.
type Table struct {
	Name      string  `json:"name"`
	Database  string  `json:"database"`
	IsCreated bool    `json:"isCreated"`
	IsDeleted bool    `json:"isDeleted"`
	IsUpdated bool    `json:"isUpdated"`
	Fields    []Field `json:"fields"`
}

// Class represents an OpenEdge class definition.
type Class struct {
	Name          string        `json:"name"`
	Inherits      []string      `json:"inherits"`
	Implements    []string      `json:"implements"`
	UseWidgetPool bool          `json:"useWidgetPool"`
	Final         bool          `json:"final"`
	Abstract      bool          `json:"abstract"`
	Serializable  bool          `json:"serializable"`
	Constructors  []Constructor `json:"constructors"`
	Methods       []Method      `json:"methods"`
}

// Interface represents an OpenEdge interface definition.
type Interface struct {
	Name     string   `json:"name"`
	Inherits []string `json:"inherits"`
	Methods  []Method `json:"methods"`
}

// Method represents a class or interface method.
type Method struct {
	Name       string      `json:"name"`
	Accessor   Accessor    `json:"accessor"`
	Static     bool        `json:"static"`
	Override   bool        `json:"override"`
	Final      bool        `json:"final"`
	Abstract   bool        `json:"abstract"`
	ReturnType string      `json:"returntype"`
	Signature  []Parameter `json:"signature"`
}

// Constructor represents a class constructor.
type Constructor struct {
	Accessor  Accessor    `json:"accessor"`
	Static    bool        `json:"static"`
	Signature []Parameter `json:"signature"`
}

// Parameter represents a method/constructor parameter.
type Parameter struct {
	Name     string        `json:"name"`
	Mode     ParameterMode `json:"mode"`
	DataType string        `json:"datatype"`
}

// MethodInvocation represents a method call on a class.
type MethodInvocation struct {
	Class   string   `json:"class"`
	Methods []string `json:"methods"`
}

// Procedure represents an internal procedure reference.
type Procedure struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
}

// Run represents a RUN statement reference.
type Run struct {
	Name       string `json:"name"`
	Persistent bool   `json:"persistent"`
	Dynamic    bool   `json:"dynamic"`
}

// TableDefinition holds a table name and its database.
type TableDefinition struct {
	Table    string `json:"table"`
	Database string `json:"database"`
}

// TempTableField represents a field in a temp-table.
type TempTableField struct {
	Name string `json:"name"`
}

// TempTable represents a temp-table definition.
type TempTable struct {
	Name   string           `json:"name"`
	Fields []TempTableField `json:"fields"`
}

// XrefFile is the per-source aggregation of all xref data.
type XrefFile struct {
	XrefFilePath string             `json:"xreffile"`
	SourceFile   string             `json:"sourcefile"`
	Class        *Class             `json:"class,omitempty"`
	Interface    *Interface         `json:"interface,omitempty"`
	CpInternal   string             `json:"cpInternal"`
	CpStream     string             `json:"cpStream"`
	Includes     []string           `json:"includes"`
	TableNames   []string           `json:"tablenames"`
	TTNames      []string           `json:"ttnames"`
	Instantiates []string           `json:"instantiates"`
	Invokes      []MethodInvocation `json:"invokes"`
	Annotations  []string           `json:"annotations"`
	Procedures   []Procedure        `json:"procedures"`
	Runs         []Run              `json:"runs"`
	Tables       []Table            `json:"tables"`
	TempTables   []TempTable        `json:"temptables"`
}

// NewXrefFile creates a new XrefFile with the xref file path normalized.
func NewXrefFile(file, xrefBaseDir string) *XrefFile {
	return &XrefFile{
		XrefFilePath: NormalizePath(file, xrefBaseDir),
		Includes:     []string{},
		TableNames:   []string{},
		TTNames:      []string{},
		Instantiates: []string{},
		Invokes:      []MethodInvocation{},
		Annotations:  []string{},
		Procedures:   []Procedure{},
		Runs:         []Run{},
		Tables:       []Table{},
		TempTables:   []TempTable{},
	}
}

// SetSourceFile sets the source file path.
func (xf *XrefFile) SetSourceFile(sourcefile string) {
	xf.SourceFile = sourcefile
}

// AddTable adds or updates a table reference. Returns the table.
func (xf *XrefFile) AddTable(tablename, db string, created, deleted, updated bool) *Table {
	// Track unique table names
	found := false
	for _, tn := range xf.TableNames {
		if tn == tablename {
			found = true
			break
		}
	}
	if !found {
		xf.TableNames = append(xf.TableNames, tablename)
	}

	// Find or create the table
	var table *Table
	for i := range xf.Tables {
		if xf.Tables[i].Name == tablename && xf.Tables[i].Database == db {
			table = &xf.Tables[i]
			break
		}
	}
	if table == nil {
		xf.Tables = append(xf.Tables, Table{
			Name:     tablename,
			Database: db,
			Fields:   []Field{},
		})
		table = &xf.Tables[len(xf.Tables)-1]
	}

	table.IsCreated = table.IsCreated || created
	table.IsDeleted = table.IsDeleted || deleted
	table.IsUpdated = table.IsUpdated || updated

	return table
}

// AddField adds or updates a field on a table.
func (xf *XrefFile) AddField(fieldname string, table *Table, isUpdated bool) {
	if fieldname == "" || strings.TrimSpace(fieldname) == "" {
		return
	}

	for i := range table.Fields {
		if table.Fields[i].Name == fieldname {
			table.Fields[i].IsUpdated = table.Fields[i].IsUpdated || isUpdated
			return
		}
	}

	table.Fields = append(table.Fields, Field{
		Name:      fieldname,
		IsUpdated: isUpdated,
	})
}

// Finish sorts the collected slices for deterministic output.
func (xf *XrefFile) Finish() {
	sort.Slice(xf.Instantiates, func(i, j int) bool {
		return strings.ToLower(xf.Instantiates[i]) < strings.ToLower(xf.Instantiates[j])
	})
	sort.Slice(xf.Includes, func(i, j int) bool {
		return strings.ToLower(xf.Includes[i]) < strings.ToLower(xf.Includes[j])
	})
	sort.Slice(xf.Invokes, func(i, j int) bool {
		return strings.ToLower(xf.Invokes[i].Class) < strings.ToLower(xf.Invokes[j].Class)
	})
	sort.Slice(xf.TableNames, func(i, j int) bool {
		return strings.ToLower(xf.TableNames[i]) < strings.ToLower(xf.TableNames[j])
	})
}

// NormalizePath normalizes a file path by replacing backslashes and removing the base dir prefix.
func NormalizePath(file, dir string) string {
	file = strings.ReplaceAll(file, "\\", "/")
	file = strings.ReplaceAll(file, "//", "/")
	dir = strings.ReplaceAll(dir, "\\", "/")
	dir = strings.ReplaceAll(dir, "//", "/")
	if dir != "" && dir != "." {
		if !strings.HasSuffix(dir, "/") {
			dir += "/"
		}
		file = strings.TrimPrefix(file, dir)
	}
	return strings.TrimPrefix(file, "/")
}
