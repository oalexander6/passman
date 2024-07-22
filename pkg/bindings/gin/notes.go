package gin_binding

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/pkg/components"
	"github.com/oalexander6/passman/pkg/entities"
	csrf "github.com/utrack/gin-csrf"
)

func (b *GinBinding) HandleGetAllNotes(ctx *gin.Context) {
	notes, err := b.services.GetAllNotes(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"notes": notes,
	})
}

func (b *GinBinding) HandleGetNoteByID(ctx *gin.Context) {
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

func (b *GinBinding) HandleCreateNote(ctx *gin.Context) {
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

	notes, err := b.services.GetAllNotes(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	csrfToken := csrf.GetToken(ctx)

	sendJSONOrHTML(ctx, http.StatusCreated, &gin.H{"note": savedNote}, components.NoteList(notes, csrfToken))
}

func (b *GinBinding) HandleUpdateNote(ctx *gin.Context) {
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

	csrfToken := csrf.GetToken(ctx)

	sendJSONOrHTML(ctx, http.StatusOK, &gin.H{"note": updatedNote}, components.NoteListItem(updatedNote, csrfToken))
}

func (b *GinBinding) HandleDeleteNote(ctx *gin.Context) {
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

	notes, err := b.services.GetAllNotes(ctx.Request.Context())
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	csrfToken := csrf.GetToken(ctx)

	sendJSONOrHTML(ctx, http.StatusOK, &gin.H{}, components.NoteList(notes, csrfToken))
}
