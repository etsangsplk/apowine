package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aporeto-inc/apowine/source/mongodb-lib"
	"github.com/aporeto-inc/apowine/source/server/internal/auth"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

// Server holds database object
type Server struct {
	mongodb       *mongodb.MongoDB
	newConnection bool
	host          []string
	database      string
	collection    string
	auth          *auth.Auth
	session       *sessions.Session
}

// NewServer creates a new server handler
func NewServer(mongo *mongodb.MongoDB, isNewConnection bool, host []string, database string, collection string, auth *auth.Auth) *Server {
	zap.L().Info("Creating a new server handler")

	return &Server{
		mongodb:       mongo,
		newConnection: isNewConnection,
		host:          host,
		database:      database,
		collection:    collection,
		auth:          auth,
	}
}

// AllDrinks returns drinks based on type in JSON format
func (s *Server) AllDrinks(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	s.checkIfUserAuthenticated(w, r)

	if s.newConnection {
		s.createDatabaseSession()
	}

	//Extracting the endpoints from URL
	drinkName := strings.SplitAfter(r.URL.RequestURI(), "/")
	decoder := json.NewDecoder(r.Body)
	data, err := s.mongodb.Read(decoder, drinkName[1], false)
	if err != nil {
		zap.L().Error("error reading data from database", zap.Error(err))
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		zap.L().Error("error in json output", zap.Error(err))
	}
}

func (s *Server) checkIfUserAuthenticated(w http.ResponseWriter, r *http.Request) {
	cookie := s.auth.GetCookie()

	s.session, _ = cookie.GetCookieStore().Get(r, "sessions")

	// Check if user is authenticated
	if auth, ok := s.session.Values["authenticated"].(bool); !ok || !auth {

		if e := s.session.Save(r, w); e != nil {
			zap.L().Error("Error in saving the session ", zap.Error(e))
		}
		s.session.Values["redirectURL"] = r.URL.String()
		fmt.Println("URL", r.URL.String())
		if err := s.session.Save(r, w); err != nil {
			zap.Error(err)
		}

		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
}

// RandomDrink returns random drink based on type in JSON format
func (s *Server) RandomDrink(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	s.checkIfUserAuthenticated(w, r)

	if s.newConnection {
		s.createDatabaseSession()
	}

	endpoint := strings.SplitAfter(r.URL.RequestURI(), "/")
	drinkName := strings.Replace(endpoint[1], "/", "", -1)
	decoder := json.NewDecoder(r.Body)
	data, err := s.mongodb.Read(decoder, drinkName, true)
	if err != nil {
		zap.L().Error("error reading data from database", zap.Error(err))
	}
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		zap.L().Error("error in json output", zap.Error(err))
	}

}
func (s *Server) createDatabaseSession() {
	mongodb, err := mongodb.NewMongoSession(s.host, "", "", s.database, s.collection)
	if err != nil {
		zap.L().Error("error creating a session", zap.Error(err))
	}
	s.mongodb = mongodb
}

// FindDrinkEndpoint finds If a drink is available in the database given ID in the URL
// Writes JSON
func (s *Server) FindDrinkEndpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	s.checkIfUserAuthenticated(w, r)

	if s.newConnection {
		s.createDatabaseSession()
	}

	endpoint := strings.SplitAfter(r.URL.RequestURI(), "/")
	drinkName := strings.Replace(endpoint[1], "/", "", -1)
	params := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	data, err := s.mongodb.ReadByID(params["id"], decoder, drinkName)
	if err != nil {
		zap.L().Error("error reading data from database", zap.Error(err))
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		zap.L().Error("error in json output", zap.Error(err))
	}
}

// CreateDrinkEndPoint creates a drink (beer or wine) with an ID to use in MongoDB
func (s *Server) CreateDrinkEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	s.checkIfUserAuthenticated(w, r)

	if s.newConnection {
		s.createDatabaseSession()
	}

	drinkName := strings.SplitAfter(r.URL.RequestURI(), "/")
	decoder := json.NewDecoder(r.Body)
	err := s.mongodb.Insert(decoder, drinkName[1])
	if err != nil {
		zap.L().Error("error inserting data from database", zap.Error(err))
	}
}

// UpdateDrinkEndPoint updates an existing drink property given its ID
func (s *Server) UpdateDrinkEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	s.checkIfUserAuthenticated(w, r)

	if s.newConnection {
		s.createDatabaseSession()
	}

	drinkName := strings.SplitAfter(r.URL.RequestURI(), "/")
	decoder := json.NewDecoder(r.Body)

	err := s.mongodb.Update(decoder, drinkName[1])
	if err != nil {
		zap.L().Error("error updating data in database", zap.Error(err))
	}
}

// DeleteDrinkEndPoint deletes a drink given ID and its type
func (s *Server) DeleteDrinkEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	s.checkIfUserAuthenticated(w, r)

	if s.newConnection {
		s.createDatabaseSession()
	}

	params := mux.Vars(r)

	err := s.mongodb.Delete(params["id"])
	if err != nil {
		zap.L().Error("error deleting data from database", zap.Error(err))
	}
}
