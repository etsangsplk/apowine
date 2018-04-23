package mongodb

import (
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
