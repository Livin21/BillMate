package main

import (
	"encoding/json"
	"net/http"

	"github.com/livin21/billmate/internal/store"
	"github.com/livin21/billmate/internal/util"
)

func (app *application) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("role") != "admin" {
		app.unAuthorized(w)
		return
	}
	users, err := app.store.Users.List(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJson(w, http.StatusOK, users)
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("role") != "admin" {
		app.unAuthorized(w)
		return
	}
	var user store.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if user.Email == "" || user.Password == "" {
		app.badRequest(w, "email and password are required")
		return
	}
	if user.Role != "" {
		if user.Role != "admin" && user.Role != "user" {
			app.badRequest(w, "role can only be admin or user")
			return
		}
	}
	err = app.store.Users.Create(r.Context(), &user)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJson(w, http.StatusCreated, &util.MessageResponse{
		Message: "User created successfully",
		Status: http.StatusCreated,
	})
}