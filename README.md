[![Build Status](https://travis-ci.org/dsoprea/go-static-site-builder.svg?branch=master)](https://travis-ci.org/dsoprea/go-static-site-builder)
[![Coverage Status](https://coveralls.io/repos/github/dsoprea/go-static-site-builder/badge.svg?branch=master)](https://coveralls.io/github/dsoprea/go-static-site-builder?branch=master)
[![GoDoc](https://godoc.org/github.com/dsoprea/go-static-site-builder?status.svg)](https://godoc.org/github.com/dsoprea/go-static-site-builder/index)


# Overview

This supports building a static website directly via Go (not via the command-line).

This project was created in order to solve the problem of producing an HTML-based browser on-the-fly to accompany other data.


# Features

- Provides a simple builder tool to populate widgets.
- Builds hierarchical website content as a general structure.
- Website structure is serializable and therefore storable so that tit can be modified and rerendered later.
- When the website is rendered, it is first rendered as intermediate content and then rendered and HTML content. This allows us to produce lightweight markup but offload the actual HTML production to a third-party tool. This also enables you to debug content issues in the HTML by inspecting the intermediate content.
- The intermediate content supports multiple dialects. This project comes with a [Markdown](https://daringfireball.net/projects/markdown) dialect.


# Example

Example from [github.com/dsoprea/go-static-site-builder/markdown](https://github.com/dsoprea/go-static-site-builder/tree/master/markdown):

```go
md := NewMarkdownDialect()
sb := sitebuilder.NewSiteBuilder("Site Title", md)

rootNode := sb.Root()

// Create content.

pb := rootNode.Builder()

lrl := sitebuilder.NewLocalResourceLocator("some/image/path")

err := pb.AddContentImage("image alt text 1", lrl)
log.PanicIf(err)

childNode1, err := rootNode.AddChild("child1", "Child Page 1")
log.PanicIf(err)

pb = childNode1.Builder()

err = pb.AddContentImage("image alt text 2", lrl)
log.PanicIf(err)

childNode2, err := rootNode.AddChild("child2", "Child Page 2")
log.PanicIf(err)

pb = childNode2.Builder()

err = pb.AddContentImage("image alt text 3", lrl)
log.PanicIf(err)

childChildNode1, err := childNode1.AddChild("childChild1", "Child's Child Page 1")
log.PanicIf(err)

pb = childChildNode1.Builder()

err = pb.AddContentImage("image alt text 4", lrl)
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
// <p><a href="file://some/image/path">image alt text 2</a></p>
//
// child2.html
// ====================
//
// <h1>Child Page 2</h1>
//
// <p><a href="file://some/image/path">image alt text 3</a></p>
//
// childChild1.html
// ====================
//
// <h1>Child&rsquo;s Child Page 1</h1>
//
// <p><a href="file://some/image/path">image alt text 4</a></p>
//
// index.html
// ====================
//
// <h1>Site Title</h1>
//
// <p><a href="file://some/image/path">image alt text 1</a></p>
```


# To Dos

- Add support for additional widgets. **This is low-cost but currently depends upon need. Contributions welcome.**
- Need to add actual support for save/restore of hierarchical structure for long-term maintenance and modification.
