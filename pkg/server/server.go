package server

import (
	"context"
	"errors"
	"event/pkg/config"
	"event/pkg/models"
	"event/pkg/repository"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	notEnoughFieldsInQueryError = errors.New("not enough fields in query")
)

type Server struct {
	server *http.Server
	events repository.EventsDb
	routes []*Route
}

func NewServer(events repository.EventsDb) *Server {
	return &Server{events: events}
}

func (s *Server) InitServer(cfg *config.Config) {
	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: s,
	}
	s.routes = []*Route{
		NewRoute("POST", CreateEventPostParser, s.loggingMiddleware(s.handleCreateEvent)),
		NewRoute("POST", UpdateEventPostParser, s.loggingMiddleware(s.handleUpdateEvent)),
		NewRoute("POST", DeleteEventPostParser, s.loggingMiddleware(s.handleDeleteEvent)),
	}
}

func (s *Server) Run() error {
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range s.routes {
		fmt.Println(r.URL.RawQuery)
		matches, err := route.requestParser(r)
		if err != nil {
			log.Println(err)
		}
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches)
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

func parseQueryParams(queryParams map[string]interface{}) (*models.Event, error) {
	var (
		userIdStr string
		name      string
		dateStr   string
		idStr     string
		ok        bool
		event     = &models.Event{}
		err       error
	)
	if userIdStr, ok = queryParams["user_id"].(string); !ok {
		return nil, notEnoughFieldsInQueryError
	}
	if idStr, ok = queryParams["id"].(string); !ok {
		return nil, notEnoughFieldsInQueryError

	}
	if name, ok = queryParams["name"].(string); !ok {
		return nil, notEnoughFieldsInQueryError
	}
	if dateStr, ok = queryParams["date"].(string); !ok {
		return nil, notEnoughFieldsInQueryError
	}

	if idStr != "" {
		if event.Id, err = strconv.ParseUint(idStr, 10, 64); err != nil {
			return nil, err
		}
	}

	if event.UserId, err = strconv.ParseUint(userIdStr, 10, 64); err != nil {
		return nil, err
	}

	event.Name = name

	if event.Date, err = time.Parse(time.RFC3339, dateStr); err != nil {
		return nil, err
	}

	return event, nil
}

func parseDeleteQueryParams(queryParams map[string]interface{}) (*models.Event, error) {
	var (
		idStr string
		ok    bool
		event = &models.Event{}
		err   error
	)
	if idStr, ok = queryParams["id"].(string); !ok {
		return nil, notEnoughFieldsInQueryError

	}
	if event.Id, err = strconv.ParseUint(idStr, 10, 64); err != nil {
		return nil, err
	}
	return event, err
}

func (s *Server) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		next.ServeHTTP(writer, request)
		log.Printf("%s %s %s", request.Method, request.RequestURI, time.Since(start))
	}
}

func (s *Server) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	queryParams := r.Context().Value(ctxKey{}).(map[string]interface{})
	event, err := parseQueryParams(queryParams)
	if err != nil {
		log.Println(err)
		return
	}
	if err = s.events.AddEvent(event); err != nil {
		log.Println(err)
		return
	}
	if _, err = w.Write([]byte(fmt.Sprintf("{ \"result\":\"event with id=%d created\"}",
		event.Id))); err != nil { //TODO cfg
		log.Println(err)
	}
}

func (s *Server) handleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	queryParams := r.Context().Value(ctxKey{}).(map[string]interface{})
	event, err := parseQueryParams(queryParams)
	if err = s.events.UpdateEvent(event); err != nil {
		log.Println(err) //отправлять ошибку
		return
	}
	if _, err = w.Write([]byte(fmt.Sprintf("{ \"result\":\"event with id=%d updated\"}",
		event.Id))); err != nil { //TODO cfg
		log.Println(err)
	}
}

func (s *Server) handleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	queryParams := r.Context().Value(ctxKey{}).(map[string]interface{})
	event, err := parseDeleteQueryParams(queryParams)
	if err = s.events.DeleteEvent(event.Id); err != nil {
		log.Println(err) //отправлять ошибку
		return
	}
	if _, err = w.Write([]byte(fmt.Sprintf("{ \"result\":\"event with id=%d deleted\"}",
		event.Id))); err != nil { //TODO cfg
		log.Println(err)
	}
}
