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
