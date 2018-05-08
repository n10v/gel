package gel

import (
	"os"
	"strings"
	"sync"

	"path/filepath"

	"github.com/bogem/gel/fsutils"
)

type Site struct {
	BaseURL         string
	ContentDir      string
	DefaultLayout   string
	DistDir         string
	LayoutDir       string
	LayoutStaticDir string
	StaticDir       string
	Title           string
	UglyURLs        bool

	layouts *Layouts
	pages   []*Page
}

func (s *Site) PagesInDir(dir string) []*Page {
	pages := []*Page{}
	for _, page := range s.pages {
		if strings.HasPrefix(page.SrcPath(), dir) {
			pages = append(pages, page)
		}
	}
	return pages
}

func (s *Site) ParsePages() error {
	pages, err := s.WalkContentDir(s.ContentDir)
	if err != nil {
		return err
	}
	s.pages = pages

	s.layouts, err = ParseLayouts(s.LayoutDir)
	return err
}

func (s *Site) WriteContent() error {
	if err := fsutils.CreateMissingDir(s.DistDir); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(s.pages))
	ec := make(chan error, len(s.pages))

	for _, page := range s.pages {
		go func() {
			defer wg.Done()

			pageDistPath := filepath.Join(s.DistDir, page.DistPath())

			if err := fsutils.CreateMissingDir(filepath.Dir(pageDistPath)); err != nil {
				ec <- err
				return
			}

			f, err := os.Create(pageDistPath)
			if err != nil {
				ec <- err
				return
			}

			layout := page.Layout
			if layout == "" {
				layout = s.DefaultLayout
			}
			err = s.layouts.ExecuteLayout(f, layout, page)
			if err != nil {
				ec <- err
				return
			}
		}()
	}
	wg.Wait()
	close(ec)

	return <-ec
}

func (s *Site) CopyStatic() error {
	for _, staticDir := range []string{s.LayoutStaticDir, s.StaticDir} {
		if staticDir == "" {
			continue
		}

		if err := fsutils.Copy(s.DistDir, staticDir); err != nil {
			return err
		}
	}
	return nil
}

func (s *Site) WalkContentDir(contentDir string) ([]*Page, error) {
	pages := []*Page{}

	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".md" {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			fileContent := make([]byte, info.Size())
			if _, err := f.Read(fileContent); err != nil {
				return err
			}
			page, err := ParsePage(fileContent, path, s)
			if err != nil {
				return err
			}

			pages = append(pages, page)
		}
		return nil
	})

	return pages, err
}
