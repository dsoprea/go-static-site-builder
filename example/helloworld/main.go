package main

import (
    "os"
    "path"

    "github.com/dsoprea/go-logging"
    "github.com/jessevdk/go-flags"

    "github.com/dsoprea/go-static-site-builder"
    "github.com/dsoprea/go-static-site-builder/markdown"
)

type rootParameters struct {
    OutputPath string `long:"output-path" description:"Path to write to" required:"true"`
}

var (
    rootArguments = new(rootParameters)
)

func main() {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.PrintError(err)
            os.Exit(-1)
        }
    }()

    p := flags.NewParser(rootArguments, flags.Default)

    _, err := p.Parse()
    if err != nil {
        os.Exit(1)
    }

    sc := sitebuilder.NewSiteContext(rootArguments.OutputPath)
    md := markdowndialect.NewMarkdownDialect()

    sb := sitebuilder.NewSiteBuilder("Site Title", md, sc)

    // Create content on root page.

    rootNode := sb.Root()
    rootPb := rootNode.Builder()

    erl1, err := sitebuilder.NewEmbeddedResourceLocator(path.Join("asset", "image1.jpg"), "", false)
    log.PanicIf(err)

    iw := sitebuilder.NewImageWidget("image alt text 1", erl1, 100, 100)

    err = rootPb.AddContentImage(iw)
    log.PanicIf(err)

    // Add a new page.

    childNode1, err := rootNode.AddChildNode("child1", "Child Page 1")
    log.PanicIf(err)

    childPb := childNode1.Builder()

    erl2, err := sitebuilder.NewEmbeddedResourceLocator(path.Join("asset", "image2.jpg"), "", false)
    log.PanicIf(err)

    iw = sitebuilder.NewImageWidget("image alt text 2", erl2, 100, 100)

    err = childPb.AddContentImage(iw)
    log.PanicIf(err)

    // Add a new page.

    childNode2, err := rootNode.AddChildNode("child2", "Child Page 2")
    log.PanicIf(err)

    childPb = childNode2.Builder()

    erl3, err := sitebuilder.NewEmbeddedResourceLocator(path.Join("asset", "image3.jpg"), "", false)
    log.PanicIf(err)

    iw = sitebuilder.NewImageWidget("image alt text 3", erl3, 100, 100)

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

    // // Print.

    // files, err := ioutil.ReadDir(tempPath)
    // log.PanicIf(err)

    // for _, fi := range files {
    //     filename := fi.Name()

    //     fmt.Printf("%s\n", filename)
    //     fmt.Printf("====================\n")
    //     fmt.Printf("\n")

    //     filepath := path.Join(tempPath, filename)
    //     content, err := ioutil.ReadFile(filepath)
    //     log.PanicIf(err)

    //     // For the [testable] example.
    //     fixedContent := strings.Replace(string(content), tempPath, "example_path", -1)

    //     _, err = os.Stdout.Write([]byte(fixedContent))
    //     log.PanicIf(err)

    //     fmt.Printf("\n")
    // }
}
