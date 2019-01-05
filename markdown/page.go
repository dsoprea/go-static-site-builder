package markdowndialect

import (
    "bytes"
    "fmt"

    "github.com/dsoprea/go-logging"

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
    b := new(bytes.Buffer)

    _, err = fmt.Fprintf(b, "# %s\n\n", sn.PageTitle)
    log.PanicIf(err)

    for _, ps := range sn.Content.Statements {
        switch ps.Type {
        case sitebuilder.ContentImage:
            iw := ps.StatementMetadata["image"].(sitebuilder.ImageWidget)

            err = ImageWidgetToMarkdown(iw, b)
            log.PanicIf(err)

            _, err = b.Write([]byte("\n\n"))
            log.PanicIf(err)

            // TODO(dustin): !! Finish handling navbar.
        default:
            log.Panicf("widget not valid")
        }
    }

    intermediateOutput := b.Bytes()
    sn.SetIntermediateOutput(intermediateOutput)

    for _, childNode := range sn.Children {
        err := md.RenderIntermediate(childNode)
        log.PanicIf(err)
    }

    return nil
}

// RenderHtml produces HTML from the dialect-specific content.
func (md *MarkdownDialect) RenderHtml(sn *sitebuilder.SiteNode) (err error) {

    // TODO(dustin): !! Finish.

    return nil
}
