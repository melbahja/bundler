package bundle

import (
	"fmt"
)

type (
	Bundle struct {
		Type    string
		Source  string
		Output  string
		Options map[string]string
	}

	Bundler struct {
		ID          string
		Name        string
		Data        map[string]interface{}
		Version     string
		Publisher   string
		Description string
		Bundles     []Bundle
	}
)

func (b Bundler) Run() error {
	for _, bundle := range b.Bundles {
		if err := b.Bundle(bundle); err != nil {
			return err
		}
	}
	return nil
}

func (b Bundler) Bundle(bundle Bundle) error {
	switch bundle.Type {
	case "msi":
		return msi(b, bundle)
	}

	return fmt.Errorf("Unknown or invalid bundle type: %s", bundle.Type)
}
