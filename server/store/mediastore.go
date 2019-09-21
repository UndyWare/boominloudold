package store

import (
	"fmt"
	"os"
)

type mediastore struct {
	baseDir string
}

func (m *mediastore) Open(baseDirectory string) error {
	//check if directory exists
	m.baseDir = baseDirectory
	return nil
}

func (m *mediastore) Close() error {
	//check if directory exists
	m.baseDir = ""
	return nil
}

func (m *mediastore) Read(filename string) (string, error) {
	fullpath := m.baseDir + "/" + filename
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		return "", nil
	}
	return fullpath, nil
}

func (m *mediastore) Write(filename string, pathToFile string) error {
	fullpath := m.baseDir + "/" + filename
	if _, err := os.Stat(fullpath); !os.IsNotExist(err) {
		return fmt.Errorf("file already exists")
	}
	return os.Rename(pathToFile, fullpath)
}
