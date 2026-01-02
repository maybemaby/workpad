package notes

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/maybemaby/workpad/api/projects"
	"github.com/maybemaby/workpad/api/utils"
	"github.com/stretchr/testify/suite"
	_ "modernc.org/sqlite"
)

func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}

	return t
}

func seedProjects(ctx context.Context, db *sqlx.DB) []projects.Project {
	projectNames := []string{"Project Alpha", "Beta Project", "Gamma", "Delta Force"}

	tx := db.MustBegin()

	defer tx.Commit()
	var projectsInserted []projects.Project

	for _, name := range projectNames {
		var proj projects.Project
		row := tx.QueryRowContext(ctx, "INSERT INTO projects (name) VALUES (?) RETURNING name, created_at", name)

		if err := row.Scan(&proj.Name, &proj.CreatedAt); err != nil {
			panic(err)
		}

		projectsInserted = append(projectsInserted, proj)
	}

	return projectsInserted
}

func seedNotes(ctx context.Context, db *sqlx.DB) []Note {
	notes := []Note{
		{HTMLContent: "<p>Note for 2026-01-01</p>", Date: mustParseTime(time.DateOnly, "2026-01-01")},
		{HTMLContent: "<p>Note for 2026-01-02</p>", Date: mustParseTime(time.DateOnly, "2026-01-02")},
		{HTMLContent: "<p>Note for 2026-01-03</p>", Date: mustParseTime(time.DateOnly, "2026-01-03")},
		{HTMLContent: "<p>Note for 2026-03-10</p>", Date: mustParseTime(time.DateOnly, "2025-01-10")},
		{HTMLContent: "<p>Note for 2026-03-15</p>", Date: mustParseTime(time.DateOnly, "2026-03-15")},
	}

	tx := db.MustBegin()

	defer tx.Commit()

	for i := range notes {
		row := tx.QueryRowContext(ctx, "INSERT INTO notes (html_content, note_date) VALUES (?, ?) RETURNING id", notes[i].HTMLContent, notes[i].Date.Format("2006-01-02"))

		if err := row.Scan(&notes[i].Id); err != nil {
			panic(err)
		}
	}

	return notes
}

type NoteStoreSuite struct {
	suite.Suite
	db  *sql.DB
	dbx *sqlx.DB
}

func (s *NoteStoreSuite) SetupTest() {
	s.db, _ = sql.Open("sqlite", ":memory:")
	s.dbx = sqlx.NewDb(s.db, "sqlite")

	err := utils.SetupSqliteDb(s.db)

	if err != nil {
		panic(err)
	}

	seedNotes(s.T().Context(), s.dbx)
	seedProjects(s.T().Context(), s.dbx)
}

func (s *NoteStoreSuite) TearDownTest() {
	s.db.Close()
}

func (s *NoteStoreSuite) TestGetNoteByDate_Found() {
	store := NewNoteService(s.dbx)

	note, err := store.GetNoteByDate(s.T().Context(), mustParseTime(time.DateOnly, "2026-01-02"))

	s.NoError(err)
	s.Equal("<p>Note for 2026-01-02</p>", note.HTMLContent)
	s.Equal("2026-01-02", note.Date.Format(time.DateOnly))
}

func (s *NoteStoreSuite) TestGetNoteByDate_NotFound() {
	store := NewNoteService(s.dbx)

	_, err := store.GetNoteByDate(s.T().Context(), mustParseTime(time.DateOnly, "2026-12-31"))

	s.ErrorIs(err, sql.ErrNoRows)
}

func (s *NoteStoreSuite) TestCreateNote_NewNote() {
	store := NewNoteService(s.dbx)

	note, err := store.CreateNote(s.T().Context(), "<p>New Note</p>", mustParseTime(time.DateOnly, "2026-04-01"))

	s.NoError(err)
	s.Equal("<p>New Note</p>", note.HTMLContent)
	s.Equal("2026-04-01", note.Date.Format(time.DateOnly))
}

func (s *NoteStoreSuite) TestCreateNote_UpdateExistingNote() {
	store := NewNoteService(s.dbx)

	note, err := store.CreateNote(s.T().Context(), "<p>Updated Note for 2026-01-02</p>", mustParseTime(time.DateOnly, "2026-01-02"))

	s.NoError(err)
	s.Equal("<p>Updated Note for 2026-01-02</p>", note.HTMLContent)
	s.Equal("2026-01-02", note.Date.Format(time.DateOnly))
}

func (s *NoteStoreSuite) TestGetNoteDatesForMonth() {
	store := NewNoteService(s.dbx)

	days, err := store.GetNoteDatesForMonth(s.T().Context(), 2026, time.January)

	s.NoError(err)
	s.ElementsMatch([]int{1, 2, 3}, days)
}

func (s *NoteStoreSuite) TestGetNoteDatesForMonth_NoNotes() {
	store := NewNoteService(s.dbx)
	days, err := store.GetNoteDatesForMonth(s.T().Context(), 2026, time.February)

	s.NoError(err)
	s.Empty(days)
}

func (s *NoteStoreSuite) TestUpdateExcerpts() {
	store := NewNoteService(s.dbx)

	excerpts := []ExcerptNode{
		{Node: "Excerpt 1", Projects: []string{"Project Alpha", "Beta Project"}},
		{Node: "Excerpt 2", Projects: []string{"Gamma"}},
	}

	err := store.UpdateExcerptsForDate(s.T().Context(), mustParseTime(time.DateOnly, "2026-01-02"), excerpts)

	s.NoError(err)
}

func (s *NoteStoreSuite) TestUpdateExcerpts_NoteNotFound() {
	store := NewNoteService(s.dbx)

	excerpts := []ExcerptNode{
		{Node: "Excerpt 1", Projects: []string{"Project Alpha", "Beta Project"}},
	}

	err := store.UpdateExcerptsForDate(s.T().Context(), mustParseTime(time.DateOnly, "2026-12-31"), excerpts)

	s.ErrorIs(err, sql.ErrNoRows)
}

func (s *NoteStoreSuite) TestUpdateExcerpts_ProjectNotFound() {
	store := NewNoteService(s.dbx)

	excerpts := []ExcerptNode{
		{Node: "Excerpt 1", Projects: []string{"Nonexistent Project"}},
	}

	err := store.UpdateExcerptsForDate(s.T().Context(), mustParseTime(time.DateOnly, "2026-01-02"), excerpts)

	s.NoError(err)
}

func TestNoteStoreSuite(t *testing.T) {
	suite.Run(t, new(NoteStoreSuite))
}
