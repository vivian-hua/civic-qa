package context

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/team-ravl/civic-qa/services/form/internal/generator"
	"github.com/team-ravl/civic-qa/services/form/internal/model"
	"github.com/team-ravl/civic-qa/services/form/internal/repository/form"
)

// HandleGetForm GET /form/{formID}
// Allows a user to get the HTML for a specific form (meant for iFrame)
func (ctx *Context) HandleGetForm(w http.ResponseWriter, r *http.Request) {
	// parse the URL parameter
	vars := mux.Vars(r)
	formID, err := strconv.ParseUint(vars["formID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid Form ID", http.StatusBadRequest)
		return
	}

	// retrieve the form
	formData, err := ctx.FormStore.GetByID(uint(formID))
	if err != nil {
		if err == form.ErrFormNotFound {
			http.Error(w, "Form Not Found", http.StatusNotFound)
			return
		}
		log.Printf("Error retrieving form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// generate the HTML
	formHTML, err := generator.GetHTML(formData)
	if err != nil {
		log.Printf("Error generating HTML: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// return the html
	w.Header().Set("content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(formHTML)
	if err != nil {
		log.Printf("Error writing: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// HandlePostForm POST /form/{formID}
// Allows a user to respond to a specific form
func (ctx *Context) HandlePostForm(w http.ResponseWriter, r *http.Request) {
	// parse the URL parameter
	vars := mux.Vars(r)
	formID, err := strconv.ParseUint(vars["formID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid Form ID", http.StatusBadRequest)
		return
	}

	// parse the form
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid Form Data", http.StatusBadRequest)
		return
	}

	// create the response
	responseFormData := r.Form
	response := &model.FormResponse{
		Name:         responseFormData.Get("name"),
		EmailAddress: responseFormData.Get("email"),
		InquiryType:  responseFormData.Get("inquiryType"),
		Subject:      responseFormData.Get("subject"),
		Body:         responseFormData.Get("body"),
		CreatedAt:    time.Now().UTC(),
		Active:       true,
		FormID:       uint(formID),
	}

	// respond
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Thank you!"))

	// get tags from analytics
	tagsStrings, err := ctx.Analytics.GetKeyPhrases(response)
	if err != nil {
		log.Printf("Error calling analytics: %v", err)
		tagsStrings = make([]string, 0)
	}

	// make tag structs
	tags := make([]model.Tag, len(tagsStrings))
	for i := range tagsStrings {
		tags[i] = model.Tag{Value: tagsStrings[i]}
	}

	// insert the responses
	response.Tags = tags
	// store the response
	err = ctx.ResponseStore.Create(response)
	if err != nil {
		log.Printf("Error storing response: %v", err)
		return
	}

}
