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
    td := NewTestDialect()
    sb := NewSiteBuilder("node title", td)

    // Generate content.

    rootNode := NewSiteNode(sb, "node_id", "node title")
    pb := rootNode.Builder()

    lrl := NewLocalResourceLocator("some/image/path")

    err := pb.AddContentImage("image alt text", lrl)
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

func TestSiteNode_AddChild(t *testing.T) {
    td := NewTestDialect()
    sb := NewSiteBuilder("node title", td)

    // Generate content.

    rootNode := NewSiteNode(sb, "node_id", "node title")
    pb := rootNode.Builder()

    lrl := NewLocalResourceLocator("some/image/path")

    err := pb.AddContentImage("image alt text", lrl)
    log.PanicIf(err)

    rootNode.AddChild("child1", "child title1")
    rootNode.AddChild("child2", "child title2")

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
    td := NewTestDialect()
    sb := NewSiteBuilder("site title", td)

    rootNode := sb.Root()
    pb := rootNode.Builder()

    lrl := NewLocalResourceLocator("some/image/path")

    err := pb.AddContentImage("image alt text", lrl)
    log.PanicIf(err)

    // Write.

    tempPath, err := ioutil.TempDir("", "gssb")
    log.PanicIf(err)

    defer os.RemoveAll(tempPath)

    err = sb.WriteToPath(tempPath)
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
    td := NewTestDialect()
    sb := NewSiteBuilder("site title", td)

    rootNode := sb.Root()
    pb := rootNode.Builder()

    lrl := NewLocalResourceLocator("some/image/path")

    err := pb.AddContentImage("image alt text", lrl)
    log.PanicIf(err)

    // Write.

    err = sb.rootNode.Render()
    log.PanicIf(err)

    tempPath, err := ioutil.TempDir("", "gssb")
    log.PanicIf(err)

    defer os.RemoveAll(tempPath)

    err = sb.writeToPath(rootNode, tempPath)
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
    td := NewTestDialect()
    sb := NewSiteBuilder("site title", td)

    rootNode := sb.Root()
    pb := rootNode.Builder()

    lrl := NewLocalResourceLocator("some/image/path")

    err := pb.AddContentImage("image alt text", lrl)
    log.PanicIf(err)

    childNode1 := rootNode.AddChild("child1", "Child1")
    rootNode.AddChild("child2", "Child2")
    childNode1.AddChild("childChild1", "ChildChild1")

    // Write.

    err = sb.rootNode.Render()
    log.PanicIf(err)

    tempPath, err := ioutil.TempDir("", "gssb")
    log.PanicIf(err)

    defer os.RemoveAll(tempPath)

    err = sb.writeToPath(rootNode, tempPath)
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
