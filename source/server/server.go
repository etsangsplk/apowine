package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aporeto-inc/apowine/source/mongodb-lib"
	"github.com/aporeto-inc/apowine/source/server/configuration"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

// Server holds database object
type Server struct {
	mongodb          *mongodb.MongoDB
	midgardToken     string
	newConnection    bool
	host             []string
	database         string
	collection       string
	session          *sessions.Session
	midgardTokenJSON *midgardtoken
	authorizedUser   *userdetails
	unauthorizedUser *userdetails
}

type userdetails struct {
	email string
	lname string
	fname string
}

type midgardtoken struct {
	Claims struct {
		Aud  string `json:"aud"`
		Data struct {
			Email        string `json:"email"`
			FamilyName   string `json:"familyName"`
			GivenName    string `json:"givenName"`
			Name         string `json:"name"`
			Organization string `json:"organization"`
			Realm        string `json:"realm"`
		} `json:"data"`
		Exp   int    `json:"exp"`
		Iat   int    `json:"iat"`
		Iss   string `json:"iss"`
		Realm string `json:"realm"`
		Sub   string `json:"sub"`
	} `json:"claims"`
}

// NewServer creates a new server handler
func NewServer(mongo *mongodb.MongoDB, host []string, cfg *configuration.Configuration) *Server {
	zap.L().Info("Creating a new server handler")

	authorizedUser := &userdetails{
		email: cfg.AuthorizedEmail,
		fname: cfg.AuthorizedGivenName,
		lname: cfg.AuthorizedFamilyName,
	}

	unauthorizedUser := &userdetails{
		email: cfg.UnauthorizedEmail,
		fname: cfg.UnsuthorizedGivenName,
		lname: cfg.UnauthorizedFamilyName,
	}

	if cfg.MakeNewConnection && mongo != nil {
		mongo.GetSession().Close()
	}

	return &Server{
		mongodb:          mongo,
		newConnection:    cfg.MakeNewConnection,
		host:             host,
		database:         cfg.MongoDatabaseName,
		collection:       cfg.MongoCollectionName,
		authorizedUser:   authorizedUser,
		unauthorizedUser: unauthorizedUser,
	}
}

// AllDrinks returns drinks based on type in JSON format
func (s *Server) AllDrinks(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	m := s.mongodb
	var err error
	if s.newConnection {

		m, err = mongodb.NewMongoSession(s.host, "", "", s.database, s.collection)
		if err != nil {
			zap.L().Error("error creating a session", zap.Error(err))
		}
		defer m.GetSession().Close()

	}

	//Extracting the endpoints from URL
	drinkName := strings.SplitAfter(r.URL.RequestURI(), "/")
	decoder := json.NewDecoder(r.Body)
	data, err := m.Read(decoder, drinkName[1], false)
	if err != nil {
		zap.L().Error("error reading data from database", zap.Error(err))
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		zap.L().Error("error in json output", zap.Error(err))
	}
}

func (s *Server) checkIfUserAuthenticated(w http.ResponseWriter, r *http.Request) error {

	url := "https://api.console.aporeto.com/auth?token=" + s.midgardToken

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	claims, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var token midgardtoken
	json.Unmarshal(claims, &token)
	s.midgardTokenJSON = &token
	err = s.generateandValidateUserModel(s.midgardTokenJSON)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) generateandValidateUserModel(tokenJSON *midgardtoken) error {
	// Check for predefined authorizing policies
	if tokenJSON.Claims.Data.Email == s.authorizedUser.email && tokenJSON.Claims.Data.GivenName == s.authorizedUser.fname && tokenJSON.Claims.Data.FamilyName == s.authorizedUser.lname {
		return nil
	} else if tokenJSON.Claims.Data.Email == s.unauthorizedUser.email && tokenJSON.Claims.Data.GivenName == s.unauthorizedUser.fname && tokenJSON.Claims.Data.FamilyName == s.unauthorizedUser.lname {
		return fmt.Errorf(s.unauthorizedUser.fname + " is not authorized to access this resource")
	}
	// Enforcer will police
	return nil
}

// RandomDrink returns random drink based on type in JSON format
func (s *Server) RandomDrink(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	m := s.mongodb
	var err error
	if s.newConnection {

		m, err = mongodb.NewMongoSession(s.host, "", "", s.database, s.collection)
		if err != nil {
			zap.L().Error("error creating a session", zap.Error(err))
		}
		defer m.GetSession().Close()

	}

	endpoint := strings.SplitAfter(r.URL.RequestURI(), "/")
	drinkName := strings.Replace(endpoint[1], "/", "", -1)
	decoder := json.NewDecoder(r.Body)
	data, err := m.Read(decoder, drinkName, true)
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

	err := s.checkIfUserAuthenticated(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := s.mongodb
	if s.newConnection {
		m, err = mongodb.NewMongoSession(s.host, "", "", s.database, s.collection)
		if err != nil {
			zap.L().Error("error creating a session", zap.Error(err))
		}
		defer m.GetSession().Close()
	}

	endpoint := strings.SplitAfter(r.URL.RequestURI(), "/")
	drinkName := strings.Replace(endpoint[1], "/", "", -1)
	params := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)
	data, err := m.ReadByID(params["id"], decoder, drinkName)
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

	err := s.checkIfUserAuthenticated(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := s.mongodb

	if s.newConnection {
		m, err = mongodb.NewMongoSession(s.host, "", "", s.database, s.collection)
		if err != nil {
			zap.L().Error("error creating a session", zap.Error(err))
		}
		defer m.GetSession().Close()

	}

	drinkName := strings.SplitAfter(r.URL.RequestURI(), "/")
	decoder := json.NewDecoder(r.Body)

	if err = m.Insert(decoder, drinkName[1]); err != nil {

		zap.L().Error("error inserting data from database", zap.Error(err))
	}
}

func (s *Server) GetToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bearerToken := r.Header.Get("Authorization")
	tokenParts := strings.SplitAfter(bearerToken, " ")
	token := tokenParts[1]
	if token != "" {
		s.midgardToken = token
		zap.L().Info("Midgard token received from client is ", zap.String("token", s.midgardToken))
	} else {
		http.Error(w, "No token received from client", http.StatusInternalServerError)
	}
}

// UpdateDrinkEndPoint updates an existing drink property given its ID
func (s *Server) UpdateDrinkEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	err := s.checkIfUserAuthenticated(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if s.newConnection {

		m := s.mongodb
		var err error
		if s.newConnection {
			m, err = mongodb.NewMongoSession(s.host, "", "", s.database, s.collection)
			if err != nil {
				zap.L().Error("error creating a session", zap.Error(err))
			}
			defer m.GetSession().Close()

		}

		drinkName := strings.SplitAfter(r.URL.RequestURI(), "/")
		decoder := json.NewDecoder(r.Body)

		if err = m.Update(decoder, drinkName[1]); err != nil {

			zap.L().Error("error updating data in database", zap.Error(err))
		}
	}
}

// DeleteDrinkEndPoint deletes a drink given ID and its type
func (s *Server) DeleteDrinkEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	err := s.checkIfUserAuthenticated(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if s.newConnection {

		m := s.mongodb

		if s.newConnection {
			m, err = mongodb.NewMongoSession(s.host, "", "", s.database, s.collection)
			if err != nil {
				zap.L().Error("error creating a session", zap.Error(err))
			}
			defer m.GetSession().Close()

		}

		params := mux.Vars(r)

		if err = m.Delete(params["id"]); err != nil {

			zap.L().Error("error deleting data from database", zap.Error(err))
		}
	}
}
