package notes

import (
	"database/sql"
	"errors"
	"net/http"
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

	utils.WriteJSON(w, r, note)
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

	utils.WriteJSON(w, r, note)
}
