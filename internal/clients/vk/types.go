package vk

import "encoding/json"

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
	Type string
}
