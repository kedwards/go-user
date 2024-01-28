package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"strings"
	// "github.com/go-chi/chi/v5"

	"github.com/kedwards/go-user/internal/models"
)

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	Version string `json:"version,omitempty"`
	ID      int    `json:"id,omitempty"`
}

type user struct {
	Email		 string `json:"email"`
	FirstName  string `json:"firstName"`
  LastName   string `json:"lastName"`
	Role			 string `json:"role"`
	UserName   string `json:"username"`
}

func (app *application) GetHealthStatus(w http.ResponseWriter, r *http.Request) {
  j := jsonResponse{
	  OK:      true,
		Message: fmt.Sprintf("User management service is healthy"),
		Version: app.version,
		Content: "",
	}

	out, err := json.MarshalIndent(j, "", "   ")
	if err != nil {
		app.errorLog.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// GetAllUsers lists all users 
func (app *application) GetAllUsers(w http.ResponseWriter, r *http.Request) {
  allUsers, err := app.DB.GetUsers()
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	var resp struct {
		Users []*models.User `json:"users"`
	}

	resp.Users = allUsers

	app.writeJSON(w, http.StatusOK, resp)
}

// CreateUser creates a new user
func (app *application) CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.User

	err := app.readJSON(w, r, &newUser)
  if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// save to database
	err = app.DB.CreateUser(newUser)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// send response
	var payload struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
	}

	payload.Error = false
	payload.Message = fmt.Sprintf("User %s has been created successfuil", newUser.Email)

	_ = app.writeJSON(w, http.StatusOK, payload, nil)
}

// UpdateUser updates a user
func (app *application) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var updatedUser models.User

	err := app.readJSON(w, r, &updatedUser)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	err = app.DB.UpdateUser(updatedUser)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// send response
	var payload struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
	}

	payload.Error = false
	payload.Message = fmt.Sprintf("User %s has been updated successfuil", updatedUser.Email)

	_ = app.writeJSON(w, http.StatusOK, payload, nil)
}

// DeleteUser deletes a user
func (app *application) DeleteUser(w http.ResponseWriter, r *http.Request) {
  // userEmail := chi.URLParam(r, "user")
	// changed to this, so we can test it more easily
	// split the URL up by /, and grab the 3rd element
	exploded := strings.Split(r.RequestURI, "/")
	userEmail := exploded[3]

	err := app.DB.DeleteUser(userEmail)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// send response
	var payload struct {
		Error   bool          `json:"error"`
		Message string        `json:"message"`
	}

	payload.Error = false
	payload.Message = fmt.Sprintf("User %s has been deleted successfuil", userEmail)

	_ = app.writeJSON(w, http.StatusOK, payload, nil)
}
