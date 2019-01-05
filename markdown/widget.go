package markdowndialect

import (
    "fmt"
    "io"

    "github.com/dsoprea/go-logging"

    "github.com/dsoprea/go-static-site-builder"
)

func ImageWidgetToMarkdown(iw sitebuilder.ImageWidget, w io.Writer) (err error) {
    uri := iw.Locator.Uri()

    _, err = fmt.Fprintf(w, "![%s](%s \"%s\")", iw.AltText, uri, iw.AltText)
    log.PanicIf(err)

    return nil
}
