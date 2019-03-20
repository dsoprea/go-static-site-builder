package markdowndialect

import (
    "fmt"
    "io"
    "strings"

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
        _, err = fmt.Fprintf(w, "![%s](%s \"%s\")\n\n", iw.AltText, uri, iw.AltText)
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

func HeadingToMarkdown(h sitebuilder.HeadingWidget, w io.Writer) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    prefix := strings.Repeat("#", h.Level)

    _, err = fmt.Fprintf(w, "%s %s\n\n", prefix, h.Text)
    log.PanicIf(err)

    return nil
}

func WriteNewline(w io.Writer) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    _, err = w.Write([]byte{'\n'})
    log.PanicIf(err)

    return nil
}

func WriteDoubleNewline(w io.Writer) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    _, err = w.Write([]byte{'\n', '\n'})
    log.PanicIf(err)

    return nil
}

func InlineLinkListToMarkdown(items []sitebuilder.LinkWidget, w io.Writer) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    for _, lw := range items {
        err = LinkWidgetToMarkdown(lw, w)
        log.PanicIf(err)

        _, err = w.Write([]byte{' '})
        log.PanicIf(err)
    }

    err = WriteDoubleNewline(w)
    log.PanicIf(err)

    return nil
}

func BulletedLinkListToMarkdown(items []sitebuilder.LinkWidget, w io.Writer) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    for _, lw := range items {
        _, err = w.Write([]byte("- "))
        log.PanicIf(err)

        err = LinkWidgetToMarkdown(lw, w)
        log.PanicIf(err)

        err = WriteNewline(w)
        log.PanicIf(err)
    }

    err = WriteNewline(w)
    log.PanicIf(err)

    return nil
}
