package bundle


import (
	"io"
	"text/template"
	"github.com/google/uuid"
)


func option(bundle Bundle, key string, def string) string {
	if v, ok := bundle.Options[key]; ok {
		return v
	}
	return def
}


func execTemplate(w io.Writer, templ string, data map[string]interface{}) error {

	funcs := template.FuncMap{
		"guid": func() string {
			return uuid.New().String()
		},
	}

	tmpl, err := template.New("BundlerTemplate").Funcs(funcs).Parse(templ)
	if err != nil {
		return err
	}

	return tmpl.Execute(w, data)
}

