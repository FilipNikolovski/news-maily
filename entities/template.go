package entities

import "time"

type Template struct {
	Name        string `json:"name"`
	HTMLPart    string `json:"html_part"`
	TextPart    string `json:"text_part"`
	SubjectPart string `json:"subject_part"`
}

type TemplateCollection struct {
	NextToken string         `json:"next_token"`
	List      []TemplateMeta `json:"list"`
}

type TemplateMeta struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}
