package sitebuilder

import ()

type WidgetType int

// These should represent every method in the PageDialectBuilder interface.
const (
    ContentImage WidgetType = 1 + iota
    Navbar
    Link
)

// Image

type ImageWidget struct {
    AltText       string
    Locator       ResourceLocator
    Width, Height int
}

// NewImageWidget creates an image widget. If width and height are zero, no
// width/height will be prescribed. If just the width is given, the height will
// be scaled accordingly.
func NewImageWidget(altText string, locator ResourceLocator, width, height int) (iw ImageWidget) {
    return ImageWidget{
        AltText: altText,
        Locator: locator,
        Height:  height,
        Width:   width,
    }
}

// Link

type LinkWidget struct {
    Text    string
    Locator ResourceLocator
}

func NewLinkWidget(text string, locator ResourceLocator) (lw LinkWidget) {
    return LinkWidget{
        Text:    text,
        Locator: locator,
    }
}

// Navbar

type NavbarWidget struct {
    Items []LinkWidget
}

func NewNavbarWidget(items []LinkWidget) NavbarWidget {
    return NavbarWidget{
        Items: items,
    }
}
