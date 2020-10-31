package bundle

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func msi(bundler Bundler, bundle Bundle) error {

	cfg, guid := option(bundle, "template", "app.wxs"), option(bundle, "guid", "")
	if guid == "" {
		return fmt.Errorf("could not find msi guid option.")
	}

	cfgData, err := ioutil.ReadFile(cfg)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("tmp_bundler_%s", cfg))
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()

	data := map[string]interface{}{
		"ID":          bundler.ID,
		"Ver":         strings.TrimPrefix(bundler.Version, "v"),
		"GUID":        guid,
		"Name":        bundler.Name,
		"Version":     bundler.Version,
		"Source":      bundle.Source,
		"Publisher":   bundler.Publisher,
		"Description": bundler.Description,
	}

	for i, v := range bundler.Data {
		data[i] = v
	}

	if err = execTemplate(f, string(cfgData), data); err != nil {
		return err
	}

	// TODO: change cmd name on windows.

	if out, err := run("wixl", []string{"-o", bundle.Output, f.Name()}, nil, bundle.Source); err != nil {
		return fmt.Errorf("wix error: %s", string(out))
	}

	return nil
}
