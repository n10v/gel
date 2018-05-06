package gel

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/bogem/gel/fsutils"
	"github.com/bogem/gel/urlutils"
	blackfriday "gopkg.in/russross/blackfriday.v1"
)

type Page struct {
	Description string `toml:"description"`
	Layout      string `toml:"layout"`
	Title       string `toml:"title"`
	Site        *Site

	contentIndex   int
	distPath       string
	rawFileContent []byte
	srcPath        string
}

func ParsePage(fileContent []byte, srcPath string, site *Site) (*Page, error) {
	relSrc, err := filepath.Rel(site.ContentDir, srcPath)
	if err != nil {
		return nil, err
	}

	p := &Page{rawFileContent: fileContent, Site: site, srcPath: relSrc}

	contentIndex, err := p.parseFrontmatter()
	if err != nil {
		return nil, err
	}

	p.contentIndex = contentIndex

	return p, nil
}

func (p *Page) Content() template.HTML {
	return template.HTML(bytes.TrimSpace(blackfriday.MarkdownBasic(p.rawFileContent[p.contentIndex:])))
}

func (p *Page) Dir() string {
	return filepath.Dir(p.SrcPath())
}

// Relative to DistDir.
func (p *Page) DistPath() string {
	if p.distPath == "" {
		if p.Site.UglyURLs || fsutils.DeleteExt(filepath.Base(p.SrcPath())) == "index" {
			p.distPath = fsutils.ChangeExt(p.SrcPath(), ".html")
		} else {
			p.distPath = urlutils.CleanURL(p.SrcPath())
		}
	}
	return p.distPath
}

// Relative to ContentDir.
func (p *Page) SrcPath() string {
	return p.srcPath
}

func (p *Page) URL() string {
	distPath := p.DistPath()
	if !p.Site.UglyURLs {
		distPath = filepath.Dir(distPath)
	}
	if distPath == "." {
		distPath = ""
	}
	return "/" + strings.TrimPrefix(distPath, "/")
}
