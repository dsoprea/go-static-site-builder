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
    b := new(bytes.Buffer)

    _, err = fmt.Fprintf(b, "# %s\n\n", sn.PageTitle)
    log.PanicIf(err)

    for _, ps := range sn.Content.Statements {
        switch ps.Type {
        case sitebuilder.ContentImage:
            iw := ps.StatementMetadata["image"].(sitebuilder.ImageWidget)

            err = ImageWidgetToMarkdown(iw, b)
            log.PanicIf(err)

            err = md.writeDoubleNewline(b)
            log.PanicIf(err)

        case sitebuilder.Navbar:
            nw := ps.StatementMetadata["navbar"].(sitebuilder.NavbarWidget)

            for _, lw := range nw.Items {
                err = LinkWidgetToMarkdown(lw, b)
                log.PanicIf(err)

                _, err = b.Write([]byte{' '})
                log.PanicIf(err)
            }

            err = md.writeDoubleNewline(b)
            log.PanicIf(err)

        case sitebuilder.Link:
            lw := ps.StatementMetadata["link"].(sitebuilder.LinkWidget)

            err = LinkWidgetToMarkdown(lw, b)
            log.PanicIf(err)

            err = md.writeDoubleNewline(b)
            log.PanicIf(err)

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

func (md *MarkdownDialect) writeDoubleNewline(w io.Writer) (err error) {
    _, err = w.Write([]byte{'\n', '\n'})
    log.PanicIf(err)

    return nil
}

// RenderHtml produces HTML from the dialect-specific content.
func (md *MarkdownDialect) RenderHtml(sn *sitebuilder.SiteNode) (err error) {

    intermediateOutput := sn.IntermediateOutput()
    output := blackfriday.Run(intermediateOutput)

    sn.SetFinalOutput(output)

    return nil
}
