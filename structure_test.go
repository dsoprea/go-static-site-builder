package sitebuilder

import (
    "fmt"
    "io/ioutil"
    "os"
    "path"
    "reflect"
    "sort"
    "testing"

    "github.com/dsoprea/go-logging"
)

func TestPageContent_Add(t *testing.T) {
    metadata := map[string]interface{}{
        "abc": 123,
    }

    ps := PageStatement{
        Type:              ContentImage,
        StatementMetadata: metadata,
    }

    pc := NewPageContent()
    pc.Add(ps)

    if len(pc.Statements) != 1 {
        t.Fatalf("Not exactly one statement.")
    } else if reflect.DeepEqual(pc.Statements[0], ps) != true {
        t.Fatalf("Metadata not correct.")
    } else if len(pc.PageMetadata) != 0 {
        t.Fatalf("There shouldn't be any page metadata.")
    }
}

func TestSiteNode_Render(t *testing.T) {
    sc := NewSiteContext("")

    td := NewTestDialect()
    sb := NewSiteBuilder("node title", td, sc)

    // Generate content.

    rootNode := NewSiteNode(sb, "node_id", "node title")
    pb := rootNode.Builder()

    lrl := NewLocalResourceLocator("some/image/path")

    iw := NewImageWidget("image alt text", lrl, 0, 0)

    err := pb.AddContentImage(iw)
    log.PanicIf(err)

    // Get output.

    err = rootNode.Render()
    log.PanicIf(err)

    actual := rootNode.FinalOutput()

    expected := `<header>node title</header>
<widget>image alt text | file://some/image/path</widget>
<footer>node title</footer>
`

    if string(actual) != expected {
        fmt.Printf("ACTUAL:\n%s", actual)

        t.Fatalf("Unexpected output.")
    }
}

func TestSiteNode_AddChildNode(t *testing.T) {
    sc := NewSiteContext("")

    td := NewTestDialect()
    sb := NewSiteBuilder("node title", td, sc)

    // Generate content.

    rootNode := NewSiteNode(sb, "node_id", "node title")
    pb := rootNode.Builder()

    lrl := NewLocalResourceLocator("some/image/path")

    iw := NewImageWidget("image alt text", lrl, 0, 0)

    err := pb.AddContentImage(iw)
    log.PanicIf(err)

    rootNode.AddChildNode("child1", "child title1")
    rootNode.AddChildNode("child2", "child title2")

    // Get output.

    err = rootNode.Render()
    log.PanicIf(err)

    actual := rootNode.FinalOutput()

    expected := `<header>node title</header>
<widget>image alt text | file://some/image/path</widget>
<footer>node title</footer>
`

    if string(actual) != expected {
        fmt.Printf("ACTUAL:\n%s", actual)

        t.Fatalf("Unexpected output.")
    }
}

func TestSiteBuilder_WriteToPath(t *testing.T) {
    tempPath, err := ioutil.TempDir("", "gssb")
    log.PanicIf(err)

    defer os.RemoveAll(tempPath)

    sc := NewSiteContext(tempPath)

    td := NewTestDialect()
    sb := NewSiteBuilder("site title", td, sc)

    rootNode := sb.Root()
    pb := rootNode.Builder()

    lrl := NewLocalResourceLocator("some/image/path")

    iw := NewImageWidget("image alt text", lrl, 0, 0)

    err = pb.AddContentImage(iw)
    log.PanicIf(err)

    // Write.

    err = sb.WriteToPath()
    log.PanicIf(err)

    // Read.

    files, err := ioutil.ReadDir(tempPath)
    log.PanicIf(err)

    if len(files) != 1 {
        for _, fi := range files {
            fmt.Printf("FILE: [%s]\n", fi.Name())
        }

        t.Fatalf("Exact files weren't produced.")
    }

    indexFilepath := path.Join(tempPath, "index.html")
    actualBytes, err := ioutil.ReadFile(indexFilepath)
    log.PanicIf(err)

    actual := string(actualBytes)

    expected := `<header>site title</header>
<widget>image alt text | file://some/image/path</widget>
<footer>site title</footer>
`

    if actual != expected {
        fmt.Printf("ACTUAL:\n%s", actual)

        t.Fatalf("Unexpected output.")
    }
}

