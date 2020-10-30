package bundler_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/melbahja/bundler/bundle"
	"gopkg.in/yaml.v2"
)

var (
	bundlerFileData = []byte(`
id: MelbahjaApp
name: App Name
version: v1.1.1
publisher: Example LTD
description: Your app short description here.
bundles:
  - type: msi
    source: dist/windows
    output: dist/MelbahjaApp.msi
    options:
      guid: 07957cbf-c2ac-4fa6-9c26-42b6666ba7f7
      wixfile: app.wxs
`)
)

func TestBundleMSI(t *testing.T) {

	defer os.RemoveAll("dist/")

	bundler := new(bundle.Bundler)

	if err := yaml.Unmarshal(bundlerFileData, bundler); err != nil {
		t.Errorf("Could not unmarshal bundler file data: %s", err.Error())
	}

	if err := build(bundler.ID, "windows", "dist/windows"); err != nil {
		t.Errorf("Build err: %s", err.Error())
	}

	if err := bundler.Run(); err != nil {
		t.Error(err)
	}
}

func build(id, gos, outDir string) error {

	o := fmt.Sprintf("%s/%s", outDir, id)
	if gos == "windows" {
		o += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", o, "../main.go")
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOOS=%s", gos))
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
