package markdowndialect

import (
    "fmt"
    "testing"

    "github.com/dsoprea/go-logging"

    "github.com/dsoprea/go-static-site-builder"
)

func TestMarkdownDialect_RenderIntermediate(t *testing.T) {
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

[image alt text](file://some/image/path)

`

    if string(actual) != expected {
        fmt.Printf("ACTUAL:\n%s", actual)

        t.Fatalf("Unexpected output.")
    }
}
