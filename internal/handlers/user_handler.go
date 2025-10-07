/*
TODO:
1. Add validation method for username and password.
2. Improve response messages

*/

package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/nabsk911/code-snippet-organizer/internal/store"
	"github.com/nabsk911/code-snippet-organizer/internal/utils"
)

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (uh *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req authRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("Failed to decode register request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload!"})
		return
	}

	if req.Username == "" || req.Password == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Username and password are required"})
		return
	}

	password_hash, err := utils.SetPasswordHash(req.Password)

	if err != nil {
		uh.logger.Printf("Failed to hash password: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error!"})
		return
	}

	user := &store.User{
		Username:     req.Username,
		PasswordHash: password_hash,
	}

	err = uh.userStore.CreateUser(user)
	if err != nil {
		uh.logger.Printf("Failed to create user in store: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Failed to create user!"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "User created successfully!"})
}

func (uh *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req authRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		uh.logger.Printf("Failed to decode login request: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid request payload!"})
		return
	}

	if req.Username == "" || req.Password == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Username and password are required"})
		return
	}

	user, err := uh.userStore.GetUserByUsername(req.Username)

	if err != nil {
		uh.logger.Printf("Failed to retrieve user %s: %v", req.Username, err)
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid credentials!"})
		return
	}

	passwordMatches, err := utils.CheckPasswordHash(req.Password, user.PasswordHash)

	if err != nil {
		uh.logger.Printf("Error checking password hash: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error!"})
		return
	}

	if !passwordMatches {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "Invalid credentials!"})
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		uh.logger.Printf("Failed to generate token: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"token": token,
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}
