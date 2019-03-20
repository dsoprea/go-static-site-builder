package markdowndialect

import (
    "bytes"
    "fmt"
    "io"

    "github.com/dsoprea/go-logging"
    "gopkg.in/russross/blackfriday.v2"

    "github.com/dsoprea/go-static-site-builder"
)

type MarkdownDialect struct {
}

func NewMarkdownDialect() (md *MarkdownDialect) {
    return &MarkdownDialect{}
}

// RenderIntermediate produces dialect-specific content that can be passed to
// RenderHtml.
func (md *MarkdownDialect) RenderIntermediate(sn *sitebuilder.SiteNode) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    // TODO(dustin): !! This is an overflow concern, especially with large embedded images.
    b := new(bytes.Buffer)

    _, err = fmt.Fprintf(b, "# %s\n\n", sn.PageTitle)
    log.PanicIf(err)

    for _, ps := range sn.Content.Statements {
        err := md.renderStatment(b, ps)
        log.PanicIf(err)
    }

    _, err = fmt.Fprintf(b, "\n")
    log.PanicIf(err)

    intermediateOutput := b.Bytes()

    sn.SetIntermediateOutput(intermediateOutput)

    for _, childNode := range sn.Children {
        err := md.RenderIntermediate(childNode)
        log.PanicIf(err)
    }

    return nil
}

func (md *MarkdownDialect) renderStatment(w io.Writer, ps sitebuilder.PageStatement) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    switch ps.Type {
    case sitebuilder.Heading:
        h := ps.StatementMetadata["heading"].(sitebuilder.HeadingWidget)

        err = HeadingToMarkdown(h, w)
        log.PanicIf(err)

    case sitebuilder.ContentImage:
        iw := ps.StatementMetadata["image"].(sitebuilder.ImageWidget)

        err = ImageWidgetToMarkdown(iw, w)
        log.PanicIf(err)

    case sitebuilder.HorizontalNavbar:
        nw := ps.StatementMetadata["horizontal_navbar"].(sitebuilder.NavbarWidget)

        err := InlineLinkListToMarkdown(nw.Items, w)
        log.PanicIf(err)

    case sitebuilder.VerticalNavbar:
        nw := ps.StatementMetadata["vertical_navbar"].(sitebuilder.NavbarWidget)

        err := BulletedLinkListToMarkdown(nw.Items, w)
        log.PanicIf(err)

    case sitebuilder.Link:
        lw := ps.StatementMetadata["link"].(sitebuilder.LinkWidget)

        err = LinkWidgetToMarkdown(lw, w)
        log.PanicIf(err)

    default:
        log.Panicf("widget not valid")
    }

    return nil
}

// RenderHtml produces HTML from the dialect-specific content.
func (md *MarkdownDialect) RenderHtml(sn *sitebuilder.SiteNode) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    intermediateOutput := sn.IntermediateOutput()
    output := blackfriday.Run(intermediateOutput)

    sn.SetFinalOutput(output)

    return nil
}
