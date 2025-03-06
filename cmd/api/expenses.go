package main

import (
	"encoding/json"
	"net/http"

	"github.com/livin21/billmate/internal/store"
)

func (app *application) createExpenseHandler(w http.ResponseWriter, r *http.Request) {
	var expense store.Expense
	err := json.NewDecoder(r.Body).Decode(&expense)
	if err != nil {
		app.serverError(w, err)
		return
	}
	if expense.Title == "" || expense.Amount == 0 {
		app.badRequest(w, "title and amount are required")
		return
	}
	parsedUserId, parseErr:= getUserIdFromContext(r.Context())
	if parseErr != nil {
		app.serverError(w, parseErr)
		return
	}
	expense.UserId = parsedUserId
	err = app.store.Expenses.Create(r.Context(), &expense)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJson(w, http.StatusCreated, expense)
}

func (app *application) listExpensesHandler(w http.ResponseWriter, r *http.Request) {
	var expenses []store.Expense
	var err error
	if r.Context().Value("role") == "admin" {
		expenses, err = app.store.Expenses.List(r.Context())
	} else {
		parsedUserId, parseErr:= getUserIdFromContext(r.Context())
		if parseErr != nil {
			app.serverError(w, parseErr)
			return
		}
		expenses, err = app.store.Expenses.ListByUser(r.Context(), parsedUserId)
	}
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJson(w, http.StatusOK, expenses)
}