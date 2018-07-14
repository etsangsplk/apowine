package mongodb

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDB holds database properties
type MongoDB struct {
	session    *mgo.Session
	collection *mgo.Collection
}

// Beer holds beer properties
type Beer struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	BeerName string        `bson:"beername" json:"beername"`
	Type     string        `bson:"type" json:"type"`
}

// Wine holds wine properties
type Wine struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	WineName string        `bson:"winename" json:"winename"`
	Type     string        `bson:"type" json:"type"`
}

// Count holds wine properties
type Count struct {
	ID    bson.ObjectId `bson:"_id" json:"id"`
	Count int           `bson:"count" json:"count"`
	Type  string        `bson:"type" json:"type"`
	Date  time.Time     `bson:"date" json:"date"`
}
