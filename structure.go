package sitebuilder

import (
    "fmt"
    "os"
    "path"
    "regexp"

    "io/ioutil"

    "github.com/dsoprea/go-logging"
)

var (
    pageIdRe *regexp.Regexp = nil
)

const (
    rootPageId                     = "index"
    defaultIdToLocalFilepathFormat = "%s.html"
    defaultPageIdRePhrase          = `^[A-Za-z0-9_,.\\-]+$`
)

// PageStatement defines one or more statements that represent a single widget
// or feature added by a single call to DialectPageBuilder.
type PageStatement struct {
    Type              WidgetType
    StatementMetadata map[string]interface{}
}

// PageContent describes all dialect-specific content for a page prior to
// generating HTML.
type PageContent struct {
    Statements   []PageStatement
    PageMetadata map[string]interface{}
}

// NewPageContent returns a new PageDialectContent struct.
func NewPageContent() (pc *PageContent) {
    return &PageContent{
        Statements:   make([]PageStatement, 0),
        PageMetadata: make(map[string]interface{}),
    }
}

// Add pushes a new statement into the page structure.
func (pc *PageContent) Add(ps PageStatement) {
    pc.Statements = append(pc.Statements, ps)
}

// SiteNode describes a single page and its children. This is the core utility
// for managing content.
//
// Specific members are public that we'd like to be able to serialize/store.
type SiteNode struct {
    sb                 *SiteBuilder
    intermediateOutput []byte
    finalOutput        []byte

    PageId    string
    PageTitle string
    Content   *PageContent

    Children []*SiteNode
}

func (sn *SiteNode) String() string {
    return fmt.Sprintf("SiteNode<PAGE-ID=[%s] PAGE-TITLE=[%s]>", sn.PageId, sn.PageTitle)
}

func NewSiteNode(sb *SiteBuilder, pageId string, pageTitle string) (sn *SiteNode) {
    content := NewPageContent()

    return &SiteNode{
        sb:        sb,
        PageId:    pageId,
        PageTitle: pageTitle,
        Content:   content,
        Children:  make([]*SiteNode, 0),
    }
}

func (sn *SiteNode) SetIntermediateOutput(intermediateOutput []byte) {
    sn.intermediateOutput = intermediateOutput
}

func (sn *SiteNode) IntermediateOutput() []byte {
    if sn.intermediateOutput == nil {
        log.Panicf("intermediate output not generated yet")
    }

    return sn.intermediateOutput
}

func (sn *SiteNode) SetFinalOutput(finalOutput []byte) {
    sn.finalOutput = finalOutput
}

func (sn *SiteNode) FinalOutput() []byte {
    if sn.finalOutput == nil {
        log.Panicf("final output not generated yet")
    }

    return sn.finalOutput
}

func (sn *SiteNode) Render() (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    err = sn.sb.dialect.RenderIntermediate(sn)
    log.PanicIf(err)

    for _, sn := range sn.Children {
        err := sn.Render()
        log.PanicIf(err)
    }

    err = sn.sb.dialect.RenderHtml(sn)
    log.PanicIf(err)

    return nil
}

// AddChildNode creates and appends a new child node for the current node and
// returns it.
func (sn *SiteNode) AddChildNode(pageId, pageTitle string) (childNode *SiteNode, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    if _, found := sn.sb.pageIndex[pageId]; found == true {
        log.Panicf("node with page-ID [%s] already exists", pageId)
    }

    if pageIdRe.MatchString(pageId) == false {
        log.Panicf("page-ID has an invalid format: [%s]", pageId)
    }

    childNode = NewSiteNode(sn.sb, pageId, pageTitle)
    sn.sb.pageIndex[childNode.PageId] = struct{}{}

    sn.Children = append(sn.Children, childNode)

    return childNode, nil
}

func (sn *SiteNode) SiteBuilder() *SiteBuilder {
    return sn.sb
}

func (sn *SiteNode) Builder() *PageBuilder {
    return NewPageBuilder(sn)
}

// SiteBuilder contains all nodes for the current site being built.
type SiteBuilder struct {
    dialect     Dialect
    rootNode    *SiteNode
    pageIndex   map[string]struct{}
    siteContext *SiteContext
}

type SiteContext struct {
    // htmlOutputPath is the path that HTML content will be written to.
    htmlOutputPath string

    // IdToLocalFilepathFormat is the filename template that we'll plug a page-
    // ID into in order to produce the final filename.
    idToLocalFilepathFormat string
}

func NewSiteContext(htmlOutputPath string) *SiteContext {
    return &SiteContext{
        htmlOutputPath:          htmlOutputPath,
        idToLocalFilepathFormat: defaultIdToLocalFilepathFormat,
    }
}

// func (sc *SiteContext) IdToLocalFilepathFormat() string {
//     return sc.toLocalFilepathFormat
// }

func (sc *SiteContext) SetIdToLocalFilepathFormat(format string) {
    sc.idToLocalFilepathFormat = format
}

func (sc *SiteContext) HtmlOutputPath() string {
    return sc.htmlOutputPath
}

func (sc *SiteContext) GetFinalPageFilename(pageId string) string {
    filename := fmt.Sprintf(sc.idToLocalFilepathFormat, pageId)
    return filename
}

func NewSiteBuilder(siteTitle string, dialect Dialect, siteContext *SiteContext) (sb *SiteBuilder) {
    pageIndex := map[string]struct{}{
        rootPageId: struct{}{},
    }

    sb = &SiteBuilder{
        dialect:     dialect,
        pageIndex:   pageIndex,
        siteContext: siteContext,
    }

    rootNode := NewSiteNode(sb, rootPageId, siteTitle)
    sb.rootNode = rootNode

    return sb
}

func (sb *SiteBuilder) PageIsValid(pageId string) bool {
    _, found := sb.pageIndex[pageId]
    return found
}

// Root is the root node (homepage) of the site.
func (sb *SiteBuilder) Context() (siteContext *SiteContext) {
    return sb.siteContext
}

// Root is the root node (homepage) of the site.
func (sb *SiteBuilder) Root() (rootNode *SiteNode) {
    return sb.rootNode
}

func (sb *SiteBuilder) WriteToPath() (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    err = sb.rootNode.Render()
    log.PanicIf(err)

    err = os.MkdirAll(sb.siteContext.htmlOutputPath, 0755)
    log.PanicIf(err)

    err = sb.writeToPath(sb.rootNode)
    log.PanicIf(err)

    return nil
}

func (sb *SiteBuilder) writeToPath(sn *SiteNode) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    filename := sb.Context().GetFinalPageFilename(sn.PageId)
    pageFilepath := path.Join(sb.siteContext.htmlOutputPath, filename)

    finalOutput := sn.FinalOutput()

    err = ioutil.WriteFile(pageFilepath, finalOutput, 0666)
    log.PanicIf(err)

    for _, childNode := range sn.Children {
        err := sb.writeToPath(childNode)
        log.PanicIf(err)
    }

    return nil
}

func init() {
    var err error

    pageIdRe, err = regexp.Compile(defaultPageIdRePhrase)
    log.PanicIf(err)
}
