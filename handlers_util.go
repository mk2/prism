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
func GetVar(req *http.Request, key string) interface{} {
	varsLock.RLock()
	value := vars[req][key]
	varsLock.RUnlock()
	return value
}

// SetVar sets the key to the value for the specified http.Request.
func SetVar(req *http.Request, key string, value interface{}) {
	varsLock.Lock()
	vars[req][key] = value
	varsLock.Unlock()
}

// OpenVars opens the vars for the specified http.Request.
// Must be called before GetVar or SetVar is called for each
// request.
func OpenVars(req *http.Request) {
	varsLock.Lock()
	if vars == nil {
		vars = map[*http.Request]map[string]interface{}{}
	}
	vars[req] = map[string]interface{}{}
	varsLock.Unlock()
}

// CloseVars closes the vars for the specified
// http.Request.
// Must be called when all var activity is completed to
// clean up any used memory.
func CloseVars(res *http.Request) {
	varsLock.Lock()
	delete(vars, res)
	varsLock.Unlock()
}

//////////////////////////////////////////////////////////////////
// Encode/Decode Utility (originaly written in Go Programming Blueprints)
//////////////////////////////////////////////////////////////////

func DecodeBody(req *http.Request, v interface{}) error {
	defer req.Body.Close()
	return json.NewDecoder(req.Body).Decode(v)
}

func EncodeBody(res http.ResponseWriter, req *http.Request, v interface{}) error {
	return json.NewEncoder(res).Encode(v)
}

//////////////////////////////////////////////////////////////////
// Response Utility (originaly written in Go Programming Blueprints)
//////////////////////////////////////////////////////////////////

func Respond(res http.ResponseWriter, req *http.Request, status int, data interface{}) {
	res.WriteHeader(status)
	if data != nil {
		EncodeBody(res, req, data)
	}
}
func RespondErr(res http.ResponseWriter, req *http.Request, status int, args ...interface{}) {
	Respond(res, req, status, map[string]interface{}{
		"error": map[string]interface{}{
			"message": fmt.Sprint(args...),
		},
	})
}
func RespondHTTPErr(res http.ResponseWriter, req *http.Request, status int) {
	RespondErr(res, req, status, http.StatusText(status))
}
