package sitebuilder

import (
    "fmt"
    "io/ioutil"
    "path"

    "github.com/dsoprea/go-logging"
)

const (
    rootPageId                     = "index"
    defaultIdToLocalFilepathFormat = "%s.html"
)

type WidgetType int

// These should represent every method in the PageDialectBuilder interface.
const (
    ContentImage WidgetType = 1 + iota
    ChildNavbar
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

// AddNode creates and appends a new child node for the current node and returns
// it.
func (sn *SiteNode) AddChild(pageId, pageTitle string) (childNode *SiteNode) {
    if _, found := sn.sb.pageIndex[pageId]; found == true {
        log.Panicf("node with page-ID [%s] already exists", pageId)
    }

    childNode = NewSiteNode(sn.sb, pageId, pageTitle)
    sn.sb.pageIndex[childNode.PageId] = struct{}{}

    sn.Children = append(sn.Children, childNode)

    return childNode
}

func (sn *SiteNode) Builder() *PageBuilder {
    return NewPageBuilder(sn)
}

// SiteBuilder contains all nodes for the current site being built.
type SiteBuilder struct {
    dialect                 Dialect
    rootNode                *SiteNode
    pageIndex               map[string]struct{}
    idToLocalFilepathFormat string
}

func NewSiteBuilder(siteTitle string, dialect Dialect) (sb *SiteBuilder) {
    pageIndex := map[string]struct{}{
        rootPageId: struct{}{},
    }

    sb = &SiteBuilder{
        dialect:                 dialect,
        pageIndex:               pageIndex,
        idToLocalFilepathFormat: defaultIdToLocalFilepathFormat,
    }

    rootNode := NewSiteNode(sb, rootPageId, siteTitle)
    sb.rootNode = rootNode

    return sb
}

func (sb *SiteBuilder) SetIdToLocalFilepathFormat(format string) {
    sb.idToLocalFilepathFormat = format
}

// Root is the root node (homepage) of the site.
func (sb *SiteBuilder) Root() (rootNode *SiteNode) {
    return sb.rootNode
}

func (sb *SiteBuilder) WriteToPath(rootPath string) (err error) {
    err = sb.rootNode.Render()
    log.PanicIf(err)

    err = sb.writeToPath(sb.rootNode, rootPath)
    log.PanicIf(err)

    return nil
}

func (sb *SiteBuilder) writeToPath(sn *SiteNode, rootPath string) (err error) {
    filename := fmt.Sprintf(sb.idToLocalFilepathFormat, sn.PageId)
    pageFilepath := path.Join(rootPath, filename)

    finalOutput := sn.FinalOutput()

    err = ioutil.WriteFile(pageFilepath, finalOutput, 0666)
    log.PanicIf(err)

    for _, childNode := range sn.Children {
        err := sb.writeToPath(childNode, rootPath)
        log.PanicIf(err)
    }

    return nil
}
