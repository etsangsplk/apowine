package mongodb

import (
	"encoding/json"
	"math/rand"
	"time"

	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// BEER ...
	BEER = "beer"
	// WINE ...
	WINE = "wine"
	// RANDOM ...
	RANDOM = "random"
)

// NewMongoSession creates a new database session with options
func NewMongoSession(host []string, username string, password string, database string, collection string) (*MongoDB, error) {

	session, err := connectToMongoDBSession(host, username, password, database)
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		session:    session,
		collection: session.DB(database).C(collection),
	}, nil
}

func connectToMongoDBSession(host []string, username string, password string, database string) (*mgo.Session, error) {

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    host,
		Username: username,
		Password: password,
		Database: database,
	})

	return session, err
}

// GetSession returns session object
func (m *MongoDB) GetSession() *mgo.Session {

	return m.session
}

// GetCollection returns collection object
func (m *MongoDB) GetCollection() *mgo.Collection {

	return m.collection
}

// Insert decodes and adds a new data into the database given a drinkType
func (m *MongoDB) Insert(data *json.Decoder, drinkType string) error {
	var beer Beer
	var wine Wine

	zap.L().Info("Inserting data into database")
	zap.L().Info("Drinktype", zap.String("type", drinkType))

	switch drinkType {
	case BEER:
		data.Decode(&beer)
		beer.ID = bson.NewObjectId()
		beer.Type = BEER
		err := m.collection.Insert(&beer)
		if err != nil {
			return err
		}
		zap.L().Info("Inserted data", zap.Any("data", beer))
	case WINE:
		data.Decode(&wine)
		wine.ID = bson.NewObjectId()
		wine.Type = WINE
		err := m.collection.Insert(&wine)
		if err != nil {
			return err
		}
		zap.L().Info("Inserted data", zap.Any("data", wine))
	}

	return nil
}

// Read decodes and lists data available in the database
func (m *MongoDB) Read(data *json.Decoder, drinkType string, isRandom bool) (interface{}, error) {
	var beers []Beer
	var wines []Wine

	data.Decode(&beers)
	data.Decode(&wines)

	zap.L().Info("Reading data from database")
	zap.L().Info("Drinktype", zap.String("type", drinkType))

	switch drinkType {
	case BEER:
		err := m.collection.Find(bson.M{}).All(&beers)
		if err != nil {
			return nil, err
		}
		zap.L().Info("data", zap.Any("data", beers))
	case WINE:
		err := m.collection.Find(bson.M{}).All(&wines)
		if err != nil {
			return nil, err
		}
		zap.L().Info("data", zap.Any("data", wines))
	case RANDOM:
		m.collection.Find(bson.M{}).All(&beers)
		m.collection.Find(bson.M{}).All(&wines)
	}

	return readRandom(beers, wines, isRandom, drinkType), nil
}

func readRandom(beers []Beer, wines []Wine, random bool, drinkType string) interface{} {
	rand.Seed(time.Now().Unix())
	if len(beers) > 0 && !random && drinkType != RANDOM {
		return beers
	} else if len(wines) > 0 && !random && drinkType != RANDOM {
		return wines
	} else if len(beers) > 0 && random && drinkType != RANDOM {
		return beers[rand.Intn(len(beers))]
	} else if len(wines) > 0 && random && drinkType != RANDOM {
		return wines[rand.Intn(len(wines))]
	} else if random {
		var drinks []interface{}
		drinks = append(drinks, beers[rand.Intn(len(beers))])
		drinks = append(drinks, wines[rand.Intn(len(wines))])
		return drinks[rand.Intn(len(drinks))]
	}

	return nil
}

// ReadByID decodes and returns data given ID and type
func (m *MongoDB) ReadByID(id string, data *json.Decoder, drinkType string) (interface{}, error) {
	zap.L().Info("Reading data from database by given ID")
	zap.L().Info("ID", zap.String("ID", id))

	var beer Beer
	var wine Wine

	switch drinkType {
	case BEER:
		data.Decode(&beer)
		err := m.collection.FindId(bson.ObjectIdHex(id)).One(&beer)
		if err != nil {
			return nil, err
		}
		zap.L().Info("data", zap.Any("data", beer))
		return beer, nil
	case WINE:
		data.Decode(&wine)
		err := m.collection.FindId(bson.ObjectIdHex(id)).One(&wine)
		if err != nil {
			return nil, err
		}
		zap.L().Info("data", zap.Any("data", wine))
		return wine, nil
	}

	return nil, nil
}

// Update decodes and modifies data in database given type
func (m *MongoDB) Update(data *json.Decoder, drinkType string) error {
	zap.L().Info("Updating data in database")
	zap.L().Info("Drinktype", zap.String("type", drinkType))

	var beer Beer
	var wine Wine

	switch drinkType {
	case BEER:
		data.Decode(&beer)
		err := m.collection.UpdateId(beer.ID, &beer)
		if err != nil {
			return err
		}
		zap.L().Info("Updating data in database", zap.Any("data", beer))
	case WINE:
		data.Decode(&wine)
		err := m.collection.UpdateId(wine.ID, &wine)
		if err != nil {
			return err
		}
		zap.L().Info("Updating data in database", zap.Any("data", wine))
	}

	return nil
}

// Delete decodes and removes a record given ID
func (m *MongoDB) Delete(id string) error {
	zap.L().Info("ID", zap.String("ID", id))

	err := m.collection.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		return err
	}

	return nil
}
