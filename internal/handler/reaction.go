package handler

import (
	"errors"
	"net/http"
	"strconv"

	"forum/internal/models"
)

func (h *Handler) createReaction(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		user := r.Context().Value(ctxKey).(models.User)

		if err := r.ParseForm(); err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
		entityID, err := strconv.Atoi(r.PostForm.Get("id"))
		object := r.PostForm.Get("object")
		if err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
		var reactionType string
		if val, ok := r.Form["like"]; ok {
			if val[0] != "like" {
				h.ErrorPage(w, http.StatusBadRequest, err)
				return
			}
			reactionType = val[0]
		} else if val, ok := r.Form["dislike"]; ok {
			if val[0] != "dislike" {
				h.ErrorPage(w, http.StatusBadRequest, err)
				return
			}
			reactionType = val[0]
		} else {
			h.ErrorPage(w, http.StatusBadRequest, err)
			return
		}
		reaction := models.Reaction{
			AuthorID: user.ID,
			Type:     reactionType,
		}
		if object == "post" {
			_, err := h.services.GetPostById(entityID)
			if err != nil {
				h.ErrorPage(w, http.StatusBadRequest, err)
				return
			}
			reaction.PostID = entityID
		} else if object == "comment" {
			_, err := h.services.Commentary.GetCommentByID(entityID)
			if err != nil {
				h.ErrorPage(w, http.StatusBadRequest, err)
				return
			}
			reaction.CommentID = entityID
		} else {
			h.ErrorPage(w, http.StatusBadRequest, err)
			return
		}

		if err := h.services.Reaction.CreateReaction(reaction); err != nil {
			h.ErrorPage(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		h.ErrorPage(w, http.StatusMethodNotAllowed, errors.New("error: method not allowed"))
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
