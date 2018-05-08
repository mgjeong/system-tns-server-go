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

package topic

import (
	"net/http"
	"strings"
	"tns/api/common"
	"tns/commons/errors"
	"tns/commons/logger"
	topicController "tns/controller/topic"
)

type Command interface {
	Handle(w http.ResponseWriter, req *http.Request)
}

type RequestHandler struct{}

var topicExecutor topicController.Command

func init() {
	topicExecutor = topicController.Executor{}
}

func (RequestHandler) Handle(w http.ResponseWriter, req *http.Request) {
	// Check URL
	url := strings.TrimPrefix(req.URL.Path, "/api/v1"+"/tns/topic")
	if len(url) != 0 {
		common.WriteError(w, errors.NotFoundURL{url})
		return
	}

	switch req.Method {
	case http.MethodPost:
		handlePostReq(w, req)
	case http.MethodGet:
		handleGetReq(w, req)
	case http.MethodDelete:
		handleDeleteReq(w, req)
	default:
		logger.Logging(logger.DEBUG, "Invalid Method")
		common.WriteError(w, errors.InvalidMethod{req.Method})
		return
	}
}

func handlePostReq(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	body, err := common.GetBodyFromReq(req)
	if err != nil {
		logger.Logging(logger.DEBUG, "GetBodyFromReq failed")
		common.WriteError(w, err)
		return
	}

	resp, err := topicExecutor.CreateTopic(body)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteResponse(w, http.StatusCreated, common.MapToJsonByte(resp))
}

func handleGetReq(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Parse query
	name := ""
	hierarchical := false // false is default

	for field, values := range req.URL.Query() {
		if len(values) != 1 {
			common.WriteError(w, errors.InvalidQuery{field}) // No any array type value so far
			return
		}

		switch field {
		case "name":
			name = values[0]
		case "hierarchical":
			if values[0] == "yes" {
				hierarchical = true
			} else if values[0] == "no" {
				hierarchical = false
			} else {
				common.WriteError(w, errors.InvalidQuery{field})
				return
			}
		default:
			logger.Logging(logger.DEBUG, "Invalid query: "+field)
			common.WriteError(w, errors.InvalidQuery{field})
			return
		}
	}

	resp, err := topicExecutor.ReadTopic(name, hierarchical)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteResponse(w, http.StatusOK, common.MapToJsonByte(resp))
}

func handleDeleteReq(w http.ResponseWriter, req *http.Request) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Parse query
	name := ""

	for field, values := range req.URL.Query() {
		if len(values) != 1 { // No any array type value so far
			common.WriteError(w, errors.InvalidQuery{field})
			return
		}

		switch field {
		case "name":
			name = values[0]
		default:
			logger.Logging(logger.DEBUG, "Invalid query: "+field)
			common.WriteError(w, errors.InvalidQuery{field})
			return
		}
	}

	err := topicExecutor.DelteTopic(name)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	common.WriteResponse(w, http.StatusOK, nil)
}
