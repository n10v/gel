package gel

import (
	"path/filepath"
	"strings"
)

type PageUtils struct {
	contentDir string
	distDir    string
}

func NewPageUtils(site *Site) *PageUtils {
	return &PageUtils{
		contentDir: site.ContentDir,
		distDir:    site.DistDir,
	}
}

func (pu *PageUtils) ResolveDistPath(pageSourcePath string) (string, error) {
	relPath, err := filepath.Rel(pu.contentDir, pageSourcePath)
	if err != nil {
		return "", err
	}
	return changeExt(filepath.Join(pu.distDir, relPath), ".html"), nil
}

func changeExt(path, newExt string) string {
	oldExt := filepath.Ext(path)
	return strings.TrimSuffix(path, oldExt) + newExt
}
