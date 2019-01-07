package markdowndialect

import (
    "bytes"
    "testing"

    "github.com/dsoprea/go-logging"

    "github.com/dsoprea/go-static-site-builder"
)

func TestImageWidgetToMarkdown(t *testing.T) {
    altText := "alt text"
    lrl := sitebuilder.NewLocalResourceLocator("some/image/path")

    iw := sitebuilder.NewImageWidget(altText, lrl, 0, 0)

    b := new(bytes.Buffer)

    err := ImageWidgetToMarkdown(iw, b)
    log.PanicIf(err)

    content := b.String()
    if content != "![alt text](file://some/image/path \"alt text\")" {
        t.Fatalf("Content not correct: [%s]", content)
    }
}

func TestLinkWidgetToMarkdown(t *testing.T) {
    text := "text"
    lrl := sitebuilder.NewLocalResourceLocator("some/file")

    lw := sitebuilder.NewLinkWidget(text, lrl)

    b := new(bytes.Buffer)

    err := LinkWidgetToMarkdown(lw, b)
    log.PanicIf(err)

    content := b.String()
    if content != "[text](file://some/file)" {
        t.Fatalf("Content not correct: [%s]", content)
    }
}
