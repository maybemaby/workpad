package notes

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type NoteStore interface {
	GetNoteByDate(ctx context.Context, date time.Time) (Note, error)
	CreateNote(ctx context.Context, htmlContent string, date time.Time) (Note, error)
}

type NoteService struct {
	db *sqlx.DB
}

func NewNoteService(db *sqlx.DB) *NoteService {
	return &NoteService{db: db}
}

func (s *NoteService) GetNoteByDate(ctx context.Context, date time.Time) (Note, error) {
	var note Note

	err := s.db.GetContext(ctx, &note, "SELECT id, html_content, note_date FROM notes WHERE date(note_date) = ?", date.Format("2006-01-02"))

	return note, err
}

func (s *NoteService) CreateNote(ctx context.Context, htmlContent string, date time.Time) (Note, error) {
	var id int

	err := s.db.QueryRowContext(ctx, `INSERT INTO notes (html_content, note_date) VALUES (?, ?) ON CONFLICT (note_date) DO UPDATE SET html_content = excluded.html_content RETURNING id`, htmlContent, date.Format("2006-01-02")).Scan(&id)

	if err != nil {
		return Note{}, err
	}

	return Note{
		Id:          id,
		HTMLContent: htmlContent,
		Date:        date,
	}, nil
}
