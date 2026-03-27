package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bfv/xref/internal/models"
)

// ParserConfig controls which xref line types are processed.
type ParserConfig struct {
	Classes      bool
	Methods      bool
	Constructors bool
	Invokes      bool
	Interfaces   bool
	News         bool // instantiates
	Procedures   bool
	Runs         bool
}

// DefaultParserConfig returns a ParserConfig with all options enabled.
func DefaultParserConfig() ParserConfig {
	return ParserConfig{
		Classes:      true,
		Methods:      true,
		Constructors: true,
		Invokes:      true,
		Interfaces:   true,
		News:         true,
		Procedures:   true,
		Runs:         true,
	}
}

// Parser parses OpenEdge .xref files into structured XrefFile data.
type Parser struct {
	config         ParserConfig
	crudIgnore     []string
	crudTypes      []string
	typesToProcess []string
}

// NewParser creates a Parser with the given config. Pass nil for defaults.
func NewParser(config *ParserConfig) *Parser {
	cfg := DefaultParserConfig()
	if config != nil {
		cfg = *config
	}

	p := &Parser{
		config:     cfg,
		crudIgnore: []string{"PUBLIC-DATA-MEMBER", "PUBLIC-PROPERTY", "SHARED", "DATA-MEMBER"},
		crudTypes:  []string{"ACCESS", "UPDATE", "CREATE", "DELETE", "REFERENCE"},
		typesToProcess: []string{
			"ANNOTATION", "CLASS", "COMPILE", "CONSTRUCTOR", "CPINTERNAL", "CPSTREAM",
			"INCLUDE", "INTERFACE", "INVOKE", "METHOD", "NEW", "PROCEDURE",
			"PRIVATE-PROCEDURE", "RUN", "SEARCH",
		},
	}

	p.applyConfig()
	return p
}

func (p *Parser) applyConfig() {
	if !p.config.Classes {
		p.removeType("CLASS")
		p.config.Constructors = false
		p.config.Methods = false
		p.config.Interfaces = false
	}
	if !p.config.Constructors {
		p.removeType("CONSTRUCTOR")
	}
	if !p.config.Methods {
		p.removeType("METHOD")
	}
	if !p.config.Interfaces {
		p.removeType("INTERFACE")
	}
	if !p.config.Invokes {
		p.removeType("INVOKE")
	}
	if !p.config.News {
		p.removeType("NEW")
	}
	if !p.config.Procedures {
		p.removeType("PROCEDURE")
		p.removeType("PRIVATE-PROCEDURE")
	}
	if !p.config.Runs {
		p.removeType("RUN")
	}
}

func (p *Parser) removeType(t string) {
	filtered := p.typesToProcess[:0]
	for _, v := range p.typesToProcess {
		if v != t {
			filtered = append(filtered, v)
		}
	}
	p.typesToProcess = filtered
}

func (p *Parser) includeType(t string) bool {
	for _, v := range p.typesToProcess {
		if v == t {
			return true
		}
	}
	return false
}

