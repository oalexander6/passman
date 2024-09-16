package httpserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/models"
)

func HandleHello() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Ok")
	}
}

func HandleGetAllNotes(m models.Models) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		notes, err := m.NoteGetAll(ctx)
		if err != nil {
			json(ctx, http.StatusInternalServerError, gin.H{"error": "Something went wrong while getting notes."})
			return
		}

		json(ctx, http.StatusOK, gin.H{"notes": notes})
	}
}

func HandleCreateNote(m models.Models) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var createNoteParams models.NoteCreateParams

		if err := ctx.ShouldBindJSON(&createNoteParams); err != nil {
			json(ctx, http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request: %s.", err)})
			return
		}

		note, err := m.NoteCreate(ctx, createNoteParams)
		if err != nil {
			json(ctx, http.StatusInternalServerError, gin.H{"error": "Something went wrong while saving note."})
			return
		}

		json(ctx, http.StatusCreated, gin.H{"note": note})
	}
}
