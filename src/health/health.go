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
  "strings"
	"time"
	"io/ioutil"
	"github.com/gorilla/mux"
	. "keepalive_model"
	. "postka_model"
	. "tns_model"
	. "db/mongo"
)

type HealthServer struct{
	Server string
	Database string
}

var tns = TNSserver{}
var health = HealthServer{}
var health_server bool = true
var health_first bool = true
var keepAliveList [1000]KAList
var list_size int

//Keep-alive init function
func InitKeepAlive(){
	if health_server{
		health_server = false	
	    println("Init Keep Alive REST POST Server")
			r := mux.NewRouter()
			r.HandleFunc("/api/v1/tns/keepalive", health.TopicKeepAlive).Methods("POST")
			if err := http.ListenAndServe(":48324", r); err != nil {
				log.Fatal(err)
			}
	}
	//1. Get all TNS topic list from DB
  for i:= 0;i<1000;i++{
	keepAliveList[i].Topic = ""
	keepAliveList[i].Status = false
}
	println("GetTopicData Check start")
  GetTopicData()

	//2. Edges will ping their service-id 
	println("start Keep Alive")
	fmt.Println(time.Now().Format(time.RFC850))
}


// POST healthcheck for Topics in TNS server
func (m *HealthServer) TopicKeepAlive(w http.ResponseWriter, r *http.Request) {
	println("Entered TopicKeepAlive")
	defer r.Body.Close()
	var kadata []POST_ka
  body, err1 := ioutil.ReadAll(r.Body);
	if err1 != nil {
	}
	fmt.Printf("Body string : %s\n",body)							 
	json.Unmarshal(body, &kadata)
  var i = 0
	var size = len(kadata)
	fmt.Printf("kadata size : %d\n", size)
	for i = 0; i < len(kadata); i++ {
	    fmt.Printf("idx : %d\n", i)
	    fmt.Printf("topic : %s\n", kadata[i].Topic)
	}
  var idx = 0
	for idx = 0; idx < len(kadata); idx++ {
	health.CheckKeepAlive(kadata[idx].Topic)
	}
	println("Keep Alive POST test done")
	respondWithJson(w, http.StatusOK, map[string]string{"result": "keep alive test success"})
}

//Check KeepAlive from List
func (m *HealthServer) CheckKeepAlive(topic string) {
	var idx = 0
	for idx = 0; idx < list_size; idx++ {
  if strings.EqualFold(topic,keepAliveList[idx].Topic) == true {
	  fmt.Printf("keepalive for loop idx : %d\n", idx)
		keepAliveList[idx].Status = true
	  fmt.Printf("keepalive selected topic : %d  %d \n", idx, keepAliveList[idx].Status)
		break	
	}
	}
}

//Close Keep-alive session
func (m *HealthServer) CloseKeepAlive() {
	var idx = 0
	for idx = 0; idx < list_size; idx++ {
  if keepAliveList[idx].Status == false {
	  fmt.Printf("keepalive erase topic idx : %d\n", idx)
		mytopic := keepAliveList[idx].Topic	
		fmt.Printf("REST DELETE topic : %s\n",mytopic)							 
		tnsdata, err := tns.DiscoverDELTopic(mytopic)
		if err != nil {
				return
		}
		if err1 := tns.Delete(tnsdata); err1 != nil {
		}
	time.Sleep(10*time.Millisecond)
	}
	}
}

//HealthCheck trigger fuction
func (m *HealthServer) TriggerKeepAlive() {
	nextTime := time.Now().Truncate(time.Minute)
  nextTime = nextTime.Add(10*time.Minute)
//  nextTime = nextTime.Add(3*time.Minute)
	go InitKeepAlive()
	time.Sleep(time.Until(nextTime))
	health_first = false	
  go health.CloseKeepAlive()
	time.Sleep(time.Minute)
	println("End Health Check")
	fmt.Println(time.Now().Format(time.RFC850))
	go health.TriggerKeepAlive()
}

// Get All Topic List for keep-alive
func GetTopicData(){
	println("Entered GetTopicData()")
	var tnsdata []TNSdata
	tnsdata, err := tns.FindAll()
	println("tns.FindAll() finished")
	list_size = len(tnsdata)
	fmt.Printf("list_size => : %d\n", list_size)
	if err != nil{
		return
	}
	var idx = 0
	for idx = 0; idx < len(tnsdata); idx++ {
	    fmt.Printf("idx : %d\n", idx)
	  	keepAliveList[idx].Topic= tnsdata[idx].Topic	 
	    fmt.Printf("topic : %s\n", keepAliveList[idx].Topic)
	}
	for i := 0;i<idx;i++{
	fmt.Printf("topic : %s\n", keepAliveList[i].Topic)
	fmt.Printf("status : %s\n", keepAliveList[i].Status)
	}
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

