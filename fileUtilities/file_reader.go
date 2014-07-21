package fileUtilities

import "io/ioutil"

var NotificationsRoot string

func ReadFile(path string) (string, error) {
    buffer, err := ioutil.ReadFile(path)
    if err != nil {
        return "", err
    }

    contents := string(buffer)

    return contents, nil
}
