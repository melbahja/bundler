package bundle

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
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

	// TODO: create this is tmp dir.
	f, err := os.Create(fmt.Sprintf("tmp%s.wxs", bundler.ID))
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

	// Fix for windows wix!!!!
	f.Close()

	if runtime.GOOS == "windows" {

		wixBin := option(bundle, "bindir", "c:\\Program Files (x86)\\Wix Toolset v3.11\\bin\\")
		if err = run(fmt.Sprintf("%s\\candle", wixBin), []string{f.Name()}, nil, bundle.Source); err != nil {
			return fmt.Errorf("wix candle error: %s", err)
		}

		err = run(
			fmt.Sprintf("%slight", wixBin),
			[]string{"-ext", "WixUIExtension", "-cultures:en-us", "-out", bundle.Output, fmt.Sprintf("tmp%s.wixobj", bundler.ID)},
			nil,
			bundle.Source,
		)

		if err != nil {
			return fmt.Errorf("wix light error: %s", err)
		}

		return nil
	}

	if err := run("wixl", []string{"-o", bundle.Output, f.Name()}, nil, bundle.Source); err != nil {
		return fmt.Errorf("wix error: %s", err)
	}

	return nil
}
