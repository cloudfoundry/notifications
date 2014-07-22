package fileUtilities

import (
    "io/ioutil"
    "os"
)

func ReadFile(path string) (string, error) {
    buffer, err := ioutil.ReadFile(path)
    if err != nil {
        return "", err
    }

    contents := string(buffer)

    return contents, nil
}

func FileExists(path string) bool {
    _, err := os.Stat(path)
    if err != nil {
        return false
    }
    return true
}
