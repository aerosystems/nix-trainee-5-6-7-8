package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/aerosystems/nix-trainee-4/internal/models"
	"github.com/go-chi/chi/v5"
)

func (h *BaseHandler) ReadComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.commentRepo.FindAll()
	if err != nil {
		log.Println(err)
	}

	if len(*comments) == 0 {
		_ = errorJSON(w, errors.New("comments do not exist"), http.StatusNotFound)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "all comments with ID were found successfully",
		Data:    comments,
	}
	_ = writeJSON(w, http.StatusOK, payload)
}

func (h *BaseHandler) ReadComment(w http.ResponseWriter, r *http.Request) {
	stringID := chi.URLParam(r, "id")

	ID, err := strconv.Atoi(stringID)
	if err != nil {
		_ = errorJSON(w, err, http.StatusBadRequest)
	}

	comment, err := h.commentRepo.FindByID(ID)
	if err != nil {
		log.Println(err)
	}

	if *comment == (models.Comment{}) {
		_ = errorJSON(w, errors.New("comment with ID "+stringID+" does not exist"), http.StatusNotFound)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "comment with ID " + stringID + " was found successfully",
		Data:    comment,
	}
	_ = writeJSON(w, http.StatusOK, payload)
}

func (h *BaseHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	stringID := chi.URLParam(r, "id")

	ID, err := strconv.Atoi(stringID)
	if err != nil {
		_ = errorJSON(w, err, http.StatusBadRequest)
	}

	var requestPayload struct {
		PostID int    `json:"postId"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Body   string `json:"body"`
	}

	err = readJSON(w, r, &requestPayload)
	if err != nil {
		_ = errorJSON(w, err, http.StatusBadRequest)
		return
	}

	comment, err := h.commentRepo.FindByID(ID)
	if err != nil {
		log.Println(err)
	}

	if *comment != (models.Comment{}) {
		_ = errorJSON(w, errors.New("comment with ID "+stringID+" exists"), http.StatusNotFound)
		return
	}

	newComment := models.Comment{
		Id:     ID,
		PostId: requestPayload.PostID,
		Name:   requestPayload.Name,
		Email:  requestPayload.Email,
		Body:   requestPayload.Body,
	}

	err = h.commentRepo.Create(&newComment)
	if err != nil {
		log.Println(err)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "comment with ID " + stringID + " was created successfully",
		Data:    newComment,
	}
	_ = writeJSON(w, http.StatusOK, payload)
}

func (h *BaseHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	stringID := chi.URLParam(r, "id")

	ID, err := strconv.Atoi(stringID)
	if err != nil {
		_ = errorJSON(w, err, http.StatusBadRequest)
	}

	var requestPayload struct {
		PostID int    `json:"postId"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Body   string `json:"body"`
	}

	err = readJSON(w, r, &requestPayload)
	if err != nil {
		_ = errorJSON(w, err, http.StatusBadRequest)
		return
	}

	comment, err := h.commentRepo.FindByID(ID)
	if err != nil {
		log.Println(err)
	}

	if *comment == (models.Comment{}) {
		_ = errorJSON(w, errors.New("comment with ID "+stringID+"does not exists"), http.StatusNotFound)
		return
	}

	newComment := comment
	newComment.Id = ID

	if requestPayload.PostID != 0 {
		newComment.PostId = requestPayload.PostID
	}
	if requestPayload.Name != "" {
		newComment.Name = requestPayload.Name
	}
	if requestPayload.Email != "" {
		newComment.Email = requestPayload.Email
	}
	if requestPayload.Body != "" {
		newComment.Body = requestPayload.Body
	}

	err = h.commentRepo.Update(newComment)
	if err != nil {
		log.Println(err)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "comment with ID " + stringID + " was updated successfully",
		Data:    nil,
	}
	_ = writeJSON(w, http.StatusOK, payload)
}

func (h *BaseHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	stringID := chi.URLParam(r, "id")

	ID, err := strconv.Atoi(stringID)
	if err != nil {
		_ = errorJSON(w, err, http.StatusBadRequest)
	}

	comment, err := h.commentRepo.FindByID(ID)
	if err != nil {
		log.Println(err)
	}

	if *comment == (models.Comment{}) {
		_ = errorJSON(w, errors.New("comment with ID "+stringID+" does not exist"), http.StatusNotFound)
		return
	}

	err = h.commentRepo.Delete(comment)
	if err != nil {
		log.Println(err)
	}

	payload := jsonResponse{
		Error:   false,
		Message: "comment with ID " + stringID + " was deleted successfully",
		Data:    nil,
	}
	_ = writeJSON(w, http.StatusOK, payload)
}
