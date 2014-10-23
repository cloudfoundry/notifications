package models

type Template struct {
    Name       string
    Text       string `json:"text"`
    HTML       string `json:"html"`
    Overridden bool   `json:"overridden"`
}
