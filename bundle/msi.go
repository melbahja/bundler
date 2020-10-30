package bundle

import (
	"os"
	"fmt"
	"io/ioutil"
	"os/exec"
)

func msi(bundler Bundler, bundle Bundle) error {

	cfg, guid := option(bundle, "wixfile", "app.wxs"), option(bundle, "guid", "")

	if cfg == "" {
		return fmt.Errorf("could not find msi wixfile option.")
	} else if guid == "" {
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
		"ID": bundler.ID,
		"GUID": guid,
		"Name": bundler.Name,
		"Version": bundler.Version,
		"Source": bundle.Source,
		"Publisher": bundler.Publisher,
		"Description": bundler.Description,
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


func run(c string, args []string, env []string, dir string) ([]byte, error) {
	cmd := exec.Command(c, args...)
	cmd.Env = os.Environ()
	if env != nil {
		cmd.Env = append(cmd.Env, env...)
	}
	return cmd.CombinedOutput()
}
