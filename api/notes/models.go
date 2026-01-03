package notes

import "time"

type Note struct {
	HTMLContent string    `json:"html_content" required:"true" db:"html_content"`
	Date        time.Time `json:"note_date" required:"true" db:"note_date"`
	Id          int       `json:"id" required:"true"`
}

type CreateNoteRequest struct {
	HTMLContent string `json:"html_content" required:"true"`
}

type NoteExcerpt struct {
	ProjectName string `json:"project_name" required:"true" db:"project_name"`
	Excerpt     string `json:"excerpt" required:"true"`
	Date        string `json:"date" required:"true" db:"note_date"`
	Id          int    `json:"id" required:"true"`
	NoteId      int    `json:"note_id" required:"true" db:"note_id"`
}

type UpdateNoteExcerptRequest struct {
	Excerpts []ExcerptNode `json:"excerpts" required:"true" nullable:"false"`
	Date     string        `json:"date" required:"true" example:"2026-01-01"`
}

type UpdateNoteExcerptData struct {
	Excerpts []ExcerptNode `json:"excerpts" required:"true"`
}

type ExcerptNode struct {
	Projects []string `json:"projects" example:"[Project A,Project B]" required:"true" nullable:"false"`
	Node     string   `json:"node" required:"true" example:"{\"type\":\"paragraph\",\"content\":[{\"type\":\"text\",\"text\":\"Sample excerpt text.\"}]}"`
}