func containsStr(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

// ParseDir parses all .xref files in dirname recursively. sourceBaseDir is optionally
// stripped from source file paths.
func (p *Parser) ParseDir(dirname string, sourceBaseDir string) []*models.XrefFile {
	var parsed []*models.XrefFile

	info, err := os.Stat(dirname)
	if err != nil || !info.IsDir() {
		return parsed
	}

	files := p.readFiles(dirname)
	for _, file := range files {
		if filepath.Ext(file) == ".xref" {
			xreffile := p.parseFile(file, dirname, sourceBaseDir)
			parsed = append(parsed, xreffile)
		}
	}

	p.postProcess(parsed)
	return parsed
}

func (p *Parser) parseFile(file, xrefBaseDir, sourceBaseDir string) *models.XrefFile {
	xreffile := models.NewXrefFile(file, xrefBaseDir)

	data, err := os.ReadFile(file)
	if err != nil {
		return xreffile
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		xl := &XrefLine{}
		xl.ParseLine(line)
		p.processXrefLine(xl, xreffile)
	}

	xreffile.SourceFile = normalizeSourceFilename(xreffile.SourceFile)
	if sourceBaseDir != "" {
		xreffile.SourceFile = strings.Replace(xreffile.SourceFile, normalizeSourceFilename(sourceBaseDir), "", 1)
	}

	xreffile.Finish()
	return xreffile
}

func (p *Parser) processXrefLine(xl *XrefLine, xf *models.XrefFile) {
	if containsStr(p.crudTypes, xl.Type) {
		parts := strings.SplitN(xl.Info, " ", 3)
		isSequence := len(parts) >= 2 && parts[1] == "SEQUENCE"
		if !isSequence {
			p.processCrud(xl, xf)
		}
		return
	}

	if !p.includeType(xl.Type) {
		return
	}

	switch xl.Type {
	case "ANNOTATION":
		p.processAnnotation(xl, xf)
	case "CLASS":
		p.processClass(xl, xf)
	case "COMPILE":
		p.processCompile(xl, xf)
	case "CONSTRUCTOR":
		p.processConstructor(xl, xf)
	case "CPINTERNAL":
		p.processCpInternal(xl, xf)
	case "CPSTREAM":
		p.processCpStream(xl, xf)
	case "INCLUDE":
		p.processInclude(xl, xf)
	case "INTERFACE":
		p.processInterface(xl, xf)
	case "INVOKE":
		p.processInvoke(xl, xf)
	case "METHOD":
		p.processMethod(xl, xf)
	case "NEW":
		p.processNew(xl, xf)
	case "PROCEDURE":
		p.processProcedure(xl, xf, false)
	case "PRIVATE-PROCEDURE":
		p.processProcedure(xl, xf, true)
	case "RUN":
		p.processRun(xl, xf)
	case "SEARCH":
		p.processSearch(xl, xf)
	}
}

func (p *Parser) processAnnotation(xl *XrefLine, xf *models.XrefFile) {
	xf.Annotations = append(xf.Annotations, xl.Info)
}

func (p *Parser) processCpInternal(xl *XrefLine, xf *models.XrefFile) {
	xf.CpInternal = xl.Info
}

func (p *Parser) processCpStream(xl *XrefLine, xf *models.XrefFile) {
	xf.CpStream = xl.Info
}

func (p *Parser) processClass(xl *XrefLine, xf *models.XrefFile) {
	entries := strings.Split(xl.Info, ",")
	if len(entries) < 7 {
		return
	}

	inheritsStr := strings.TrimSpace(strings.Replace(entries[1], "INHERITS", "", 1))
	var inherits []string
	if inheritsStr != "" {
		inherits = strings.Fields(inheritsStr)
	}

	implementsStr := strings.TrimSpace(strings.Replace(entries[2], "IMPLEMENTS", "", 1))
	var implements []string
	if implementsStr != "" {
		implements = strings.Fields(implementsStr)
	}

	classObj := &models.Class{
		Name:          entries[0],
		Inherits:      inherits,
		Implements:    implements,
		UseWidgetPool: entries[3] == "USE-WIDGET-POOL",
		Final:         entries[4] == "FINAL",
		Abstract:      entries[5] == "ABSTRACT",
		Serializable:  entries[6] == "SERIALIZABLE",
		Constructors:  []models.Constructor{},
		Methods:       []models.Method{},
	}

	xf.Class = classObj
}

func (p *Parser) processInterface(xl *XrefLine, xf *models.XrefFile) {
	entries := strings.Split(xl.Info, ",")
	if len(entries) < 2 {
		return
	}

	inheritsStr := strings.TrimSpace(strings.Replace(entries[1], "INHERITS", "", 1))
	var inherits []string
	if inheritsStr != "" {
		inherits = strings.Fields(inheritsStr)
	}

	xf.Interface = &models.Interface{
		Name:     entries[0],
		Inherits: inherits,
		Methods:  []models.Method{},
	}
}

func (p *Parser) processInclude(xl *XrefLine, xf *models.XrefFile) {
	includeFile := strings.ReplaceAll(xl.Info, "\"", "")
	if idx := strings.Index(includeFile, " "); idx >= 0 {
		includeFile = includeFile[:idx]
	}

	if !containsStr(xf.Includes, includeFile) {
		xf.Includes = append(xf.Includes, includeFile)
	}
}

func (p *Parser) processCrud(xl *XrefLine, xf *models.XrefFile) {
	tablePart := strings.Fields(xl.Info)
	if len(tablePart) == 0 {
		return
	}

	if containsStr(p.crudIgnore, tablePart[0]) {
		return
	}

	if len(tablePart) >= 2 && (tablePart[1] == "WORKFILE" || tablePart[1] == "TEMPTABLE") {
		return
	}

	if xl.Type == "ACCESS" || xl.Type == "UPDATE" || xl.Type == "REFERENCE" {
		tableInfo := strings.SplitN(tablePart[0], ".", 2)
		if len(tableInfo) == 2 {
			table := xf.AddTable(tableInfo[1], tableInfo[0], false, false, xl.Type == "UPDATE")
			if len(tablePart) > 1 {
				xf.AddField(tablePart[1], table, xl.Type == "UPDATE")
			}
		}
	} else if xl.Type == "CREATE" || xl.Type == "DELETE" {
		tableInfo := strings.SplitN(tablePart[0], ".", 2)
		if len(tableInfo) == 2 {
			xf.AddTable(tableInfo[1], tableInfo[0], xl.Type == "CREATE", xl.Type == "DELETE", false)
		}
	}
}

func (p *Parser) processNew(xl *XrefLine, xf *models.XrefFile) {
	if !containsStr(xf.Instantiates, xl.Info) {
		xf.Instantiates = append(xf.Instantiates, xl.Info)
	}
}

func (p *Parser) processInvoke(xl *XrefLine, xf *models.XrefFile) {
	parts := strings.SplitN(xl.Info, ":", 2)
	if len(parts) < 2 {
		return
	}
	fqclass := parts[0]
	method := parts[1]

	var classObj *models.MethodInvocation
	for i := range xf.Invokes {
		if xf.Invokes[i].Class == fqclass {
			classObj = &xf.Invokes[i]
			break
		}
	}
	if classObj == nil {
		xf.Invokes = append(xf.Invokes, models.MethodInvocation{
			Class:   fqclass,
			Methods: []string{},
		})
		classObj = &xf.Invokes[len(xf.Invokes)-1]
	}

	if !containsStr(classObj.Methods, method) {
		classObj.Methods = append(classObj.Methods, method)
	}
}

func (p *Parser) processCompile(xl *XrefLine, xf *models.XrefFile) {
	xf.SetSourceFile(xl.Info)
}

func (p *Parser) processProcedure(xl *XrefLine, xf *models.XrefFile, isPrivate bool) {
	procInfo := strings.SplitN(xl.Info, ",", 2)
	xf.Procedures = append(xf.Procedures, models.Procedure{
		Name:    procInfo[0],
		Private: isPrivate,
	})
}

func (p *Parser) processRun(xl *XrefLine, xf *models.XrefFile) {
	procInfo := strings.Fields(xl.Info)
	if len(procInfo) == 0 {
		return
	}

	name := procInfo[0]

	// Check if already tracked
	for _, r := range xf.Runs {
		if r.Name == name {
			return
		}
	}

	persistent := len(procInfo) > 1 && procInfo[1] == "PERSISTENT"
	dynamic := strings.HasPrefix(strings.ToLower(name), "value(")

	xf.Runs = append(xf.Runs, models.Run{
		Name:       name,
		Persistent: persistent,
		Dynamic:    dynamic,
	})
}

func (p *Parser) processSearch(xl *XrefLine, xf *models.XrefFile) {
	searchInfo := strings.Fields(xl.Info)

	if !containsStr(searchInfo, "TEMPTABLE") {
		return
	}

	var ttname string
	if len(searchInfo) > 0 && searchInfo[0] == "DATA-MEMBER" {
		if len(searchInfo) > 1 {
			parts := strings.SplitN(searchInfo[1], ":", 2)
			if len(parts) == 2 {
				ttname = parts[1]
			}
		}
	} else if len(searchInfo) > 0 {
		ttname = searchInfo[0]
	}

	if ttname == "" {
		return
	}

	for _, tt := range xf.TTNames {
		if strings.EqualFold(tt, ttname) {
			return
		}
	}
	xf.TTNames = append(xf.TTNames, ttname)
}

func (p *Parser) processMethod(xl *XrefLine, xf *models.XrefFile) {
	method := p.extractMethod(xl, false)
	if xf.Class != nil {
		xf.Class.Methods = append(xf.Class.Methods, method)
	} else if xf.Interface != nil {
		xf.Interface.Methods = append(xf.Interface.Methods, method)
	}
}

func (p *Parser) processConstructor(xl *XrefLine, xf *models.XrefFile) {
	method := p.extractMethod(xl, true)
	constructor := models.Constructor{
		Accessor:  method.Accessor,
		Static:    method.Static,
		Signature: method.Signature,
	}

	if xf.Class != nil {
		xf.Class.Constructors = append(xf.Class.Constructors, constructor)
	}
}

func (p *Parser) extractMethod(xl *XrefLine, isConstructor bool) models.Method {
	methodInfo := strings.Split(xl.Info, ",")

	method := models.Method{
		Accessor:  models.AccessorPublic,
		Signature: []models.Parameter{},
	}

	if len(methodInfo) > 0 {
		method.Accessor = models.Accessor(strings.ToLower(methodInfo[0]))
	}
	if len(methodInfo) > 1 {
		method.Static = methodInfo[1] == "STATIC"
	}
	if len(methodInfo) > 2 {
		method.Override = methodInfo[2] == "OVERRIDE"
	}
	if len(methodInfo) > 3 {
		method.Final = methodInfo[3] == "FINAL"
	}
	if len(methodInfo) > 4 {
		method.Abstract = methodInfo[4] == "ABSTRACT"
	}
	if len(methodInfo) > 5 {
		method.Name = methodInfo[5]
	}
	if !isConstructor && len(methodInfo) > 6 {
		method.ReturnType = methodInfo[6]
	}

	startIdx := 7
	if isConstructor {
		startIdx = 6
	}
	for i := startIdx; i < len(methodInfo); i++ {
		param := p.extractParameter(methodInfo[i])
		if param != nil {
			method.Signature = append(method.Signature, *param)
		}
	}

	return method
}

func (p *Parser) extractParameter(paramString string) *models.Parameter {
	if paramString == "" {
		return nil
	}

	paramInfo := strings.Fields(paramString)
	if len(paramInfo) < 2 {
		return nil
	}

	param := &models.Parameter{
		Name: paramInfo[1],
		Mode: models.ParameterMode(strings.ToLower(paramInfo[0])),
	}

	if len(paramInfo) > 2 {
		param.DataType = strings.ToLower(paramInfo[2])
	}

	return param
}

func (p *Parser) readFiles(dirname string) []string {
	var filelist []string

	entries, err := os.ReadDir(dirname)
	if err != nil {
		return filelist
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirname, entry.Name())
		if entry.IsDir() {
			filelist = append(filelist, p.readFiles(fullPath)...)
		} else {
			filelist = append(filelist, fullPath)
		}
	}

	return filelist
}

