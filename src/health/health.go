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

package health

import (
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	. "db/mongo"
	. "rest"
)

type HealthServer struct{
	Server string
	Database string
}

var tns = TNSserver{}
var rest = RESTServer{}
var health = HealthServer{}
var health_server bool = true
var health_first = false

//Keep-alive init function
func InitKeepAlive(){
	if health_server{
		health_server = false	
			r := mux.NewRouter()
			r.HandleFunc("/api/v1/tns/health", rest.TopicKeepAlive).Methods("POST")
			if err := http.ListenAndServe(":48324", r); err != nil {
				log.Fatal(err)
			}
	}
	// TO DO
	//1. Get all TNS topic list from DB
  GetTopicData()

	//2. Edges will ping their service-id 
	println("start Health Check")
	fmt.Println(time.Now().Format(time.RFC850))
}

//Close Keep-alive session
func CloseKeepAlive(){
// TO DO
	//1. sort  unpinged topics(by cid)
	//2. query to DB to DELETE unpinged topics(by cid)
}

//HealthCheck trigger fuction
func (m *HealthServer) TriggerKeepAlive(){
	nextTime := time.Now().Truncate(time.Minute)
  nextTime = nextTime.Add(time.Minute)
	if health_first{
	go InitKeepAlive()
	}
	time.Sleep(time.Until(nextTime))
	health_first = true	
	println("End Health Check")
	fmt.Println(time.Now().Format(time.RFC850))
	go health.TriggerKeepAlive()
}

// Get All Topic List for keep-alive
func GetTopicData(){
// TO DO

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

