package prism

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

//////////////////////////////////////////////////////////////////
// VAR SYSTEM (originaly written in Go Programming Blueprints)
//////////////////////////////////////////////////////////////////

var vars map[*http.Request]map[string]interface{}
var varsLock sync.RWMutex

// GetVar gets the value of the key for the specified http.Request.
func GetVar(r *http.Request, key string) interface{} {

	varsLock.RLock()
	value := vars[r][key]
	varsLock.RUnlock()
	return value
}

// SetVar sets the key to the value for the specified http.Request.
func SetVar(r *http.Request, key string, value interface{}) {

	varsLock.Lock()
	vars[r][key] = value
	varsLock.Unlock()
}

// OpenVars opens the vars for the specified http.Request.
// Must be called before GetVar or SetVar is called for each
// request.
func OpenVars(r *http.Request) {

	varsLock.Lock()
	if vars == nil {
		vars = map[*http.Request]map[string]interface{}{}
	}
	vars[r] = map[string]interface{}{}
	varsLock.Unlock()
}

// CloseVars closes the vars for the specified
// http.Request.
// Must be called when all var activity is completed to
// clean up any used memory.
func CloseVars(r *http.Request) {

	varsLock.Lock()
	delete(vars, r)
	varsLock.Unlock()
}

//////////////////////////////////////////////////////////////////
// Encode/Decode Utility (originaly written in Go Programming Blueprints)
//////////////////////////////////////////////////////////////////

func DecodeBody(r *http.Request, v interface{}) error {

	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func EncodeBody(res http.ResponseWriter, req *http.Request, v interface{}) error {

	return json.NewEncoder(res).Encode(v)
}

//////////////////////////////////////////////////////////////////
// Response Utility (originaly written in Go Programming Blueprints)
//////////////////////////////////////////////////////////////////

func Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {

	w.WriteHeader(status)
	if data != nil {
		EncodeBody(w, r, data)
	}
}
func RespondErr(w http.ResponseWriter, r *http.Request, status int, args ...interface{}) {

	Respond(w, r, status, map[string]interface{}{
		"error": map[string]interface{}{
			"message": fmt.Sprint(args...),
		},
	})
}
func RespondHTTPErr(w http.ResponseWriter, r *http.Request, status int) {

	RespondErr(w, r, status, http.StatusText(status))
}
