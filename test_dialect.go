package sitebuilder

import (
    "bytes"
    "fmt"
    "strings"

    "github.com/dsoprea/go-logging"
)

type TestDialect struct {
}

func NewTestDialect() (md *TestDialect) {
    return &TestDialect{}
}

// RenderIntermediate produces dialect-specific content that can be passed to
// RenderHtml.
func (md *TestDialect) RenderIntermediate(sn *SiteNode) (err error) {
    b := new(bytes.Buffer)

    _, err = fmt.Fprintf(b, "## page-top | %s ##\n", sn.PageTitle)
    log.PanicIf(err)

    for _, ps := range sn.Content.Statements {
        switch ps.Type {
        case ContentImage:
            iw := ps.StatementMetadata["image"].(ImageWidget)

            uri := iw.Locator.Uri()

            _, err := fmt.Fprintf(b, "## widget | image | %s | %s ##\n", iw.AltText, uri)
            log.PanicIf(err)
            // TODO(dustin): !! Finish handling navbar.
        default:
            log.Panicf("widget not valid")
        }
    }

    for _, childNode := range sn.Children {
        err := md.RenderIntermediate(childNode)
        log.PanicIf(err)
    }

    _, err = fmt.Fprintf(b, "## page-bottom | %s ##\n", sn.PageTitle)
    log.PanicIf(err)

    intermediateOutput := b.Bytes()
    sn.SetIntermediateOutput(intermediateOutput)

    return nil
}

// RenderHtml produces HTML from the dialect-specific content.
func (md *TestDialect) RenderHtml(sn *SiteNode) (err error) {
    s := string(sn.IntermediateOutput())
    lines := strings.Split(s, "\n")

    b := new(bytes.Buffer)

    for _, line := range lines {
        if strings.HasPrefix(line, "## page-top | ") == true && strings.HasSuffix(line, " ##") == true {
            // Example: ## page-title | node_title ##

            pivot := 14
            content := line[pivot : len(line)-3]

            _, err := fmt.Fprintf(b, "<header>%s</header>\n", content)
            log.PanicIf(err)
        } else if strings.HasPrefix(line, "## widget | image | ") == true && strings.HasSuffix(line, " ##") == true {
            // Example: ## widget | image | image alt text | file://some/image/path ##

            pivot := 20
            content := line[pivot : len(line)-3]

            _, err := fmt.Fprintf(b, "<widget>%s</widget>\n", content)
            log.PanicIf(err)
        } else if strings.HasPrefix(line, "## page-bottom | ") == true && strings.HasSuffix(line, " ##") == true {
            // Example: ## page-title | node_title ##

            pivot := 17
            content := line[pivot : len(line)-3]

            _, err := fmt.Fprintf(b, "<footer>%s</footer>\n", content)
            log.PanicIf(err)
        } else if len(line) == 0 {
            continue
        } else {
            log.Panicf("intermediate line not valid: [%s]", line)
        }
    }

    finalOutput := b.Bytes()
    sn.SetFinalOutput(finalOutput)

    return nil
}
