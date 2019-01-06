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

func (pb *PageBuilder) AddContentImage(iw ImageWidget) (err error) {
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

func (pb *PageBuilder) AddNavbar(nw NavbarWidget) (err error) {
    metadata := map[string]interface{}{
        "navbar": nw,
    }

    ps := PageStatement{
        Type:              Navbar,
        StatementMetadata: metadata,
    }

    pb.sn.Content.Add(ps)

    return nil
}

func (pb *PageBuilder) AddLink(lw LinkWidget) (err error) {
    metadata := map[string]interface{}{
        "link": lw,
    }

    ps := PageStatement{
        Type:              Link,
        StatementMetadata: metadata,
    }

    pb.sn.Content.Add(ps)

    return nil
}
