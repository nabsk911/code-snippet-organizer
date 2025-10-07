package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/nabsk911/code-snippet-organizer/internal/store"
	"github.com/nabsk911/code-snippet-organizer/internal/utils"
)

type snippetRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Code        string `json:"code"`
	Language    string `json:"language"`
}

type SnippetHandler struct {
	snippetStore store.SnippetStore
	logger       *log.Logger
}

func NewSnippetHandler(logger *log.Logger, snippetStore store.SnippetStore) *SnippetHandler {
	return &SnippetHandler{
		logger:       logger,
		snippetStore: snippetStore,
	}
}

func (sh *SnippetHandler) HandleCreateSnippet(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")

	var req snippetRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		sh.logger.Printf("Error decoding snippet %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload!"})
		return
	}

	snippet := store.Snippet{
		Title:       req.Title,
		Description: req.Description,
		Code:        req.Code,
		Language:    req.Language,
		UserID:      userID.(int),
	}

	createdSnippet, err := sh.snippetStore.CreateSnippet(&snippet)

	if err != nil {
		sh.logger.Printf("Error creating snippet %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error!"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"data": createdSnippet, "message": "Snippet created successfully!"})
}

func (sh *SnippetHandler) HandleGetSnippetsByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")

	snippets, err := sh.snippetStore.GetSnippetsByUserID(userID.(int))

	if err != nil {
		sh.logger.Printf("Error getting snippets %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error!"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": snippets})
}

func (sh *SnippetHandler) HandleDeleteSnippet(w http.ResponseWriter, r *http.Request) {
	snippetID := r.PathValue("id")

	snippetIDInt, err := strconv.Atoi(snippetID)

	if err != nil {
		sh.logger.Printf("Error parsing snippet ID %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid snippet ID!"})
		return
	}

	err = sh.snippetStore.DeleteSnippet(snippetIDInt)

	if err != nil {
		sh.logger.Printf("Error deleting snippet %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error!"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "Snippet deleted successfully!"})
}

func (sh *SnippetHandler) HandleUpdateSnippet(w http.ResponseWriter, r *http.Request) {
	var req snippetRequest

	snippetID := r.PathValue("id")

	snippetIDInt, err := strconv.Atoi(snippetID)

	if err != nil {
		sh.logger.Printf("Error parsing snippet ID %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid snippet ID!"})
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		sh.logger.Printf("Error decoding snippet %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload!"})
		return
	}

	snippet := store.Snippet{
		ID:          snippetIDInt,
		Title:       req.Title,
		Description: req.Description,
		Code:        req.Code,
		Language:    req.Language,
	}

	updatedSnippet, err := sh.snippetStore.UpdateSnippet(&snippet)

	if err != nil {
		sh.logger.Printf("Error updating snippet %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error!"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": updatedSnippet, "message": "Snippet updated successfully!"})
}

func (sh *SnippetHandler) HandleSearchSnippets(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	title := r.URL.Query().Get("title")
	language := r.URL.Query().Get("language")

	snippets, err := sh.snippetStore.SearchSnippets(title, language, userID.(int))
	if err != nil {
		sh.logger.Printf("Error searching snippets %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error!"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": snippets})
}
