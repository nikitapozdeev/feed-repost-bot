package vk

import (
	"bytes"
	"encoding/json"
	"html/template"
)

type Response struct {
	Response json.RawMessage `json:"response"`
}

type PostsResponse struct {
	Count int    `json:"count"`
	Posts []Post `json:"items"`
}

type Post struct {
	ID          int          `json:"id"`
	Timestamp   int64        `json:"date"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Type  string `json:"type"`
	Photo Photo  `json:"photo"`
}

type Photo struct {
	Sizes []Size `json:"sizes"`
}

type Size struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Url    string `json:"url"`
	Type   string `json:"type"`
}

func (p *Post) Recipient() string {
	return ""
}

func (p *Post) HTML() (string, error) {
	messageTemplate, err := template.New("msg").Parse(`
		{{ .Text }}
		{{if .ImageUrl}} <a href="{{ .ImageUrl }}">&#8205;</a>{{end}}
	`)
	if err != nil {
		return "", err
	}

	templateData := struct {
		Text     string
		ImageUrl string
	}{
		Text:     p.Text,
		ImageUrl: p.FindLargestImage(),
	}

	var htmlMessage bytes.Buffer
	if err := messageTemplate.Execute(&htmlMessage, templateData); err != nil {
		return "", err
	}

	return htmlMessage.String(), nil
}

func (p *Post) FindLargestImage() string {
	if len(p.Attachments) > 0 {
		sizes := p.Attachments[0].Photo.Sizes
		for _, size := range sizes {
			if size.Type == "z" {
				return size.Url
			}
		}
	}

	return ""
}
