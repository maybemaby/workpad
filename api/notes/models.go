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
	Id          int    `json:"id" required:"true"`
	ProjectName string    `json:"project_name" required:"true" db:"project_name"`
	NoteId      int    `json:"note_id" required:"true" db:"note_id"`
	Excerpt     string `json:"excerpt" required:"true"`
	Date        string `json:"date" required:"true" db:"note_date"`
}

type UpdateNoteExcerptRequest struct {
	Date     string        `json:"date" required:"true" example:"2026-01-01"`
	Excerpts []ExcerptNode `json:"excerpts" required:"true" nullable:"false"`
}

type UpdateNoteExcerptData struct {
	Excerpts []ExcerptNode `json:"excerpts" required:"true"`
}

type ExcerptNode struct {
	Node     string   `json:"node" required:"true" example:"{\"type\":\"paragraph\",\"content\":[{\"type\":\"text\",\"text\":\"Sample excerpt text.\"}]}"`
	Projects []string `json:"projects" example:"[Project A,Project B]" required:"true" nullable:"false"`
}
