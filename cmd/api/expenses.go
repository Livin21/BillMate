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
	err = app.store.Expenses.Create(r.Context(), &expense)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJson(w, http.StatusCreated, expense)
}

func (app *application) listExpensesHandler(w http.ResponseWriter, r *http.Request) {
	expenses, err := app.store.Expenses.List(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.writeJson(w, http.StatusOK, expenses)
}