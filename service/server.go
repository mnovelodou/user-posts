package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"user_posts/business_logic"

	"github.com/gorilla/mux"
)

type Server struct {
	logic *business_logic.UserPostLogic
	srv   *http.Server
}

func NewServer(logic *business_logic.UserPostLogic) *Server {
	return &Server{logic: logic}
}

// Start Initialize the server and start listening
func (s *Server) Start(addr string) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/v1/user-posts/{id:[0-9]+}", s.GetUserPosts).Methods("GET")
	s.srv = &http.Server{
		Addr: addr,
		Handler: router,
	}

	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
}

// GetUserPosts process all requests from path /v1/user-posts/{id}
func (s *Server) GetUserPosts(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	IDStr := vars["id"] // this must be in vars map if we reached this point
	ID, _:= strconv.Atoi(IDStr) // this must be numeric if we reached this point
	user, err := s.logic.GetUserPost(request.Context(), ID)
	if err != nil {
		if _, ok := err.(*business_logic.IDNotFound); ok {
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (s *Server) Shutdown(ctx context.Context) {
	s.srv.Shutdown(ctx)
}