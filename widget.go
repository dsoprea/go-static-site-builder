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
    AltText string
    Locator ResourceLocator
}

func NewImageWidget(altText string, locator ResourceLocator) (iw ImageWidget) {
    return ImageWidget{
        AltText: altText,
        Locator: locator,
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
