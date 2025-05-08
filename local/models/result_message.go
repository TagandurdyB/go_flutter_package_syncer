package models

type ResultMessage struct {
	Local       string   `json:"local"`
	Server      string   `json:"server"`
	DiffMessage string   `json:"diff_message"`
	Diff        []string `json:"diff"`
}