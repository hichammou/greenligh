package main

import (
	"errors"
	"net/http"

	"greenlight.hichammou/internal/data"
	"greenlight.hichammou/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateUser(v, user)
	if !v.Valide() {
		app.faildValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateEmail) {
			v.AddError("email", "a user with this email address already exists")
			app.faildValidationResponse(w, r, v.Errors)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.logger.PrintInfo("email", map[string]string{"email": user.Email})

	// Use the background helper to execute an anonymous function that sends the welcome email
	app.background(func() {
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			app.logger.PrintError(err, nil)
			return
		}

	})

	success := "your registration completed successfully"
	err = app.writeJSON(w, http.StatusAccepted, envelope{"message": success, "user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
