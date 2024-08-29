package sip

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type FileLoader func(string, any) error

var Files = []string{}
var FileLoaders = map[string]FileLoader{
	".json": LoadYamlFile,
	".yaml": LoadYamlFile,
	".yml":  LoadYamlFile,
}

func RegisterFile(filename string) {
	Files = append(Files, filename)
}

func RegisterFileLoader(ext string, f FileLoader) {
	FileLoaders[ext] = f
}

func LoadFiles(v any) error {
	for _, filename := range Files {
		if err := LoadFile(filename, v); err != nil {
			return err
		}
	}

	return nil
}

func LoadFile(filename string, v any) error {
	ext := filepath.Ext(filename)
	loader, ok := FileLoaders[ext]
	if !ok {
		return fmt.Errorf("no loader defined for extension %s", ext)
	}

	return loader(filename, v)
}

func LoadYamlFile(filename string, v any) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	return yaml.NewDecoder(f).Decode(v)
}
