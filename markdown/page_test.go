package markdowndialect

import (
    "fmt"
    "io/ioutil"
    "os"
    "path"
    "testing"

    "github.com/dsoprea/go-logging"

    "github.com/dsoprea/go-static-site-builder"
)

func TestMarkdownDialect_RenderIntermediate_Image(t *testing.T) {
    md := NewMarkdownDialect()

    sb := sitebuilder.NewSiteBuilder("site title", md)

    rootNode := sb.Root()
    pb := rootNode.Builder()

    lrl := sitebuilder.NewLocalResourceLocator("some/image/path")

    err := pb.AddContentImage("image alt text", lrl)
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
    md := NewMarkdownDialect()
    sb := sitebuilder.NewSiteBuilder("Site Title", md)

    rootNode := sb.Root()

    // Create content.

    rootPb := rootNode.Builder()

    lrl := sitebuilder.NewLocalResourceLocator("some/image/path")

    err := rootPb.AddContentImage("image alt text 1", lrl)
    log.PanicIf(err)

    childNode1, err := rootNode.AddChild("child1", "Child Page 1")
    log.PanicIf(err)

    childPb := childNode1.Builder()

    err = childPb.AddContentImage("image alt text 2", lrl)
    log.PanicIf(err)

    childNode2, err := rootNode.AddChild("child2", "Child Page 2")
    log.PanicIf(err)

    childPb = childNode2.Builder()

    err = childPb.AddContentImage("image alt text 3", lrl)
    log.PanicIf(err)

    childChildNode1, err := childNode1.AddChild("childChild1", "Child's Child Page 1")
    log.PanicIf(err)

    childPb = childChildNode1.Builder()

    err = childPb.AddContentImage("image alt text 4", lrl)
    log.PanicIf(err)

    items := []sitebuilder.NavbarItem{
        {
            Name:   "Child1",
            PageId: "child1",
        },
        {
            Name:   "Child2",
            PageId: "child2",
        },
    }

    err = rootPb.AddChildrenNavbar(items)
    log.PanicIf(err)

    // Render and write.

    tempPath, err := ioutil.TempDir("", "")
    log.PanicIf(err)

    defer os.RemoveAll(tempPath)

    err = sb.WriteToPath(tempPath)
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

        _, err = os.Stdout.Write(content)
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
    // <p><a href="child1.html">Child1</a> <a href="child2.html">Child2</a></p>
}

func TestMarkdownDialect_RenderIntermediate_ChildrenNavbar(t *testing.T) {
    md := NewMarkdownDialect()

    sb := sitebuilder.NewSiteBuilder("site title", md)
    rootNode := sb.Root()
    pb := rootNode.Builder()

    items := []sitebuilder.NavbarItem{
        {
            Name:   "Child1",
            PageId: "child1",
        },
        {
            Name:   "Child2",
            PageId: "child2",
        },
    }

    err := pb.AddChildrenNavbar(items)
    log.PanicIf(err)

    // The child nodes can be added after or before the navbar, or even later,
    // but the page-IDs must be valid by the time we render.

    _, err = rootNode.AddChild("child1", "Child Page 1")
    log.PanicIf(err)

    _, err = rootNode.AddChild("child2", "Child Page 2")
    log.PanicIf(err)

    err = md.RenderIntermediate(rootNode)
    log.PanicIf(err)

    actual := rootNode.IntermediateOutput()

    expected := "# site title\n\n[Child1](child1.html) [Child2](child2.html) \n\n"

    if string(actual) != expected {
        fmt.Printf("ACTUAL:\n=====\n%s=====\n", actual)
        fmt.Printf("EXPECTED:\n=====\n%s=====\n", expected)

        t.Fatalf("Unexpected output.")
    }
}
