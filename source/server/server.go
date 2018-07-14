package server

import (
	"encoding/json"
	"net/http"
	"strings"

	mongodb "github.com/aporeto-inc/apowine/source/mongodb-lib"
	"github.com/aporeto-inc/apowine/source/server/configuration"
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
	session       *sessions.Session
	beerReqCount  int
	wineReqCount  int
	mongoData     chan mongoData
	stop          chan struct{}
}

type mongoData struct {
	data      *json.Decoder
	drinkName string
	count     int
	m         *mongodb.MongoDB
}

// NewServer creates a new server handler
func NewServer(mongo *mongodb.MongoDB, host []string, cfg *configuration.Configuration) *Server {
	zap.L().Info("Creating a new server handler")

	if cfg.MakeNewConnection && mongo != nil {
		mongo.GetSession().Close()
	}

	return &Server{
		mongodb:       mongo,
		newConnection: cfg.MakeNewConnection,
		host:          host,
		database:      cfg.MongoDatabaseName,
		collection:    cfg.MongoCollectionName,
		mongoData:     make(chan mongoData, 500),
		stop:          make(chan struct{}),
	}
}

// Start pushing count to db
func (s *Server) Start() error {

	go func() {
		for {
			select {
			case mongodata := <-s.mongoData:
				if err := mongodata.m.InsertOrUpdateCount(mongodata.data, mongodata.drinkName, mongodata.count); err != nil {
					zap.L().Error("error creating count record", zap.Error(err))
				}
				if s.newConnection {
					mongodata.m.GetSession().Close()
				}
			case <-s.stop:
				return
			}
		}
	}()

	return nil
}

// Stop the goroutine
func (s *Server) Stop() error {

	s.stop <- struct{}{}

	return nil
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
	}

	endpoint := strings.SplitAfter(r.URL.RequestURI(), "/")
	drinkName := strings.Replace(endpoint[1], "/", "", -1)
	decoder := json.NewDecoder(r.Body)
	data, err := m.Read(decoder, drinkName, true)
	if err != nil {
		zap.L().Error("error reading data from database", zap.Error(err))
	}

	if err = json.NewEncoder(w).Encode(data); err != nil {
		zap.L().Error("error in json output", zap.Error(err))
	}

	var count int
	switch drinkName {
	case mongodb.BEER:
		s.beerReqCount++
		count = s.beerReqCount
	case mongodb.WINE:
		s.wineReqCount++
		count = s.wineReqCount
	}

	s.mongoData <- mongoData{
		data:      decoder,
		drinkName: drinkName,
		count:     count,
		m:         m,
	}
}

// GetBeerCount returns number of beers served
func (s *Server) GetBeerCount(w http.ResponseWriter, r *http.Request) {
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
	data, err := m.ReadCount(decoder, drinkName)
	if err != nil {
		zap.L().Error("error reading data from database", zap.Error(err))
	}

	if err = json.NewEncoder(w).Encode(data); err != nil {
		zap.L().Error("error in json output", zap.Error(err))
	}
}

// GetWineCount returns number of wines served
func (s *Server) GetWineCount(w http.ResponseWriter, r *http.Request) {
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
	data, err := m.ReadCount(decoder, drinkName)
	if err != nil {
		zap.L().Error("error reading data from database", zap.Error(err))
	}

	if err = json.NewEncoder(w).Encode(data); err != nil {
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
	var err error

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

	if err = json.NewEncoder(w).Encode(data); err != nil {
		zap.L().Error("error in json output", zap.Error(err))
	}
}

// CreateDrinkEndPoint creates a drink (beer or wine) with an ID to use in MongoDB
func (s *Server) CreateDrinkEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var err error
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

// UpdateDrinkEndPoint updates an existing drink property given its ID
func (s *Server) UpdateDrinkEndPoint(w http.ResponseWriter, r *http.Request) {
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

	drinkName := strings.SplitAfter(r.URL.RequestURI(), "/")
	decoder := json.NewDecoder(r.Body)

	if err = m.Update(decoder, drinkName[1]); err != nil {
		zap.L().Error("error updating data in database", zap.Error(err))
	}

}

// DeleteDrinkEndPoint deletes a drink given ID and its type
func (s *Server) DeleteDrinkEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var err error

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
