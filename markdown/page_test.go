package markdowndialect

import (
    "fmt"
    "io/ioutil"
    "os"
    "path"
    "strings"
    "testing"

    "github.com/dsoprea/go-logging"

    "github.com/dsoprea/go-static-site-builder"
)

func TestMarkdownDialect_RenderIntermediate_Image(t *testing.T) {
    sc := sitebuilder.NewSiteContext("")
    md := NewMarkdownDialect()

    sb := sitebuilder.NewSiteBuilder("site title", md, sc)

    rootNode := sb.Root()
    pb := rootNode.Builder()

    lrl := sitebuilder.NewLocalResourceLocator("some/image/path")

    iw := sitebuilder.NewImageWidget("image alt text", lrl)

    err := pb.AddContentImage(iw)
    log.PanicIf(err)

    err = md.RenderIntermediate(rootNode)
    log.PanicIf(err)

    actual := rootNode.IntermediateOutput()

    expected := `# site title

![image alt text](file://some/image/path "image alt text")

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

    lrl := sitebuilder.NewLocalResourceLocator("some/image/path")

    iw := sitebuilder.NewImageWidget("image alt text 1", lrl)

    err = rootPb.AddContentImage(iw)
    log.PanicIf(err)

    childNode1, err := rootNode.AddChildNode("child1", "Child Page 1")
    log.PanicIf(err)

    childPb := childNode1.Builder()

    iw = sitebuilder.NewImageWidget("image alt text 2", lrl)

    err = childPb.AddContentImage(iw)
    log.PanicIf(err)

    childNode2, err := rootNode.AddChildNode("child2", "Child Page 2")
    log.PanicIf(err)

    childPb = childNode2.Builder()

    iw = sitebuilder.NewImageWidget("image alt text 3", lrl)

    err = childPb.AddContentImage(iw)
    log.PanicIf(err)

    childChildNode1, err := childNode1.AddChildNode("childChild1", "Child's Child Page 1")
    log.PanicIf(err)

    childPb = childChildNode1.Builder()

    iw = sitebuilder.NewImageWidget("image alt text 4", lrl)

    err = childPb.AddContentImage(iw)
    log.PanicIf(err)

    items := []sitebuilder.LinkWidget{
        sitebuilder.NewLinkWidget("Child1", sitebuilder.NewSitePageLocalResourceLocator(sb, "child1")),
        sitebuilder.NewLinkWidget("Child2", sitebuilder.NewSitePageLocalResourceLocator(sb, "child2")),
    }

    nw := sitebuilder.NewNavbarWidget(items)

    err = rootPb.AddNavbar(nw)
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

        // For the [testable] example.
        fixedContent := strings.Replace(string(content), tempPath, "example_path", -1)

        _, err = os.Stdout.Write([]byte(fixedContent))
        log.PanicIf(err)

        fmt.Printf("\n")
    }

    // Output:
    // child1.html
    // ====================
    //
    // <h1>Child Page 1</h1>
    //
    // <p><img src="file://some/image/path" alt="image alt text 2" title="image alt text 2" /></p>
    //
    // child2.html
    // ====================
    //
    // <h1>Child Page 2</h1>
    //
    // <p><img src="file://some/image/path" alt="image alt text 3" title="image alt text 3" /></p>
    //
    // childChild1.html
    // ====================
    //
    // <h1>Child&rsquo;s Child Page 1</h1>
    //
    // <p><img src="file://some/image/path" alt="image alt text 4" title="image alt text 4" /></p>
    //
    // index.html
    // ====================
    //
    // <h1>Site Title</h1>
    //
    // <p><img src="file://some/image/path" alt="image alt text 1" title="image alt text 1" /></p>
    //
    // <p><a href="file://example_path/child1.html">Child1</a> <a href="file://example_path/child2.html">Child2</a></p>
}

func TestMarkdownDialect_RenderIntermediate_Navbar(t *testing.T) {
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

    err := pb.AddNavbar(nw)
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

    expected := "# site title\n\n[Child1](file://child1.html) [Child2](file://child2.html) \n\n"

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

    expected := "# site title\n\n[Child1](file://child1.html)\n\n"

    if string(actual) != expected {
        fmt.Printf("ACTUAL:\n=====\n%s=====\n", actual)
        fmt.Printf("EXPECTED:\n=====\n%s=====\n", expected)

        t.Fatalf("Unexpected output.")
    }
}

// TODO(dustin): Also, test that we get an error with the link widget if not a valid page-ID.
