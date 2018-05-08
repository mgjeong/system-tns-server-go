/*******************************************************************************
 * Copyright 2017 Samsung Electronics All Rights Reserved.
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

package api

import (
	"fmt"
	"net/http"
	"strings"
	"tns/api/common"
	"tns/api/keepalive"
	"tns/api/topic"
	"tns/commons/errors"
	"tns/commons/logger"
	keepaliveController "tns/controller/keepalive"
	topicDB "tns/db/topic"
)

type RequestHandler struct{}

var Handler RequestHandler

var config = Config{}
var topicHandler topic.Command
var keepAliveHandler keepalive.Command
var keepaliveExecutor keepaliveController.Command
var topicDbExecutor topicDB.Command

func init() {
	topicHandler = topic.RequestHandler{}
	keepAliveHandler = keepalive.RequestHandler{}
	keepaliveExecutor = keepaliveController.Executor{}
	topicDbExecutor = topicDB.Executor{}
}

func RunServer(filePath string) {
	logger.Logging(logger.DEBUG, "RUN TNS Server")

	err := config.Read(filePath)
	if err != nil {
		logger.Logging(logger.ERROR, "Failed to read config file")
		return
	}

	dbUrl := config.Database.Ip + ":" + fmt.Sprint(config.Database.Port)
	err = topicDbExecutor.Connect(dbUrl, config.Database.Name, config.Database.Collection)
	if err != nil {
		logger.Logging(logger.ERROR, "Failed to connect to DB")
		return
	}

	err = keepaliveExecutor.InitKeepAlive(config.Server.KeepAliveInterval)
	if err != nil {
		logger.Logging(logger.ERROR, "Failed to initialize KeepAlive")
		return
	}

	svrUrl := config.Server.Ip + ":" + fmt.Sprint(config.Server.Port)
	http.ListenAndServe(svrUrl, &Handler)
}

func (RequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "IN receive msg", req.Method, req.URL.Path)
	defer logger.Logging(logger.DEBUG, "OUT")

	switch url := req.URL.Path; {
	case !strings.Contains(url, "/api/v1"):
		logger.Logging(logger.DEBUG, "Unknown URL")
		common.WriteError(w, errors.NotFoundURL{})

	case strings.Contains(url, "/tns/topic"):
		topicHandler.Handle(w, req)

	case strings.Contains(url, "/tns/keepalive"):
		keepAliveHandler.Handle(w, req)

	default:
		logger.Logging(logger.DEBUG, "Unknown URL")
		common.WriteError(w, errors.NotFoundURL{url})
	}
}
