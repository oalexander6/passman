package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/pkg/models"
)

func (s *Server) handleGetAllNotes(ctx *gin.Context) {
	notes, err := s.models.NoteGetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	json(ctx, http.StatusOK, gin.H{"notes": notes})
}

func (s *Server) handleCreateNote(ctx *gin.Context) {
	var createNoteParams models.NoteCreateParams

	if err := ctx.ShouldBind(&createNoteParams); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	note, err := s.models.NoteCreate(ctx, createNoteParams)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	json(ctx, http.StatusCreated, gin.H{"note": note})
}
