package markdowndialect

import (
    "fmt"
    "io"

    "github.com/dsoprea/go-logging"

    "github.com/dsoprea/go-static-site-builder"
)

func ImageWidgetToMarkdown(iw sitebuilder.ImageWidget, w io.Writer) (err error) {
    uri := iw.Locator.Uri()

    if iw.Width != 0 || iw.Height != 0 {
        if iw.Width != 0 && iw.Height == 0 {
            _, err = fmt.Fprintf(w, `<img src="%s" width="%d" alt="%s" />`, uri, iw.Width, iw.AltText)
            log.PanicIf(err)
        } else if iw.Width == 0 && iw.Height != 0 {
            _, err = fmt.Fprintf(w, `<img src="%s" height="%d" alt="%s" />`, uri, iw.Height, iw.AltText)
            log.PanicIf(err)
        } else if iw.Width != 0 && iw.Height != 0 {
            _, err = fmt.Fprintf(w, `<img src="%s" width="%d" height="%d" alt="%s" />`, uri, iw.Width, iw.Height, iw.AltText)
            log.PanicIf(err)
        }

        _, err = fmt.Fprintf(w, "<br /><br />\n\n")
        log.PanicIf(err)
    } else {
        _, err = fmt.Fprintf(w, "![%s](%s \"%s\")", iw.AltText, uri, iw.AltText)
        log.PanicIf(err)
    }

    return nil
}

func LinkWidgetToMarkdown(lw sitebuilder.LinkWidget, w io.Writer) (err error) {
    uri := lw.Locator.Uri()

    _, err = fmt.Fprintf(w, "[%s](%s)", lw.Text, uri)
    log.PanicIf(err)

    return nil
}
