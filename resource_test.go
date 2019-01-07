package sitebuilder

import (
    "bytes"
    "io/ioutil"
    "os"
    "testing"

    "github.com/dsoprea/go-logging"
)

func TestNewEmbeddedResourceLocator(t *testing.T) {
    raw := []byte{1, 2, 3}

    erl, err := NewEmbeddedResourceLocator("mime/type", raw)
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

func TestNewEmbeddedResourceLocatorWithFilepath_DetectMimetype_Deferred(t *testing.T) {
    f, err := ioutil.TempFile("", "resource*.png")
    log.PanicIf(err)

    defer os.Remove(f.Name())

    raw := []byte{1, 2, 3}

    _, err = f.Write(raw)
    log.PanicIf(err)

    err = f.Sync()
    log.PanicIf(err)

    erl, err := NewEmbeddedResourceLocatorWithFilepath(f.Name(), "", false)
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

func TestNewEmbeddedResourceLocatorWithFilepath_NoDetectMimetype_Deferred(t *testing.T) {
    f, err := ioutil.TempFile("", "resource")
    log.PanicIf(err)

    defer os.Remove(f.Name())

    raw := []byte{1, 2, 3}

    _, err = f.Write(raw)
    log.PanicIf(err)

    err = f.Sync()
    log.PanicIf(err)

    erl, err := NewEmbeddedResourceLocatorWithFilepath(f.Name(), "image/png", false)
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

func TestNewEmbeddedResourceLocatorWithFilepath_DetectMimetype_NoDefer(t *testing.T) {
    f, err := ioutil.TempFile("", "resource*.png")
    log.PanicIf(err)

    defer os.Remove(f.Name())

    raw := []byte{1, 2, 3}

    _, err = f.Write(raw)
    log.PanicIf(err)

    err = f.Sync()
    log.PanicIf(err)

    erl, err := NewEmbeddedResourceLocatorWithFilepath(f.Name(), "", true)
    log.PanicIf(err)

    if erl.Base64EncodedData == "" {
        log.Panicf("encoded data *should* have been read/set by now")
    }

    uri := erl.Uri()

    if uri != "data:image/png;base64,AQID" {
        log.Panicf("encoding not correct: [%v]", uri)
    }
}