func TestSiteBuilder_writeToPath_Simple(t *testing.T) {
    tempPath, err := ioutil.TempDir("", "gssb")
    log.PanicIf(err)

    defer os.RemoveAll(tempPath)

    sc := NewSiteContext(tempPath)

    td := NewTestDialect()
    sb := NewSiteBuilder("site title", td, sc)

    rootNode := sb.Root()
    pb := rootNode.Builder()

    lrl := NewLocalResourceLocator("some/image/path")

    iw := NewImageWidget("image alt text", lrl, 0, 0)

    err = pb.AddContentImage(iw)
    log.PanicIf(err)

    // Write.

    err = sb.rootNode.Render()
    log.PanicIf(err)

    err = sb.writeToPath(rootNode)
    log.PanicIf(err)

    // Read.

    files, err := ioutil.ReadDir(tempPath)
    log.PanicIf(err)

    if len(files) != 1 {
        for _, fi := range files {
            fmt.Printf("FILE: [%s]\n", fi.Name())
        }

        t.Fatalf("Exact files weren't produced.")
    }

    indexFilepath := path.Join(tempPath, "index.html")
    actualBytes, err := ioutil.ReadFile(indexFilepath)
    log.PanicIf(err)

    actual := string(actualBytes)

    expected := `<header>site title</header>
<widget>image alt text | file://some/image/path</widget>
<footer>site title</footer>
`

    if actual != expected {
        fmt.Printf("ACTUAL:\n%s", actual)

        t.Fatalf("Unexpected output.")
    }
}

func TestSiteBuilder_writeToPath_Tree(t *testing.T) {
    tempPath, err := ioutil.TempDir("", "gssb")
    log.PanicIf(err)

    defer os.RemoveAll(tempPath)

    sc := NewSiteContext(tempPath)

    td := NewTestDialect()
    sb := NewSiteBuilder("site title", td, sc)

    rootNode := sb.Root()
    pb := rootNode.Builder()

    lrl := NewLocalResourceLocator("some/image/path")

    iw := NewImageWidget("image alt text", lrl, 0, 0)

    err = pb.AddContentImage(iw)
    log.PanicIf(err)

    childNode1, err := rootNode.AddChildNode("child1", "Child1")
    log.PanicIf(err)

    _, err = rootNode.AddChildNode("child2", "Child2")
    log.PanicIf(err)

    _, err = childNode1.AddChildNode("childChild1", "ChildChild1")
    log.PanicIf(err)

    // Write.

    err = sb.rootNode.Render()
    log.PanicIf(err)

    err = sb.writeToPath(rootNode)
    log.PanicIf(err)

    // Read.

    files, err := ioutil.ReadDir(tempPath)
    log.PanicIf(err)

    expectedFiles := []string{
        "index.html",
        "child1.html",
        "child2.html",
        "childChild1.html",
    }

    actualFiles := make([]string, 0)
    for _, fi := range files {
        actualFiles = append(actualFiles, fi.Name())
    }

    ss := sort.StringSlice(expectedFiles)
    ss.Sort()

    ss = sort.StringSlice(actualFiles)
    ss.Sort()

    if reflect.DeepEqual(actualFiles, expectedFiles) == false {
        t.Fatalf("Exact files weren't produced: %v", actualFiles)
    }
}

func TestSiteNode_AddChildNode_InvalidFormat(t *testing.T) {
    sc := NewSiteContext("")

    td := NewTestDialect()
    sb := NewSiteBuilder("node title", td, sc)

    // Generate content.

    _, err := sb.Root().AddChildNode("invalid child id", "child title1")
    if err == nil || err.Error() != "page-ID has an invalid format" {
        t.Fatalf("page-ID has an invalid format")
    }
}
