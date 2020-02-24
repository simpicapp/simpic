package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var validate = validator.New()

// writeJSON marshals the given payload and writes it to the ResponseWriter.
func writeJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

// writeError writes out a JSON error response.
func writeError(w http.ResponseWriter, code int, error string) {
	type Error struct {
		Error string `json:"error"`
	}

	writeJSON(w, code, Error{Error: error})
}

func paginate(w http.ResponseWriter, r *http.Request, generator func(offset, count int) (interface{}, error)) {
	offset := 0
	param, ok := r.URL.Query()["offset"]
	if ok && len(param) > 0 {
		offset, _ = strconv.Atoi(param[0])
	}

	res, err := generator(offset, 100)
	if err != nil {
		log.Printf("unable to paginate: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

// bind attempts to decode the request body as JSON and bind it to the given struct.
// If the body cannot be decoded or the struct fails validation, an appropriate
// response will be written and the func will return false.
func bind(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return false
	}

	if err := validate.StructCtx(r.Context(), v); err != nil {
		writeError(w, http.StatusUnprocessableEntity, buildError(err.(validator.ValidationErrors)))
		return false
	}

	return true
}

// buildError creates a user-friendly error message from a collection of validation errors.
func buildError(errors validator.ValidationErrors) string {
	var str strings.Builder
	for _, err := range errors {
		if str.Len() > 0 {
			str.WriteRune('\n')
		}
		str.WriteString(err.Field())
		str.WriteRune(' ')
		str.WriteString(describeTag(err.Tag(), err.Param()))
	}
	return str.String()
}

// describeTag converts a validation tag to a user-friendly description of the failure mode.
func describeTag(tag string, value string) string {
	if tag == "required" {
		return "is required"
	} else if tag == "min" {
		return fmt.Sprintf("must have a length of at least %s", value)
	} else if tag == "max" {
		return fmt.Sprintf("must have a length of at most %s", value)
	} else {
		log.Printf("Don't know how to describe validation tag %s\n", tag)
		return "is invalid"
	}
}
