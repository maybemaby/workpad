package notes

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/maybemaby/workpad/api/utils"
)

type NoteHandler struct {
	noteStore NoteStore
}

func NewNoteHandler(noteStore NoteStore) *NoteHandler {
	return &NoteHandler{noteStore: noteStore}
}

type GetNoteByDateRequest struct {
	Date string `query:"date" example:"2026-01-01" required:"true"`
}

type GetMonthNotesRequest struct {
	Year  int `query:"year" example:"2026" required:"true"`
	Month int `query:"month" example:"1" required:"true"`
}

func (h *NoteHandler) GetNoteByDate(w http.ResponseWriter, r *http.Request) {

	date := r.URL.Query().Get("date")

	parsedDate, err := time.Parse("2006-01-02", date)

	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD.", http.StatusBadRequest)
		return
	}

	note, err := h.noteStore.GetNoteByDate(r.Context(), parsedDate)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Note not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	err = utils.WriteJSON(w, r, note)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req CreateNoteRequest

	if err := utils.ReadJSON(r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	currentDate := time.Now().Local()

	note, err := h.noteStore.CreateNote(r.Context(), req.HTMLContent, currentDate)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = utils.WriteJSON(w, r, note)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *NoteHandler) GetMonthNotes(w http.ResponseWriter, r *http.Request) {
	month := r.URL.Query().Get("month")
	year := r.URL.Query().Get("year")

	monthInt, err := strconv.Atoi(month)

	if err != nil || monthInt < 1 || monthInt > 12 {
		http.Error(w, "Invalid month parameter", http.StatusBadRequest)
		return
	}

	yearInt, err := strconv.Atoi(year)

	if err != nil || yearInt < 1 {
		http.Error(w, "Invalid year parameter", http.StatusBadRequest)
		return
	}

	days, err := h.noteStore.GetNoteDatesForMonth(r.Context(), yearInt, time.Month(monthInt))

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, r, days)
}

func (h *NoteHandler) UpdateNoteExcerpts(w http.ResponseWriter, r *http.Request) {

	var data UpdateNoteExcerptRequest

	if err := utils.ReadJSON(r, &data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	parsedDate, err := time.Parse(time.DateOnly, data.Date)

	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD.", http.StatusBadRequest)
		return
	}

	err = h.noteStore.UpdateExcerptsForDate(r.Context(), parsedDate, data.Excerpts)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type GetExcerptsForProjectRequest struct {
	Project string `path:"project" example:"Project A" required:"true"`
}

func (h *NoteHandler) GetExcerptsForProject(w http.ResponseWriter, r *http.Request) {
	projectName := r.PathValue("project")

	excerpts, err := h.noteStore.GetExcerptsForProject(r.Context(), projectName)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, r, excerpts)
}
