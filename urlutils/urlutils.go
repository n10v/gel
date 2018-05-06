package urlutils

import (
	"path/filepath"

	"github.com/bogem/gel/fsutils"
)

func CleanURL(url string) string {
	base := filepath.Base(url)
	return filepath.Join(filepath.Dir(url), fsutils.DeleteExt(base), "index.html")
}
