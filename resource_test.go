package sitebuilder

import (
    "bytes"
    "io/ioutil"
    "os"
    "path"
    "testing"

    "github.com/dsoprea/go-logging"
)

func TestNewEmbeddedResourceLocatorWithBytes(t *testing.T) {
    raw := []byte{1, 2, 3}

    erl, err := NewEmbeddedResourceLocatorWithBytes("mime/type", raw)
    log.PanicIf(err)

    uri := erl.Uri()

    if uri != "data:mime/type;base64,AQID" {
        log.Panicf("encoding not correct")
    }
}

func TestNewEmbeddedResourceLocatorWithReader(t *testing.T) {
    raw := []byte{1, 2, 3}
    b := bytes.NewBuffer(raw)

    erl, err := NewEmbeddedResourceLocatorWithReader("mime/type", b)
    log.PanicIf(err)

    uri := erl.Uri()

    if uri != "data:mime/type;base64,AQID" {
        log.Panicf("encoding not correct")
    }
}

func getTempFilepath(filename string) (tempPath string, f *os.File) {
    tempPath, err := ioutil.TempDir("", "")
    log.PanicIf(err)

    filepath := path.Join(tempPath, filename)

    f, err = os.Create(filepath)
    log.PanicIf(err)

    return tempPath, f
}

func NewEmbeddedResourceLocator_DetectMimetype_Deferred(t *testing.T) {
    tempPath, f := getTempFilepath("resource.png")
    defer os.RemoveAll(tempPath)

    raw := []byte{1, 2, 3}

    _, err := f.Write(raw)
    log.PanicIf(err)

    err = f.Sync()
    log.PanicIf(err)

    erl, err := NewEmbeddedResourceLocator(f.Name(), "", false)
    log.PanicIf(err)

    if erl.Base64EncodedData != "" {
        log.Panicf("encoded data *should not* have been read/set yet")
    }

    uri := erl.Uri()

    if erl.Base64EncodedData == "" {
        log.Panicf("encoded data *should* have been read/set by now")
    }

    if uri != "data:image/png;base64,AQID" {
        log.Panicf("encoding not correct: [%v]", uri)
    }
}

func NewEmbeddedResourceLocator_NoDetectMimetype_Deferred(t *testing.T) {
    f, err := ioutil.TempFile("", "resource")
    log.PanicIf(err)

    defer os.Remove(f.Name())

    raw := []byte{1, 2, 3}

    _, err = f.Write(raw)
    log.PanicIf(err)

    err = f.Sync()
    log.PanicIf(err)

    erl, err := NewEmbeddedResourceLocator(f.Name(), "image/png", false)
    log.PanicIf(err)

    if erl.Base64EncodedData != "" {
        log.Panicf("encoded data *should not* have been read/set yet")
    }

    uri := erl.Uri()

    if erl.Base64EncodedData == "" {
        log.Panicf("encoded data *should* have been read/set by now")
    }

    if uri != "data:image/png;base64,AQID" {
        log.Panicf("encoding not correct: [%v]", uri)
    }
}

func NewEmbeddedResourceLocator_DetectMimetype_NoDefer(t *testing.T) {
    tempPath, f := getTempFilepath("resource.png")
    defer os.RemoveAll(tempPath)

    raw := []byte{1, 2, 3}

    _, err := f.Write(raw)
    log.PanicIf(err)

    err = f.Sync()
    log.PanicIf(err)

    erl, err := NewEmbeddedResourceLocator(f.Name(), "", true)
    log.PanicIf(err)

    if erl.Base64EncodedData == "" {
        log.Panicf("encoded data *should* have been read/set by now")
    }

    uri := erl.Uri()

    if uri != "data:image/png;base64,AQID" {
        log.Panicf("encoding not correct: [%v]", uri)
    }
}
