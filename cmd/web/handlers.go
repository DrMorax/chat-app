package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"chatapp/internal/models"
	"chatapp/internal/validator"
	"chatapp/ui/components"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var buf bytes.Buffer

	err := components.HomeTemplate().Render(ctx, &buf)
	if err != nil {
		http.Error(w, "Error rendering the HTML response", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

type SignupJSON struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	validator.Validator
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	var form SignupJSON

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		app.serverError(w, err)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		jsonErrors := map[string]interface{}{
			"non_field_errors": form.NonFieldErrors,
			"field_errors":     form.FieldErrors,
		}

		headers := make(http.Header)
		headers.Set("Content-Type", "application/json")
		app.writeJSON(w, http.StatusUnprocessableEntity, jsonErrors, headers)

		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			jsonErrors := map[string]interface{}{
				"non_field_errors": form.NonFieldErrors,
				"field_errors":     form.FieldErrors,
			}
			headers := make(http.Header)
			headers.Set("Content-Type", "application/json")
			app.writeJSON(w, http.StatusUnprocessableEntity, jsonErrors, headers)
			return
		} else {
			app.serverError(w, err)
			return
		}
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type LoginJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	validator.Validator
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	var form LoginJSON

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		jsonErrors := map[string]interface{}{
			"non_field_errors": form.NonFieldErrors,
			"field_errors":     form.FieldErrors,
		}

		headers := make(http.Header)
		headers.Set("Content-Type", "application/json")
		app.writeJSON(w, http.StatusUnprocessableEntity, jsonErrors, headers)

		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("email", "Email or password is incorrect")
			jsonErrors := map[string]interface{}{
				"non_field_errors": form.NonFieldErrors,
				"field_errors":     form.FieldErrors,
			}
			headers := make(http.Header)
			headers.Set("Content-Type", "application/json")
			app.writeJSON(w, http.StatusUnprocessableEntity, jsonErrors, headers)
			return
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