func normalizeSourceFilename(s string) string {
	return strings.ReplaceAll(s, "\\", "/")
}

func (p *Parser) postProcess(xreffiles []*models.XrefFile) {
	if p.config.Methods || p.config.Constructors {
		p.fixClassNames(xreffiles)
	}
}

func (p *Parser) fixClassNames(xreffiles []*models.XrefFile) {
	classes := p.getClassNames(xreffiles)

	for _, xf := range xreffiles {
		if xf.Class == nil {
			continue
		}

		if p.config.Methods {
			fixParameterDatatype(classes, methodSignatures(xf.Class.Methods))
		}

		if p.config.Constructors {
			fixParameterDatatype(classes, constructorSignatures(xf.Class.Constructors))
		}
	}
}

func fixParameterDatatype(classes []string, signatures [][]models.Parameter) {
	for _, sig := range signatures {
		for k := range sig {
			for _, className := range classes {
				if strings.EqualFold(className, sig[k].DataType) {
					sig[k].DataType = className
					break
				}
			}
		}
	}
}

func methodSignatures(methods []models.Method) [][]models.Parameter {
	sigs := make([][]models.Parameter, len(methods))
	for i := range methods {
		sigs[i] = methods[i].Signature
	}
	return sigs
}

func (p *Parser) getClassNames(xreffiles []*models.XrefFile) []string {
	var classes []string
	for _, xf := range xreffiles {
		if xf.Class != nil {
			classes = append(classes, xf.Class.Name)
		}
	}
	return classes
}

func constructorSignatures(constructors []models.Constructor) [][]models.Parameter {
	sigs := make([][]models.Parameter, len(constructors))
	for i := range constructors {
		sigs[i] = constructors[i].Signature
	}
	return sigs
}
