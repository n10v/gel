package gel

import (
	"bytes"
	"html/template"

	blackfriday "gopkg.in/russross/blackfriday.v1"
)

type Page struct {
	Description string `toml:"description"`
	Layout      string `toml:"layout"`
	Title       string `toml:"title"`

	contentIndex   int
	rawFileContent []byte
	sourcePath     string
}

func ParsePage(fileContent []byte, sourcePath string) (*Page, error) {
	p := &Page{rawFileContent: fileContent, sourcePath: sourcePath}

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

func (p *Page) SourcePath() string {
	return p.sourcePath
}
