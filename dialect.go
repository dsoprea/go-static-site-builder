package sitebuilder

// DialectPageBuilder defines basic operations for a dialect that add statements
// and metadata to a PageDialectContent struct that it is already equipped with.
type DialectPageBuilder interface {
    // AddContentImage adds one large image, dead-center, at the next position
    // in the page.
    AddContentImage(altText string, locator ResourceLocator) (err error)

    // AddChildNavbar adds a navbar with links for all children.
    AddChildNavbar() (err error)
}

// Dialect defines high-level, dialect-specific translation operations.
type Dialect interface {
    // RenderIntermediate produces dialect-specific content that can be passed
    // to RenderHtml.
    RenderIntermediate(sn *SiteNode) (err error)

    // RenderHtml produces HTML from the dialect-specific content.
    RenderHtml(sn *SiteNode) (err error)
}
