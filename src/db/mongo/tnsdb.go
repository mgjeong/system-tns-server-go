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
}

// Find topic list of TNS Server
func (m *TNSserver) FindAll() ([]TNSdata, error) {
	var tns []TNSdata
	err := db.C(COLLECTION).Find(bson.M{}).All(&tns)
	return tns, err
}

// Find a topic list by its id
func (m *TNSserver) FindById(id string) (TNSdata, error) {
	var tns TNSdata
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&tns)
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

func (m *TNSserver) ResolutionTNS(mytopic string) ([]TNSdata, error) {
	println("Entered ResoutionTNS")
	var tns []TNSdata
	err := db.C(COLLECTION).Find(bson.M{"topic": bson.RegEx{Pattern: mytopic, Options: "i"}}).All(&tns)
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
	return err
}

// Delete an existing topic list
func (m *TNSserver) Delete(tns TNSdata) error {
	err := db.C(COLLECTION).Remove(&tns)
	return err
}

// Update an existing topic list
func (m *TNSserver) Update(tns TNSdata) error {
//	err := db.C(COLLECTION).UpdateId(movie.ID, &tns)
	err := db.C(COLLECTION).UpdateId(tns.ID, &tns)
	return err
}

