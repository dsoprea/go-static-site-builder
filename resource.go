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
    if found := splrl.sb.PageIsValid(splrl.PageId); found == false {
        log.Panicf("resource refers to invalid page-ID [%s]", splrl.PageId)
    }

    outputPath := splrl.sb.Context().HtmlOutputPath()
    filename := splrl.sb.Context().GetFinalPageFilename(splrl.PageId)
    filepath := path.Join(outputPath, filename)

    return fmt.Sprintf("file://%s", filepath)
}

type ResourceLocator interface {
    Uri() string
}
