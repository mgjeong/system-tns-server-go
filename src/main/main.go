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
	"log"
	"net/http"
	"github.com/gorilla/mux"
	. "config"
	. "db/mongo"
	. "rest"
	. "health"
)

var config = Config{}
var tns = TNSserver{}
var rest = RESTServer{}
var health = HealthServer{}


// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	println("Entered init")
	config.Read()
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
  go health.TriggerHealthCheck()
	if err := http.ListenAndServe(":" + config.Port, r); err != nil {
		log.Fatal(err)
	}
}
