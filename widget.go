package sitebuilder

import ()

type WidgetType int

// These should represent every method in the PageDialectBuilder interface.
const (
    ContentImage WidgetType = 1 + iota
    ChildrenNavbar
)

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

type NavbarItem struct {
    Name   string
    PageId string
}

type NavbarWidget struct {
    Items []NavbarItem
}

func NewNavbarWidget(items []NavbarItem) NavbarWidget {
    return NavbarWidget{
        Items: items,
    }
}
