package markdowndialect

import (
    "bytes"
    "fmt"
    "os"
    "path"
    "testing"

    "io/ioutil"

    "github.com/dsoprea/go-logging"

    "github.com/dsoprea/go-static-site-builder"
)

func TestMarkdownDialect_RenderIntermediate_Image(t *testing.T) {
    sc := sitebuilder.NewSiteContext("")
    md := NewMarkdownDialect()

    sb := sitebuilder.NewSiteBuilder("site title", md, sc)

    rootNode := sb.Root()
    pb := rootNode.Builder()

    lrl := sitebuilder.NewLocalResourceLocator("/some/image/path")
    iw := sitebuilder.NewImageWidget("image alt text", lrl, 0, 0)

    err := pb.AddContentImage(iw)
    log.PanicIf(err)

    err = md.RenderIntermediate(rootNode)
    log.PanicIf(err)

    actual := rootNode.IntermediateOutput()

    expected := `# site title

![image alt text](file:///some/image/path "image alt text")


`

    if string(actual) != expected {
        fmt.Printf("ACTUAL:\n%s", actual)

        t.Fatalf("Unexpected output.")
    }
}

// ExampleMarkdownDialect_RenderHtml is a wholistic usage example. It is named
// in such a way as to show up in the documentation.
func ExampleMarkdownDialect_RenderHtml() {
    tempPath, err := ioutil.TempDir("", "")
    log.PanicIf(err)

    defer os.RemoveAll(tempPath)

    sc := sitebuilder.NewSiteContext(tempPath)
    md := NewMarkdownDialect()

    sb := sitebuilder.NewSiteBuilder("Site Title", md, sc)

    rootNode := sb.Root()

    // Create content.

    rootPb := rootNode.Builder()

    lrl := sitebuilder.NewLocalResourceLocator("/some/image/path")

    iw := sitebuilder.NewImageWidget("image alt text 1", lrl, 0, 0)

    err = rootPb.AddContentImage(iw)
    log.PanicIf(err)

    childNode1, err := rootNode.AddChildNode("child1", "Child Page 1")
    log.PanicIf(err)

    childPb := childNode1.Builder()

    iw = sitebuilder.NewImageWidget("image alt text 2", lrl, 0, 0)

    err = childPb.AddContentImage(iw)
    log.PanicIf(err)

    childNode2, err := rootNode.AddChildNode("child2", "Child Page 2")
    log.PanicIf(err)

    childPb = childNode2.Builder()

    iw = sitebuilder.NewImageWidget("image alt text 3", lrl, 0, 0)

    err = childPb.AddContentImage(iw)
    log.PanicIf(err)

    childChildNode1, err := childNode1.AddChildNode("childChild1", "Child's Child Page 1")
    log.PanicIf(err)

    childPb = childChildNode1.Builder()

    iw = sitebuilder.NewImageWidget("image alt text 4", lrl, 0, 0)

    err = childPb.AddContentImage(iw)
    log.PanicIf(err)

    items := []sitebuilder.LinkWidget{
        sitebuilder.NewLinkWidget("Child1", sitebuilder.NewSitePageLocalResourceLocator(sb, "child1")),
        sitebuilder.NewLinkWidget("Child2", sitebuilder.NewSitePageLocalResourceLocator(sb, "child2")),
    }

    nw := sitebuilder.NewNavbarWidget(items)

    err = rootPb.AddHorizontalNavbar(nw)
    log.PanicIf(err)

    // Render and write.

    err = sb.WriteToPath()
    log.PanicIf(err)

    // Print.

    files, err := ioutil.ReadDir(tempPath)
    log.PanicIf(err)

    for _, fi := range files {
        filename := fi.Name()

        fmt.Printf("%s\n", filename)
        fmt.Printf("====================\n")
        fmt.Printf("\n")

        filepath := path.Join(tempPath, filename)
        content, err := ioutil.ReadFile(filepath)
        log.PanicIf(err)

        _, err = os.Stdout.Write([]byte(content))
        log.PanicIf(err)

        fmt.Printf("\n")
    }

    // Output:
    // child1.html
    // ====================
    //
    // <h1>Child Page 1</h1>
    //
    // <p><img src="file:///some/image/path" alt="image alt text 2" title="image alt text 2" /></p>
    //
    // child2.html
    // ====================
    //
    // <h1>Child Page 2</h1>
    //
    // <p><img src="file:///some/image/path" alt="image alt text 3" title="image alt text 3" /></p>
    //
    // childChild1.html
    // ====================
    //
    // <h1>Child&rsquo;s Child Page 1</h1>
    //
    // <p><img src="file:///some/image/path" alt="image alt text 4" title="image alt text 4" /></p>
    //
    // index.html
    // ====================
    //
    // <h1>Site Title</h1>
    //
    // <p><img src="file:///some/image/path" alt="image alt text 1" title="image alt text 1" /></p>
    //
    // <p><a href="child1.html">Child1</a> <a href="child2.html">Child2</a></p>
}

