package gel

import (
	"bytes"
	"errors"

	"github.com/BurntSushi/toml"
)

var errNoFrontmatter = errors.New("no frontmatter found")

var frontmatterDelim = []byte("+++")

func (p *Page) parseFrontmatter() (contentIndex int, err error) {
	startIndex, endIndex := findFrontmatter(p.rawFileContent)
	if startIndex < 0 || endIndex <= 0 {
		return -1, errNoFrontmatter
	}

	if err := toml.Unmarshal(p.rawFileContent[startIndex:endIndex], p); err != nil {
		return -1, err
	}

	return endIndex + len(frontmatterDelim), nil
}

func findFrontmatter(b []byte) (startIndex, endIndex int) {
	startIndex = bytes.Index(b, frontmatterDelim)
	if startIndex < 0 {
		return -1, -1
	}
	startIndex += len(frontmatterDelim)

	endIndex = startIndex + bytes.Index(b[startIndex:], frontmatterDelim)
	return startIndex, endIndex
}
