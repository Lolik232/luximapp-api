package domain

import "time"

type News struct {
	ID        string     `json:"id,omitempty"`
	Title     string     `json:"title,omitempty"`
	Text      string     `json:"text,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func (n *News) Sanitize() {
	n.ID = ""
	n.CreatedAt = nil
}
