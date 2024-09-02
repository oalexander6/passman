package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/pkg/models"
)

func (s *Server) handleGetAllNotes(ctx *gin.Context) {
	notes, err := s.models.NoteGetAll(ctx)
	if err != nil {
		json(ctx, http.StatusInternalServerError, gin.H{"error": "Something went wrong while getting notes."})
		return
	}

	json(ctx, http.StatusOK, gin.H{"notes": notes})
}

func (s *Server) handleCreateNote(ctx *gin.Context) {
	var createNoteParams models.NoteCreateParams

	if err := ctx.ShouldBind(&createNoteParams); err != nil {
		json(ctx, http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request: %s.", err)})
		return
	}

	note, err := s.models.NoteCreate(ctx, createNoteParams)
	if err != nil {
		json(ctx, http.StatusInternalServerError, gin.H{"error": "Something went wrong while saving note."})
		return
	}

	json(ctx, http.StatusCreated, gin.H{"note": note})
}