func TestMarkdownDialect_RenderIntermediate_HorizontalNavbar(t *testing.T) {
    sc := sitebuilder.NewSiteContext("")
    md := NewMarkdownDialect()

    sb := sitebuilder.NewSiteBuilder("site title", md, sc)
    rootNode := sb.Root()
    pb := rootNode.Builder()

    items := []sitebuilder.LinkWidget{
        sitebuilder.NewLinkWidget("Child1", sitebuilder.NewSitePageLocalResourceLocator(sb, "child1")),
        sitebuilder.NewLinkWidget("Child2", sitebuilder.NewSitePageLocalResourceLocator(sb, "child2")),
    }

    nw := sitebuilder.NewNavbarWidget(items)

    err := pb.AddHorizontalNavbar(nw)
    log.PanicIf(err)

    // The child nodes can be added after or before the navbar, or even later,
    // but the page-IDs must be valid by the time we render.

    _, err = rootNode.AddChildNode("child1", "Child Page 1")
    log.PanicIf(err)

    _, err = rootNode.AddChildNode("child2", "Child Page 2")
    log.PanicIf(err)

    err = md.RenderIntermediate(rootNode)
    log.PanicIf(err)

    actual := rootNode.IntermediateOutput()

    expected := "# site title\n\n[Child1](child1.html) [Child2](child2.html) \n\n\n"

    if string(actual) != expected {
        fmt.Printf("ACTUAL:\n=====\n%s=====\n", actual)
        fmt.Printf("EXPECTED:\n=====\n%s=====\n", expected)

        t.Fatalf("Unexpected output.")
    }
}

func TestMarkdownDialect_RenderIntermediate_Link(t *testing.T) {
    sc := sitebuilder.NewSiteContext("")
    md := NewMarkdownDialect()

    sb := sitebuilder.NewSiteBuilder("site title", md, sc)
    rootNode := sb.Root()
    pb := rootNode.Builder()

    splrl := sitebuilder.NewSitePageLocalResourceLocator(sb, "child1")
    lw := sitebuilder.NewLinkWidget("Child1", splrl)

    err := pb.AddLink(lw)
    log.PanicIf(err)

    _, err = rootNode.AddChildNode("child1", "Child Page 1")
    log.PanicIf(err)

    err = md.RenderIntermediate(rootNode)
    log.PanicIf(err)

    actual := rootNode.IntermediateOutput()

    expected := "# site title\n\n[Child1](child1.html)\n"

    if string(actual) != expected {
        fmt.Printf("ACTUAL:\n=====\n%s=====\n", actual)
        fmt.Printf("EXPECTED:\n=====\n%s=====\n", expected)

        t.Fatalf("Unexpected output.")
    }
}

func TestMarkdownDialect_renderStatement_Image(t *testing.T) {
    sc := sitebuilder.NewSiteContext("")
    md := NewMarkdownDialect()

    sb := sitebuilder.NewSiteBuilder("site title", md, sc)
    rootNode := sitebuilder.NewSiteNode(sb, "node_id", "node title")

    pb := rootNode.Builder()

    lrl := sitebuilder.NewLocalResourceLocator("/some/image/path")
    iw := sitebuilder.NewImageWidget("image alt text", lrl, 0, 0)

    err := pb.AddContentImage(iw)
    log.PanicIf(err)

    len_ := len(rootNode.Content.Statements)
    if len_ != 1 {
        t.Fatalf("Not exactly one statement: (%d)", len_)
    }

    ps := rootNode.Content.Statements[0]

    b := new(bytes.Buffer)

    err = md.renderStatment(b, ps)
    log.PanicIf(err)

    actual := b.String()
    expected := "![image alt text](file:///some/image/path \"image alt text\")\n\n"

    if actual != expected {
        t.Fatalf("Image not rendered to Markdown correct:\nACTUAL:\n%s\n\nEXPECTED:\n%s", actual, expected)
    }
}

