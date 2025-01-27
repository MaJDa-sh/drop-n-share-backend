package controllers

import (
	"drop-n-share/internal/middleware"
	"drop-n-share/internal/models"
	"drop-n-share/internal/services"
	"drop-n-share/internal/views"
	"encoding/json"
	"net/http"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	var signParams models.SignParams

	if err := json.NewDecoder(r.Body).Decode(&signParams); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid input"})
		return
	}

	result, err := uc.userService.SignUp(&signParams)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	token, err := middleware.GenerateJWT(result.ID, result.Username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	userView := views.NewUserView(result.ID, result.Username, token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userView)
}

func (uc *UserController) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	var signParams models.SignParams

	if err := json.NewDecoder(r.Body).Decode(&signParams); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid input"})
		return
	}

	result, err := uc.userService.SignIn(&signParams)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: err.Error()})
		return
	}

	token, err := middleware.GenerateJWT(result.ID, result.Username)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	userView := views.NewUserView(result.ID, result.Username, token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userView)
}

func (uc *UserController) HandleGetUserByJWT(w http.ResponseWriter, r *http.Request) {

	claims, err := middleware.GetClaimsFromToken(r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid token"})
		return
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid user ID in token"})
		return
	}

	user, err := uc.userService.GetUserByID(int(userID))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "User not found"})
		return
	}

	userView := views.ProfileView(user.ID, user.Username, user.Files)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userView)
}
