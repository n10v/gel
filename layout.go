package gel

import (
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/bogem/gel/pools"
)

type Layouts struct {
	tf        *templateFuncs
	tmpl      *template.Template
	layoutDir string
}

func ParseLayouts(layoutDir string) (*Layouts, error) {
	layouts := &Layouts{layoutDir: layoutDir}
	tmpl := template.New("root")
	tf := templateFuncs{
		tmpl: tmpl,
	}
	layouts.tmpl = tmpl.Funcs(tf.builtinFuncs())

	err := filepath.Walk(layoutDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		buf := pools.GetBytesBuffer()
		defer pools.PutBytesBuffer(buf)

		if _, err := buf.ReadFrom(f); err != nil {
			return err
		}

		name, err := filepath.Rel(layoutDir, path)
		if err != nil {
			return err
		}
		if filepath.Base(name)[0] == '.' {
			return nil
		}

		tmpl, err = tmpl.New(name).Parse(buf.String())
		return err
	})
	if err != nil {
		return nil, err
	}

	return layouts, err
}

func (l *Layouts) SetFuncs(funcs template.FuncMap) {
	l.tmpl = l.tmpl.Funcs(funcs)
}

func (l *Layouts) ExecuteLayout(w io.Writer, name string, page *Page) error {
	return l.tmpl.ExecuteTemplate(w, name, page)

}
