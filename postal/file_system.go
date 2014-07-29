package postal

import (
    "io/ioutil"
    "os"
)

type FileSystemInterface interface {
    Exists(string) bool
    Read(string) (string, error)
}

type FileSystem struct{}

func NewFileSystem() FileSystem {
    return FileSystem{}
}

func (fs FileSystem) Exists(path string) bool {
    _, err := os.Stat(path)
    if err != nil {
        return false
    }

    return true
}

func (fs FileSystem) Read(path string) (string, error) {
    bytes, err := ioutil.ReadFile(path)
    if err != nil {
        return "", err
    }

    return string(bytes), nil
}
