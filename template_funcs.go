package gel

import (
	"fmt"
	"html/template"
	"path"
	"strings"

	"github.com/bogem/gel/pools"
)

var (
	partialsDir   = "partials/"
	componentsDir = "components/"
)

type templateFuncs struct {
	layoutDir string
	tmpl      *template.Template
}

func (tf *templateFuncs) builtinFuncs() template.FuncMap {
	return map[string]interface{}{
		"component": tf.component,
		"partial":   tf.partial,
	}
}

func (tf *templateFuncs) component(name string, data interface{}) (template.HTML, error) {
	name = addPrefixIfNeeded(name, componentsDir)
	name = addPrefixIfNeeded(name, partialsDir)
	return tf.partial(name, data)
}

func (tf *templateFuncs) partial(name string, data interface{}) (template.HTML, error) {
	name = addPrefixIfNeeded(name, partialsDir)
	name = addSuffixIfNeeded(name, ".html")
	name = path.Clean(name)

	tmpl := tf.tmpl.Lookup(name)
	if tmpl == nil {
		return template.HTML(""), fmt.Errorf("partial %q not found", name)
	}

	buf := pools.GetBytesBuffer()
	defer pools.PutBytesBuffer(buf)

	if err := tmpl.Execute(buf, data); err != nil {
		return template.HTML(""), err
	}

	return template.HTML(buf.String()), nil

}

func addPrefixIfNeeded(s, prefix string) string {
	if !strings.HasPrefix(s, prefix) {
		return prefix + s
	}
	return s
}

func addSuffixIfNeeded(s, suffix string) string {
	if !strings.HasSuffix(s, suffix) {
		return s + suffix
	}
	return s
}
