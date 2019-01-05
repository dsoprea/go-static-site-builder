package sitebuilder

import (
    "fmt"
)

type LocalResourceLocator struct {
    localFilepath string
}

func NewLocalResourceLocator(localFilepath string) (lrl *LocalResourceLocator) {
    return &LocalResourceLocator{
        localFilepath: localFilepath,
    }
}

func (lrl *LocalResourceLocator) Uri() string {
    return fmt.Sprintf("file://%s", lrl.localFilepath)
}

type ResourceLocator interface {
    Uri() string
}
