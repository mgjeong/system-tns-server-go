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
		"fmt"
	//	"io/ioutil"
		"net/http"
	//	"log"
		"database/sql"
		_ "github.com/go-sql-driver/mysql"
)

func TNShandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Topic Name Service Server Start!!!")
}

func main() {
    
		fmt.Println("Go MySQL DB start")
		db, err := sql.Open("mysql", "jinhyuck:siso2010@tcp(127.0.0.1:3306)/test")
    // if there is an error opening the connection, handle it
		 if err != nil {
			 panic(err.Error())
		 }
		  
		// defer the close till after the main function has finished
		// executing 
		defer db.Close()

  	http.HandleFunc("/", TNShandler)
		go http.ListenAndServeTLS(":8081","cert.pem", "key.pem", nil)
		http.ListenAndServe(":8080", nil)
	
}
