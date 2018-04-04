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
	. "rest"
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


// Discover topic
// GET list of tns topics including keyword check
func AllTNSServerList(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	mytopic := queryValues.Get("topic")						 
	fmt.Printf("query : %v\n",queryValues)							 
	fmt.Printf("topic : %s\n",mytopic)							 
	tnsdata, err := tns.DiscoverTopic(mytopic)
//	tnsdata, err := tns.FindAll()
		if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, tnsdata)
}

// Resolution list of tns topics :: PUT => need to change to GET
func ResolutionTopic_PUT(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var tnsdata TNSdata
	if err := json.NewDecoder(r.Body).Decode(&tnsdata); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
	}
	tnsdata.ID = bson.NewObjectId()
	mytopic := tnsdata.Topic	
	tnsdata_res, err := tns.DiscoverTopic(mytopic)	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
			return
	}
	respondWithJson(w, http.StatusOK, tnsdata_res)
}

// Discover topic by keyword
// GET/{topic}
func DiscoverByTopic(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tnsdata_res, err := tns.DiscoverTopic(params["topic"])	
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid topic")
			return
	}
	respondWithJson(w, http.StatusOK, tnsdata_res)
}

// Discover topic by id
// GET/{id}
func FindTNSList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tnsdata, err := tns.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid id")
		return
	}
	// TO DO error check for all cases				
	respondWithJson(w, http.StatusOK, tnsdata)
}

// Register Topic
// POST a new list
func CreateTopicList(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var tnsdata TNSdata
	if err := json.NewDecoder(r.Body).Decode(&tnsdata); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// Validation CHECK for duplicate TOPIC
	if err := tns.CheckDuplicate(tnsdata); err != false {
//		respondWithError(w, http.StatusBadRequest, "Duplicated Topic")
    	respondWithJson(w, http.StatusOK, map[string]string{"result": "duplicated"})
		return
	}
	println("Now will be Updated")
	tnsdata.ID = bson.NewObjectId()
	if err := tns.Insert(tnsdata); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//respondWithJson(w, http.StatusCreated, tnsdata)
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// PUT update an existing lists
func UpdateTopicList(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var tnsdata TNSdata
	if err := json.NewDecoder(r.Body).Decode(&tnsdata); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := tns.Update(tnsdata); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// Unregister Topic
// DELETE an existing lists
func DeleteTNSList(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var tnsdata TNSdata
	if err := json.NewDecoder(r.Body).Decode(&tnsdata); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := tns.Delete(tnsdata); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// POST healthcheck for Topics in TNS server
func TopicHealthcheck(w http.ResponseWriter, r *http.Request) {
// TODO
// GET topic and check for existing TNSDB
// after all check for TNSDB, if there is unchecked topic, than delete it 	
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
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
