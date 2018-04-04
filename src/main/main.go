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

package main

import (
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
	. "config"
	. "tns_model"
	. "db/mongo"
)

var config = Config{}
var tns = TNSserver{}
var rest = RESTServer{}

//HealthCheck function
func HealthCheck(){
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/tns/health", rest.TopicHealthcheck).Methods("POST")
	// TO DO
	//1. Get all TNS topic list from DB
	//2. Edges will ping their service-id 
	println("start Health Check")
	fmt.Println(time.Now().Format(time.RFC850))
}


//HealthCheck trigger fuction
func TriggerHealthCheck(){
	nextTime := time.Now().Truncate(time.Minute)
  nextTime = nextTime.Add(10*time.Minute)
 // nextTime = nextTime.Add(time.Minute)
	time.Sleep(time.Until(nextTime))
	HealthCheck()
	go TriggerHealthCheck()
}


// POST healthcheck for Topics in TNS server
func TopicHealthcheck(w http.ResponseWriter, r *http.Request) {
// TODO
// GET topic and check for existing TNSDB
// after all check for TNSDB, if there is unchecked topic, than delete it 	
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	tns.Connect()
	tns.Database = config.Database
	tns.Connect()
}

// Define HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/tns/topic", rest.AllTNSServerList).Methods("GET")
	r.HandleFunc("/api/v1/tns/topic", rest.CreateTopicList).Methods("POST")
	r.HandleFunc("/api/v1/tns/topic", rest.DeleteTNSList).Methods("DELETE")
	//r.HandleFunc("/api/v1/tns/topic/{topic}", DiscoverByTopic).Methods("GET")
//	r.HandleFunc("/api/v1/tns/health", TopicHealthcheck).Methods("POST")
	TriggerHealthCheck()
	if err := http.ListenAndServe(":" + config.Port, r); err != nil {
		log.Fatal(err)
	}
}
