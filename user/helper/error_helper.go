package helper

import (
	"encoding/json"
	"net/http"
)

func (e *ErrorJsons) AddMessageAndDataEntry(message, data string) {
	errJson := ErrorJson{}
	errJson.SetMessage(message)
	errJson.SetData(data)
	e.AddErrorJson(errJson)
}

func (e *ErrorJsons) AddMessageEntry(message string) {
	errJson := ErrorJson{}
	errJson.SetMessage(message)
	e.AddErrorJson(errJson)
}

func (e *ErrorJsons) JsonError(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	b, err := json.Marshal(e)
	if err != nil {
		ServerError(w)
		return
	}
	http.Error(w, string(b), code)
}

func (e *ErrorJson) SetMessage(message string) {
	e.Message = message
}

func (e *ErrorJson) SetData(data string) {
	e.Data = data
}

func (e *ErrorJsons) AddErrorJson(errorJson ErrorJson) {
	e.Errors = append(e.Errors, errorJson)
}

type ErrorJson struct {
	Message string `json:"message"` // no space between json: and value
	Data    string `json:"data"`    // optional
}

type ErrorJsons struct {
	Errors []ErrorJson `json:"errors"`
}

func ServerError(w http.ResponseWriter) {
	res := ErrorJsons{}
	res.AddMessageEntry(http.StatusText(http.StatusInternalServerError))
	res.JsonError(w, http.StatusInternalServerError)
	return
}

func ClientError(w http.ResponseWriter) {
	res := ErrorJsons{}
	res.AddMessageEntry(http.StatusText(http.StatusBadRequest))
	res.JsonError(w, http.StatusBadRequest)
	return
}

func NotFound(w http.ResponseWriter) {
	res := ErrorJsons{}
	res.AddMessageEntry(http.StatusText(http.StatusNotFound))
	res.JsonError(w, http.StatusNotFound)
	return
}
