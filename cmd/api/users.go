package main

import (
	"encoding/json"
	"net/http"

	"github.com/livin21/billmate/internal/store"
	"github.com/livin21/billmate/internal/util"
)

func (app *application) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.store.Users.List(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJson(w, http.StatusOK, users)
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user store.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		app.serverError(w, err)
		return
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