func TestMarkdownDialect_renderStatement_HorizontalNavbar(t *testing.T) {
    sc := sitebuilder.NewSiteContext("")
    md := NewMarkdownDialect()

    sb := sitebuilder.NewSiteBuilder("site title", md, sc)
    rootNode := sitebuilder.NewSiteNode(sb, "node_id", "node title")

    pb := rootNode.Builder()

    items := []sitebuilder.LinkWidget{
        sitebuilder.NewLinkWidget("Child1", sitebuilder.NewLocalResourceLocator("/some/image/path1")),
        sitebuilder.NewLinkWidget("Child2", sitebuilder.NewLocalResourceLocator("/some/image/path2")),
    }

    nw := sitebuilder.NewNavbarWidget(items)

    err := pb.AddHorizontalNavbar(nw)
    log.PanicIf(err)

    len_ := len(rootNode.Content.Statements)
    if len_ != 1 {
        t.Fatalf("Not exactly one statement: (%d)", len_)
    }

    ps := rootNode.Content.Statements[0]

    b := new(bytes.Buffer)

    err = md.renderStatment(b, ps)
    log.PanicIf(err)

    actual := b.String()
    expected := "[Child1](file:///some/image/path1) [Child2](file:///some/image/path2) \n\n"

    if actual != expected {
        t.Fatalf("Horizontal navbar not rendered to Markdown correct:\nACTUAL:\n[%s]\n\nEXPECTED:\n[%s]", actual, expected)
    }
}

func TestMarkdownDialect_renderStatement_VerticalNavbar(t *testing.T) {
    sc := sitebuilder.NewSiteContext("")
    md := NewMarkdownDialect()

    sb := sitebuilder.NewSiteBuilder("site title", md, sc)
    rootNode := sitebuilder.NewSiteNode(sb, "node_id", "node title")

    pb := rootNode.Builder()

    items := []sitebuilder.LinkWidget{
        sitebuilder.NewLinkWidget("Child1", sitebuilder.NewLocalResourceLocator("/some/image/path1")),
        sitebuilder.NewLinkWidget("Child2", sitebuilder.NewLocalResourceLocator("/some/image/path2")),
    }

    nw := sitebuilder.NewNavbarWidget(items)

    err := pb.AddVerticalNavbar(nw, "test heading")
    log.PanicIf(err)

    b := new(bytes.Buffer)

    len_ := len(rootNode.Content.Statements)
    if len_ != 2 {
        t.Fatalf("Not exactly two statements: (%d)", len_)
    }

    ps1 := rootNode.Content.Statements[0]

    err = md.renderStatment(b, ps1)
    log.PanicIf(err)

    ps2 := rootNode.Content.Statements[1]

    err = md.renderStatment(b, ps2)
    log.PanicIf(err)

    actual := b.String()
    expected := `# test heading

- [Child1](file:///some/image/path1)
- [Child2](file:///some/image/path2)

`

    if actual != expected {
        t.Fatalf("Vertical navbar not rendered to Markdown correct:\nACTUAL:\n[%s]\n\nEXPECTED:\n[%s]", actual, expected)
    }
}

func TestMarkdownDialect_renderStatement_Link(t *testing.T) {
    sc := sitebuilder.NewSiteContext("")
    md := NewMarkdownDialect()

    sb := sitebuilder.NewSiteBuilder("site title", md, sc)
    rootNode := sitebuilder.NewSiteNode(sb, "node_id", "node title")

    pb := rootNode.Builder()

    lw := sitebuilder.NewLinkWidget("SomeLink", sitebuilder.NewLocalResourceLocator("/some/image/path"))

    err := pb.AddLink(lw)
    log.PanicIf(err)

    len_ := len(rootNode.Content.Statements)
    if len_ != 1 {
        t.Fatalf("Not exactly one statement: (%d)", len_)
    }

    ps := rootNode.Content.Statements[0]

    b := new(bytes.Buffer)

    err = md.renderStatment(b, ps)
    log.PanicIf(err)

    actual := b.String()
    expected := "[SomeLink](file:///some/image/path)"

    if actual != expected {
        t.Fatalf("Link not rendered to Markdown correct:\nACTUAL:\n[%s]\n\nEXPECTED:\n[%s]", actual, expected)
    }
}

func TestMarkdownDialect_renderStatement_Heading(t *testing.T) {
    sc := sitebuilder.NewSiteContext("")
    md := NewMarkdownDialect()

    sb := sitebuilder.NewSiteBuilder("site title", md, sc)
    rootNode := sitebuilder.NewSiteNode(sb, "node_id", "node title")

    pb := rootNode.Builder()

    hw := sitebuilder.NewHeadingWidget(1, "some heading")

    err := pb.AddHeading(hw)
    log.PanicIf(err)

    len_ := len(rootNode.Content.Statements)
    if len_ != 1 {
        t.Fatalf("Not exactly one statement: (%d)", len_)
    }

    ps := rootNode.Content.Statements[0]

    b := new(bytes.Buffer)

    err = md.renderStatment(b, ps)
    log.PanicIf(err)

    actual := b.String()
    expected := "# some heading\n\n"

    if actual != expected {
        t.Fatalf("Heading not rendered to Markdown correct:\nACTUAL:\n[%s]\n\nEXPECTED:\n[%s]", actual, expected)
    }
}

// TODO(dustin): Also, test that we get an error with the link widget if not a valid page-ID.
