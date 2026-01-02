package notes

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type NoteStore interface {
	GetNoteByDate(ctx context.Context, date time.Time) (Note, error)
	CreateNote(ctx context.Context, htmlContent string, date time.Time) (Note, error)
	GetNoteDatesForMonth(ctx context.Context, year int, month time.Month) ([]int, error)
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

func (s *NoteService) GetNoteDatesForMonth(ctx context.Context, year int, month time.Month) ([]int, error) {
	var days []int

	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	rows, err := s.db.QueryxContext(ctx, `SELECT strftime('%d', note_date) FROM notes WHERE note_date >= ? AND note_date <= ?`, startDate.Format(time.DateOnly), endDate.Format(time.DateOnly))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var day int
		if err := rows.Scan(&day); err != nil {
			return nil, err
		}
		days = append(days, day)
	}

	return days, nil
}
