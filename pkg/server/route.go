package server

import (
	"encoding/json"
	"net/http"
)

var (
	RegexpUserId    = `([0-9]+)`
	RegexpDate      = `((19|20)\d\d-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01]))`
	RegexpEventName = `([a-z|A-Z]+)`
)

type Route struct {
	method        string
	requestParser func(req *http.Request) (map[string]interface{}, error)
	handler       http.HandlerFunc
}

func NewRoute(method string, requestParser func(req *http.Request) (map[string]interface{}, error), handler http.HandlerFunc) *Route {
	return &Route{method: method, requestParser: requestParser, handler: handler}
}

type ctxKey struct{}

func getField(r *http.Request, index int) string {
	return r.Context().Value(ctxKey{}).([]string)[index]
}

func CreateEventPostParser(req *http.Request) (map[string]interface{}, error) {
	if req.URL.Path != "/create_event" {
		return nil, nil
	}
	return ParseRequestBody(req)

}

func DeleteEventPostParser(req *http.Request) (map[string]interface{}, error) {
	if req.URL.Path != "/delete_event" {
		return nil, nil
	}
	return ParseRequestBody(req)

}

func UpdateEventPostParser(req *http.Request) (map[string]interface{}, error) {
	if req.URL.Path != "/update_event" {
		return nil, nil
	}
	return ParseRequestBody(req)
}

func ParseRequestBody(r *http.Request) (map[string]interface{}, error) {
	bodyBytes := make([]byte, 1024)
	n, _ := r.Body.Read(bodyBytes)
	bodyBytes = bodyBytes[:n]
	params := make(map[string]interface{})
	if err := json.Unmarshal(bodyBytes, &params); err != nil {
		return nil, err
	}
	return params, nil
}
