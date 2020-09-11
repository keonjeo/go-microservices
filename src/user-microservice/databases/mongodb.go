/*
 * @File: databases.mongodb.go
 * @Description: Handles MongoDB connections
 * @Author: Nguyen Truong Duong (seedotech@gmail.com)
 */
package databases

import (
	"time"

	log "github.com/sirupsen/logrus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"user_micro_service/common"
	"user_micro_service/models"
)

// MongoDB manages MongoDB connection
type MongoDB struct {
	MgDbSession  *mgo.Session
	DatabaseName string
}

// Init initializes mongo database
func (db *MongoDB) Init() error {
	db.DatabaseName = common.Config.MgDbName

	// DialInfo holds options for establishing a session with a MongoDB cluster.
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{common.Config.MgAddrs}, // Get HOST + PORT
		Timeout:  60 * time.Second,
		Database: db.DatabaseName,            // Database name
		Username: common.Config.MgDbUsername, // Username
		Password: common.Config.MgDbPassword, // Password
	}

	// Create a session which maintains a pool of socket connections
	// to the DB MongoDB database.
	var err error
	db.MgDbSession, err = mgo.DialWithInfo(dialInfo)

	if err != nil {
		log.Debug("Can't connect to mongo, go error: ", err)
		return err
	}

	return db.initData()
}

// InitData initializes default data
func (db *MongoDB) initData() error {
	var err error
	var count int

	// Check if user collection has at least one document
	sessionCopy := db.MgDbSession.Copy()
	defer sessionCopy.Close()

	// Get a collection to execute the query against.
	collection := sessionCopy.DB(db.DatabaseName).C(common.ColUsers)
	count, err = collection.Find(bson.M{}).Count()

	if count < 1 {
		// Create admin/admin account
		var user models.User
		user = models.User{ID: bson.NewObjectId(), Name: "admin", Password: "admin"}
		err = collection.Insert(&user)
	}

	return err
}

// Close the existing connection
func (db *MongoDB) Close() {
	if db.MgDbSession != nil {
		db.MgDbSession.Close()
	}
}
