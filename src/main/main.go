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
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
	. "config"
	. "tns_model"
	. "db/mongo"
)

var config = Config{}
var tns = TNSserver{}

// GET list of tns topics
func AllTNSServerList(w http.ResponseWriter, r *http.Request) {
	  tnsdata, err := tns.FindAll()
		if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, tnsdata)
}

// GET a list by its topic
func FindTNSList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tnsdata, err := tns.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid topic")
		return
	}
	// TO DO error check for all cases				
	respondWithJson(w, http.StatusOK, tnsdata)
}

// POST a new list
func CreateTopicList(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var tnsdata TNSdata
	if err := json.NewDecoder(r.Body).Decode(&tnsdata); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	tnsdata.ID = bson.NewObjectId()
	if err := tns.Insert(tnsdata); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, tnsdata)
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
	r.HandleFunc("/tnsdb", AllTNSServerList).Methods("GET")
	r.HandleFunc("/tnsdb", CreateTopicList).Methods("POST")
	r.HandleFunc("/tnsdb", UpdateTopicList).Methods("PUT")
	r.HandleFunc("/tnsdb", DeleteTNSList).Methods("DELETE")
	r.HandleFunc("/tnsdb/{topic}", FindTNSList).Methods("GET")
	if err := http.ListenAndServe(":" + config.Port, r); err != nil {
		log.Fatal(err)
	}
}
