package markdowndialect

import (
    "bytes"
    "testing"

    "github.com/dsoprea/go-logging"

    "github.com/dsoprea/go-static-site-builder"
)

func TestImageWidgetToMarkdown(t *testing.T) {
    altText := "alt text"
    lrl := sitebuilder.NewLocalResourceLocator("/some/image/path")

    iw := sitebuilder.NewImageWidget(altText, lrl, 0, 0)

    b := new(bytes.Buffer)

    err := ImageWidgetToMarkdown(iw, b)
    log.PanicIf(err)

    content := b.String()
    if content != "![alt text](file:///some/image/path \"alt text\")\n\n" {
        t.Fatalf("Content not correct: [%s]", content)
    }
}

func TestLinkWidgetToMarkdown(t *testing.T) {
    text := "text"
    lrl := sitebuilder.NewLocalResourceLocator("/some/file")

    lw := sitebuilder.NewLinkWidget(text, lrl)

    b := new(bytes.Buffer)

    err := LinkWidgetToMarkdown(lw, b)
    log.PanicIf(err)

    content := b.String()
    if content != "[text](file:///some/file)" {
        t.Fatalf("Content not correct: [%s]", content)
    }
}

func TestInlineLinkListToMarkdown(t *testing.T) {
    items := []sitebuilder.LinkWidget{
        sitebuilder.NewLinkWidget("Child1", sitebuilder.NewLocalResourceLocator("/some/image/path1")),
        sitebuilder.NewLinkWidget("Child2", sitebuilder.NewLocalResourceLocator("/some/image/path2")),
    }

    b := new(bytes.Buffer)

    err := InlineLinkListToMarkdown(items, b)
    log.PanicIf(err)

    actual := b.String()
    expected := "[Child1](file:///some/image/path1) [Child2](file:///some/image/path2) \n\n"

    if actual != expected {
        t.Fatalf("Inline link-list to Markdown not correct:\nACTUAL:\n[%s]\n\nEXPECTED:\n[%s]", actual, expected)
    }
}

func TestBulletedLinkListToMarkdown(t *testing.T) {
    items := []sitebuilder.LinkWidget{
        sitebuilder.NewLinkWidget("Child1", sitebuilder.NewLocalResourceLocator("/some/image/path1")),
        sitebuilder.NewLinkWidget("Child2", sitebuilder.NewLocalResourceLocator("/some/image/path2")),
    }

    b := new(bytes.Buffer)

    err := BulletedLinkListToMarkdown(items, b)
    log.PanicIf(err)

    actual := b.String()
    expected := "- [Child1](file:///some/image/path1)\n- [Child2](file:///some/image/path2)\n\n"

    if actual != expected {
        t.Fatalf("Bulleted link-list to Markdown not correct:\nACTUAL:\n[%s]\n\nEXPECTED:\n[%s]", actual, expected)
    }
}

func TestHeadingToMarkdown(t *testing.T) {
    b := new(bytes.Buffer)

    hw := sitebuilder.NewHeadingWidget(1, "some heading")

    err := HeadingToMarkdown(hw, b)
    log.PanicIf(err)

    actual := b.String()
    expected := "# some heading\n\n"

    if actual != expected {
        t.Fatalf("Heading to Markdown not correct:\nACTUAL:\n[%s]\n\nEXPECTED:\n[%s]", actual, expected)
    }
}

func TestWriteNewline(t *testing.T) {
    b := new(bytes.Buffer)

    err := WriteNewline(b)
    log.PanicIf(err)

    actual := b.String()
    expected := "\n"

    if actual != expected {
        t.Fatalf("Single-newline to Markdown not correct:\nACTUAL:\n[%s]\n\nEXPECTED:\n[%s]", actual, expected)
    }
}

func TestWriteDoubleNewline(t *testing.T) {
    b := new(bytes.Buffer)

    err := WriteDoubleNewline(b)
    log.PanicIf(err)

    actual := b.String()
    expected := "\n\n"

    if actual != expected {
        t.Fatalf("Double-newline to Markdown not correct:\nACTUAL:\n[%s]\n\nEXPECTED:\n[%s]", actual, expected)
    }
}
