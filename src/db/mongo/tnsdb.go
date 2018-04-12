/*******************************************************************************
 * Copyright 2018 Samsung Electronics All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *******************************************************************************/

package tns

import (
	"log"

	"fmt"
	. "tns_model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TNSserver struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "tns"
)

// Establish a connection to database
func (m *TNSserver) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
	println("DB connected")
}

// Find topic list of TNS Server
func (m *TNSserver) FindAll() ([]TNSdata, error) {
	var tns []TNSdata
	err := db.C(COLLECTION).Find(bson.M{}).All(&tns)
	return tns, err
}

// Find a topic list by its id - NOT USED WILL BE DEPRECATED
func (m *TNSserver) FindById(id string) (TNSdata, error) {
	var tns TNSdata
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&tns)
	println("Discovery all Respond")
	return tns, err
}

func (m *TNSserver) CheckDuplicate(tns TNSdata) bool {
	//println("Enter CheckDuplicate")
	mytopic := tns.Topic
	println(mytopic)
	hit, err := db.C(COLLECTION).Find(bson.M{"topic":mytopic}).Count()
	if err != nil{
		panic(err)
	}
	if hit == 0{
		println("There is no duplicated")
			return false
	}
	println("There is hit duplicated")
	return true
}

// discover topic with keywords
func (m *TNSserver) DiscoverTopic(mytopic string) ([]TNSdata, error) {
	var tns []TNSdata
	err := db.C(COLLECTION).Find(bson.M{"topic": bson.RegEx{Pattern: mytopic, Options: "i"}}).All(&tns)
	println("Discover with topic name")
	return tns,err
}

// discover topic for Delete
func (m *TNSserver) DiscoverDELTopic(mytopic string) (TNSdata, error) {
	println("Enter DiscoverDELTopic")
	var tns TNSdata
	err := db.C(COLLECTION).Find(bson.M{"topic": bson.RegEx{Pattern: mytopic, Options: "i"}}).One(&tns)
	println("Discover with topic name")
	return tns,err
}

// Find a topic list by its topic
func (m *TNSserver) FindByTopic(topic string) (TNSdata, error) {
	var tns TNSdata
	//err := db.C(COLLECTION).FindId(bson.ObjectIdHex(topic)).One(&tns)
	err := db.C(COLLECTION).Find(bson.ObjectIdHex(topic)).One(&tns)
	return tns, err
}

// Insert a topic list into database
func (m *TNSserver) Insert(tns TNSdata) error {
	err := db.C(COLLECTION).Insert(&tns)
	println("A New Topic Registered")
	return err
}

// Delete an existing topic list
func (m *TNSserver) Delete(tns TNSdata) error {
	println("Enter Delete of tns")	
	myID := tns.ID
	mytopic := tns.Topic
	myEndpoint := tns.Endpoint
	mySchema := tns.Schema
	fmt.Printf("DELETE id : %s\n",myID)							 
	fmt.Printf("DELETE topic : %s\n",mytopic)							 
	fmt.Printf("DELETE Endpoint : %s\n",myEndpoint)							 
	fmt.Printf("DELETE schema : %s\n",mySchema)							 
	err := db.C(COLLECTION).Remove(&tns)
	println("A Topic Deleted")
	return err
}

// Update an existing topic list
func (m *TNSserver) Update(tns TNSdata) error {
//	err := db.C(COLLECTION).UpdateId(movie.ID, &tns)
	err := db.C(COLLECTION).UpdateId(tns.ID, &tns)
	return err
}

