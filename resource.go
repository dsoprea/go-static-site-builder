package sitebuilder

import (
    "fmt"
    "path"

    "github.com/dsoprea/go-logging"
)

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

type ProjectPageLocalResourceLocator struct {
    sb     *SiteBuilder
    PageId string
}

func NewProjectPageLocalResourceLocator(sb *SiteBuilder, pageId string) (pplrl *ProjectPageLocalResourceLocator) {
    return &ProjectPageLocalResourceLocator{
        sb:     sb,
        PageId: pageId,
    }
}

func (pplrl *ProjectPageLocalResourceLocator) Uri() string {
    if found := pplrl.sb.PageIsValid(pplrl.PageId); found == false {
        log.Panicf("resource refers to invalid page-ID [%s]", pplrl.PageId)
    }

    outputPath := pplrl.sb.Context().HtmlOutputPath()
    filename := pplrl.sb.Context().GetFinalPageFilename(pplrl.PageId)
    filepath := path.Join(outputPath, filename)

    return fmt.Sprintf("file://%s", filepath)
}

type ResourceLocator interface {
    Uri() string
}
