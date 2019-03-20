package sitebuilder

import (
    "github.com/dsoprea/go-logging"
)

type PageBuilder struct {
    sn *SiteNode
}

func NewPageBuilder(sn *SiteNode) *PageBuilder {
    return &PageBuilder{
        sn: sn,
    }
}

func (pb *PageBuilder) AddHeading(h HeadingWidget) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    metadata := map[string]interface{}{
        "heading": h,
    }

    ps := PageStatement{
        Type:              Heading,
        StatementMetadata: metadata,
    }

    pb.sn.Content.Add(ps)

    return nil
}

func (pb *PageBuilder) AddContentImage(iw ImageWidget) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

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

func (pb *PageBuilder) AddHorizontalNavbar(nw NavbarWidget) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    metadata := map[string]interface{}{
        "horizontal_navbar": nw,
    }

    ps := PageStatement{
        Type:              HorizontalNavbar,
        StatementMetadata: metadata,
    }

    pb.sn.Content.Add(ps)

    return nil
}

func (pb *PageBuilder) AddVerticalNavbar(nw NavbarWidget, text string) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    // Add heading.

    h := NewHeadingWidget(1, text)

    err = pb.AddHeading(h)
    log.PanicIf(err)

    // Add vertical navbar.

    metadata := map[string]interface{}{
        "vertical_navbar": nw,
    }

    ps := PageStatement{
        Type:              VerticalNavbar,
        StatementMetadata: metadata,
    }

    pb.sn.Content.Add(ps)

    return nil
}

func (pb *PageBuilder) AddLink(lw LinkWidget) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

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
