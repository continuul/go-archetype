package generate

import (
	"bytes"
	"fmt"
	"github.com/continuul/go-archetype/generated/archetype"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	// archetype to use
	archetypeName string
	archetypePath string

	config Config

	// Cmd represents the generate command
	Cmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a starter project from an archetype",
		Long: `Generates a new project from an archetype, or
updates the actual project if using a partial archetype.`,
		Run: func(cmd *cobra.Command, args []string) {
			if _, err := archetype.AssetDir(archetypeName); err != nil {
				fmt.Errorf("archetypes does not exist: %s", err)
			}
			result, err := parse(archetypeName, "archetype-metadata.yaml")
			if err != nil {
				fmt.Errorf("archetypes does not have metadata: %s", err)
			}
			config = result
			if err := RestoreAssets(archetypeName, ""); err != nil {
				fmt.Errorf("failed to write file: %s", err)
			}
		},
	}
)

func init() {
	Cmd.PersistentFlags().StringVarP(&archetypeName, "name", "n", "", "name of selected archetype")
	Cmd.PersistentFlags().StringVarP(&archetypePath, "path", "p", "", "path to target directory")
}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := archetype.Asset(name)
	if err != nil {
		return err
	}
	info, err := archetype.AssetInfo(name)
	if err != nil {
		return err
	}
	path := _filePath(dir, filepath.Dir(name))
	if path != "" {
		err = os.MkdirAll(path, os.FileMode(0755))
		if err != nil {
			return err
		}
	}
	var params = make(map[string]interface{})
	for k, v := range config.Parameters {
		params[k] = v.Default
	}
	tmpl, err := template.New(name).Parse(string(data))
	if err != nil {
		return err
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, params)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(_filePath(dir, name), b.Bytes(), info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// Parameter of the archetype
type Parameter struct {
	Name        string `yaml:"name,omitempty"`
	Description string `yaml:"description,omitempty"`
	Default     string `yaml:"default,omitempty"`
	Type        string `yaml:"type,omitempty"`
}

// Config for the archetype
type Config struct {
	Parameters map[string]Parameter `yaml:"parameters,omitempty" mapstructure:"parameters"`
}

// Parse parses an archetype config file.
func parse(dir, name string) (c Config, err error) {
	var data []byte
	data, err = archetype.Asset(filepath.Join(dir, name))

	var raw map[string]interface{}
	err = yaml.Unmarshal(data, &raw)

	var result Config
	config := &mapstructure.DecoderConfig{
		Result:  &result,
		TagName: "yaml",
	}
	d, err := mapstructure.NewDecoder(config)
	if err != nil {
		log.Fatal(err)
	}
	if err := d.Decode(raw); err != nil {
		log.Fatal(err)
	}
	return result, nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := archetype.AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _changePath(path string) string {
	i := 0
	a := strings.Split(path, "/")
	copy(a[i:], a[i+1:])
	a[len(a)-1] = archetypePath
	a = a[:len(a)-1]
	path = strings.Join(a, "/")
	return path
}

func _filePath(dir, name string) string {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return _changePath(filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...))
}
