package gin_binding

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/pkg/components"
	"github.com/oalexander6/passman/pkg/entities"
	"github.com/oalexander6/passman/pkg/pages"
)

func (b *GinBinding) GetAllNotes(ctx *gin.Context) {
	notes, err := b.services.GetAllNotes(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"notes": notes,
	})
}

func (b *GinBinding) GetNoteByID(ctx *gin.Context) {
	noteID := ctx.Param("id")
	if noteID == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}

	note, err := b.services.GetNoteByID(ctx, noteID)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"note": note,
	})
}

func (b *GinBinding) CreateNote(ctx *gin.Context) {
	var noteInput entities.NoteInput

	if err := ctx.ShouldBind(&noteInput); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	savedNote, err := b.services.CreateNote(ctx, noteInput)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	sendJSONOrHTML(
		ctx,
		http.StatusCreated,
		&gin.H{
			"note": savedNote,
		},
		components.NoteListItem(savedNote),
	)
}

func (b *GinBinding) UpdateNote(ctx *gin.Context) {
	noteID := ctx.Param("id")
	if noteID == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var note entities.Note

	if err := ctx.ShouldBind(&note); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if noteID != note.ID {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	updatedNote, err := b.services.UpdateNote(ctx, note)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"note": updatedNote,
	})
}

func (b *GinBinding) DeleteNote(ctx *gin.Context) {
	noteID := ctx.Param("id")
	if noteID == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := b.services.DeleteNoteByID(ctx, noteID)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (b *GinBinding) ViewHomePage(ctx *gin.Context) {
	notes, err := b.services.GetAllNotes(ctx)
	if err != nil {
		sendJSONOrHTML(ctx, http.StatusInternalServerError, &gin.H{}, pages.Error())
		return
	}

	sendJSONOrHTML(ctx, http.StatusOK, &gin.H{"message": "OK"}, pages.Dashboard(notes))
}
