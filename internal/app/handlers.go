package app

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/EwvwGeN/medods_assignment/internal/domain/models"
	"github.com/gorilla/mux"
)

func (s *server) RegisterUser() http.HandlerFunc {
	log := s.log.With(slog.String("handler", "register_user"))
	return func(w http.ResponseWriter, r *http.Request) {
		req := models.RegisterRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("cant decode request body")
			http.Error(w, "error while decoidng response object", http.StatusBadRequest)
			return
		}
		if req.Email == "" {
			log.Info("empty email in request")
			http.Error(w, "email cant be empty", http.StatusBadRequest)
			return
		}
		uuid, err := s.auth.RegisterUser(req.Email)
		if err != nil {
			log.Info("error while registration", slog.String("email", req.Email), slog.String("error", err.Error()))
			http.Error(w, "error while registration", http.StatusInternalServerError)
			return
		}
		res := &models.RegisterResponse{
			UUID: uuid,
		}
		jsonRes, err := json.Marshal(&res)
		if err != nil {
			log.Error("cant marshal response", slog.Any("response", res), slog.String("error", err.Error()))
			http.Error(w, "error while registration", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonRes)
	}
}

func (s *server) CreateTokenPair() http.HandlerFunc {
	log := s.log.With(slog.String("handler", "create_token_pair"))
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		uuid, ok := vars["uuid"]
		if !ok {
			log.Warn("cant get uuid from params", slog.Any("vars", vars))
			http.Error(w, "cant get uuid", http.StatusBadRequest)
			return
		}
		token, refresh, err := s.auth.CreateTokenPair(uuid)
		if err != nil {
			log.Info("cant create token pair", slog.String("uuid", uuid), slog.String("error", err.Error()))
			http.Error(w, "cant create token pair", http.StatusInternalServerError)
			return
		}
		res := &models.TokenPair{
			AccessToken: token,
			RefreshToken: refresh,
		}
		jsonRes, err := json.Marshal(&res)
		if err != nil {
			log.Error("cant marshal response", slog.Any("response", res), slog.String("error", err.Error()))
			http.Error(w, "error while creating token pair", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonRes)
	}
}

func (s *server) RefrashToken() http.HandlerFunc {
	log := s.log.With(slog.String("handler", "refresh_token"))
	return func(w http.ResponseWriter, r *http.Request) {
		req := models.RefreshRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("cant decode request body")
			http.Error(w, "error while decoidng response object", http.StatusBadRequest)
			return
		}
		if req.RefreshToken == "" {
			log.Info("empty refresh token in request")
			http.Error(w, "refresh token cant be empty", http.StatusBadRequest)
			return
		}
		token, refresh, err := s.auth.RefreshToken(req.RefreshToken)
		if err != nil {
			log.Info("cant refresh token", slog.String("error", err.Error()))
			http.Error(w, "cant refresh token", http.StatusInternalServerError)
			return
		}
		res := &models.RefreshResponse{
			TokenPair: models.TokenPair{
				AccessToken: token,
				RefreshToken: refresh,
		}}
		jsonRes, err := json.Marshal(&res)
		if err != nil {
			log.Error("cant marshal response", slog.Any("response", res), slog.String("error", err.Error()))
			http.Error(w, "error refreshing token", http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonRes)
	}
}