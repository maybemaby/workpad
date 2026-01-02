package notes

import "time"

type Note struct {
	Id          int       `json:"id" required:"true"`
	HTMLContent string    `json:"html_content" required:"true" db:"html_content"`
	Date        time.Time `json:"note_date" required:"true" db:"note_date"`
}

type CreateNoteRequest struct {
	HTMLContent string `json:"html_content" required:"true"`
}

type NoteExcerpt struct {
	Id        int    `json:"id" required:"true"`
	ProjectId int    `json:"project_id" required:"true" db:"project_id"`
	NoteId    int    `json:"note_id" required:"true" db:"note_id"`
	Excerpt   string `json:"excerpt" required:"true"`
	Date      string `json:"date" required:"true" db:"note_date"`
}

type CreateNoteExcerptRequest struct {
}

type ExcerptNode struct {
	Node     string   `json:"node" required:"true"`
	Projects []string `json:"project" required:"true"`
}
