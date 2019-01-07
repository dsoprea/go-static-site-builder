package sitebuilder

import (
    "encoding/base64"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "mime"
    "os"
    "path"
    "path/filepath"

    "github.com/dsoprea/go-logging"
)

const (
    // MaxEmbeddedResourceSize is the maximum allowed size of any embedded
    // resource.
    MaxEmbeddedResourceSize = 1024 * 1024 * 20
)

var (
    ErrEmbeddedResourceTooLarge = errors.New("embedded resource will be too large")
)

// A local file-path.

type LocalResourceLocator struct {
    LocalFilepath string
}

func NewLocalResourceLocator(localFilepath string) (lrl *LocalResourceLocator) {
    return &LocalResourceLocator{
        LocalFilepath: localFilepath,
    }
}

func (lrl *LocalResourceLocator) Uri() string {
    return fmt.Sprintf("file://%s", lrl.LocalFilepath)
}

// A locator that points to the final output page for a node.

type SitePageLocalResourceLocator struct {
    sb     *SiteBuilder
    PageId string
}

func NewSitePageLocalResourceLocator(sb *SiteBuilder, pageId string) (splrl *SitePageLocalResourceLocator) {
    return &SitePageLocalResourceLocator{
        sb:     sb,
        PageId: pageId,
    }
}

func (splrl *SitePageLocalResourceLocator) Uri() string {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.Panic(err)
        }
    }()

    if found := splrl.sb.PageIsValid(splrl.PageId); found == false {
        log.Panicf("resource refers to invalid page-ID [%s]", splrl.PageId)
    }

    outputPath := splrl.sb.Context().HtmlOutputPath()
    filename := splrl.sb.Context().GetFinalPageFilename(splrl.PageId)
    filepath := path.Join(outputPath, filename)

    return fmt.Sprintf("file://%s", filepath)
}

// Embedded data (rather than any local or remote references).

type EmbeddedResourceLocator struct {
    MimeType          string
    Base64EncodedData string
    Filepath          string
}

func NewEmbeddedResourceLocator(mimeType string, raw []byte) (erl *EmbeddedResourceLocator, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
            log.Panic(err)
        }
    }()

    if len(raw) > MaxEmbeddedResourceSize {
        log.Panic(ErrEmbeddedResourceTooLarge)
    }

    encoded := base64.StdEncoding.EncodeToString(raw)

    erl = &EmbeddedResourceLocator{
        MimeType:          mimeType,
        Base64EncodedData: encoded,
    }

    return erl, nil
}

func NewEmbeddedResourceLocatorWithReader(mimeType string, r io.Reader) (erl *EmbeddedResourceLocator, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
            log.Panic(err)
        }
    }()

    // We read one more byte tha we allow just so we know for certain whether we
    // should throw an error for it being too large.
    lr := io.LimitReader(r, MaxEmbeddedResourceSize+1)

    raw, err := ioutil.ReadAll(lr)
    log.PanicIf(err)

    if len(raw) > MaxEmbeddedResourceSize {
        log.Panic(ErrEmbeddedResourceTooLarge)
    }

    encoded := base64.StdEncoding.EncodeToString(raw)

    erl = &EmbeddedResourceLocator{
        MimeType:          mimeType,
        Base64EncodedData: encoded,
    }

    return erl, nil
}

// NewEmbeddedResourceLocatorWithFilepath will read the given file and then
// embed it. If `mimeType` is an empty-string, we will detect it based on the
// extension. If `readImmediately` is `false`, we'll read and embed it
// immediately rather than defer until the URI is actually requested.
func NewEmbeddedResourceLocatorWithFilepath(localFilepath, mimeType string, readImmediately bool) (erl *EmbeddedResourceLocator, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
            log.Panic(err)
        }
    }()

    fi, err := os.Stat(localFilepath)
    log.PanicIf(err)

    if fi.Size() > int64(MaxEmbeddedResourceSize) {
        log.Panic(ErrEmbeddedResourceTooLarge)
    }

    if mimeType == "" {
        ext := filepath.Ext(localFilepath)
        if ext == "" {
            log.Panicf("no mime-type was given but file does not have an extension: [%s]", localFilepath)
        }

        mimeType = mime.TypeByExtension(ext)
        if mimeType == "" {
            log.Panicf("no mime-type given and no mime-type could be determined: [%s]", localFilepath)
        }
    }

    erl = &EmbeddedResourceLocator{
        MimeType: mimeType,
        Filepath: localFilepath,
    }

    if readImmediately == true {
        err = erl.materialize()
        log.PanicIf(err)
    }

    return erl, nil
}

func (erl *EmbeddedResourceLocator) materialize() (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
            log.Panic(err)
        }
    }()

    if erl.Base64EncodedData == "" {
        raw, err := ioutil.ReadFile(erl.Filepath)
        log.PanicIf(err)

        encoded := base64.StdEncoding.EncodeToString(raw)
        erl.Base64EncodedData = encoded
    }

    return nil
}

func (erl *EmbeddedResourceLocator) Uri() string {
    if erl.Base64EncodedData == "" {
        if erl.Filepath == "" {
            log.Panicf("no data present in embedded resource locator but no file-path stored to read from")
        }

        err := erl.materialize()
        log.PanicIf(err)
    }

    return fmt.Sprintf("data:%s;base64,%s", erl.MimeType, erl.Base64EncodedData)
}

// Interface.

type ResourceLocator interface {
    Uri() string
}
