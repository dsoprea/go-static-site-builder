package sitebuilder

import ()

type PageBuilder struct {
    sn *SiteNode
}

func NewPageBuilder(sn *SiteNode) *PageBuilder {
    return &PageBuilder{
        sn: sn,
    }
}

func (pb *PageBuilder) AddContentImage(altText string, locator ResourceLocator) (err error) {
    iw := NewImageWidget(altText, locator)

    metadata := map[string]interface{}{
        "image": iw,
    }

    ps := PageStatement{
        Type:              ContentImage,
        StatementMetadata: metadata,
    }

    pb.sn.Content.Add(ps)

    return nil
}

func (pb *PageBuilder) AddChildrenNavbar(items []NavbarItem) (err error) {
    nw := NewNavbarWidget(items)

    metadata := map[string]interface{}{
        "children_navbar": nw,
    }

    ps := PageStatement{
        Type:              ChildrenNavbar,
        StatementMetadata: metadata,
    }

    pb.sn.Content.Add(ps)

    return nil
}
