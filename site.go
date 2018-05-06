package gel

import (
	"html/template"
	"os"

	"path/filepath"
)

type Site struct {
	ContentDir string
	DistDir    string

	DefaultLayout string
	LayoutDir     string

	TemplateFuncs template.FuncMap

	layouts   *Layouts
	pages     []*Page
	pageUtils *PageUtils
}

func (s *Site) GetPageUtils() *PageUtils {
	if s.pageUtils == nil {
		s.pageUtils = NewPageUtils(s)
	}
	return s.pageUtils
}

func (s *Site) ParsePages() error {
	pages, err := WalkContentDir(s.ContentDir)
	if err != nil {
		return err
	}
	s.pages = pages

	s.layouts, err = ParseLayouts(s.LayoutDir)
	return err
}

func (s *Site) WriteContent() error {
	if err := createMissingDir(s.DistDir); err != nil {
		return err
	}

	for _, page := range s.pages {
		pageDistPath, err := s.GetPageUtils().ResolveDistPath(page.SourcePath())
		if err != nil {
			return err
		}

		f, err := os.Create(pageDistPath)
		if err != nil {
			return err
		}

		layout := page.Layout
		if layout == "" {
			layout = s.DefaultLayout
		}
		err = s.layouts.ExecuteLayout(f, layout, page)
		if err != nil {
			return err
		}
	}

	return nil
}

func createMissingDir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(path, 0755); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func WalkContentDir(contentDir string) ([]*Page, error) {
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
			page, err := ParsePage(fileContent, path)
			if err != nil {
				return err
			}

			pages = append(pages, page)
		}
		return nil
	})

	return pages, err
}
