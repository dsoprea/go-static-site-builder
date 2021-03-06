[![Build Status](https://travis-ci.org/dsoprea/go-static-site-builder.svg?branch=master)](https://travis-ci.org/dsoprea/go-static-site-builder)
[![Coverage Status](https://coveralls.io/repos/github/dsoprea/go-static-site-builder/badge.svg?branch=master)](https://coveralls.io/github/dsoprea/go-static-site-builder?branch=master)
[![GoDoc](https://godoc.org/github.com/dsoprea/go-static-site-builder?status.svg)](https://godoc.org/github.com/dsoprea/go-static-site-builder)


# Overview

This supports building a static website directly via Go (not via the command-line).

This project was created in order to solve the problem of producing an HTML-based browser on-the-fly to accompany other data.


# Features

- Provides a simple builder type to populate widgets into nodes.
- Expresses website content as a general, hierarchical node structure.
- Website structure is serializable and therefore storable so that it can be stored, recalled, modified, and rerendered later.
- When the website is rendered, it is first rendered as intermediate content and then rendered as HTML content. This allows us to focus on producing lightweight markup while being able to offload the actual HTML production to a third-party tool that specializes in that. This consequently enables you to debug content issues in the HTML by inspecting the intermediate content.
- The intermediate content supports multiple dialects. This project comes with a [Markdown](https://daringfireball.net/projects/markdown) dialect.
- Images can be embedded directly into the HTML content.


# Example

Aside from a functioning command example at [example/helloworld](https://godoc.org/github.com/dsoprea/go-static-site-builder/example/helloworld), this is from the example at [MarkdownDialect.RenderHtml](https://godoc.org/github.com/dsoprea/go-static-site-builder/markdown#example-MarkdownDialect-RenderHtml):

```go
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
```

Output:

```
child1.html
====================

<h1>Child Page 1</h1>

<p><img src="file://some/image/path" alt="image alt text 2" title="image alt text 2" /></p>

child2.html
====================

<h1>Child Page 2</h1>

<p><img src="file://some/image/path" alt="image alt text 3" title="image alt text 3" /></p>

childChild1.html
====================

<h1>Child&rsquo;s Child Page 1</h1>

<p><img src="file://some/image/path" alt="image alt text 4" title="image alt text 4" /></p>

index.html
====================

<h1>Site Title</h1>

<p><img src="file://some/image/path" alt="image alt text 1" title="image alt text 1" /></p>

<p><a href="child1.html">Child1</a> <a href="child2.html">Child2</a></p>
```


# To Dos

- Add support for additional widgets. **This is low-cost but currently depends upon need. Contributions welcome.**
- Need to add actual support for save/restore of hierarchical structure for long-term maintenance and modification.
