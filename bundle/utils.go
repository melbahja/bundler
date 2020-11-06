package bundle

import (
	"io"
	"os"
	"os/exec"
	"text/template"

	"github.com/google/uuid"
)

var (
	Stdout = os.Stdout
	Stderr = os.Stderr
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

func run(c string, args []string, env []string, dir string) error {
	cmd := exec.Command(c, args...)
	cmd.Env = os.Environ()
	cmd.Stderr = Stderr
	cmd.Stdout = Stdout
	if env != nil {
		cmd.Env = append(cmd.Env, env...)
	}
	return cmd.Run()
}
