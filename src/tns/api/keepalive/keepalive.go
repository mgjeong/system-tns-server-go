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
package keepalive

import (
	"net/http"
	"strings"
	"tns/api/common"
	"tns/commons/errors"
	"tns/commons/logger"
	keepaliveController "tns/controller/keepalive"
)

type Command interface {
	Handle(w http.ResponseWriter, req *http.Request)
}

type RequestHandler struct{}

var keepaliveExecutor keepaliveController.Command

func init() {
	keepaliveExecutor = keepaliveController.Executor{}
}

func (RequestHandler) Handle(w http.ResponseWriter, req *http.Request) {
	// Check URL
	url := strings.TrimPrefix(req.URL.Path, "/api/v1"+"/tns/keepalive")
	if len(url) != 0 {
		logger.Logging(logger.DEBUG, "Invalid URL")
		common.WriteError(w, errors.NotFoundURL{url})
		return
	}

	switch req.Method {
	case http.MethodPost:
		handlePostReq(w, req)
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

	resp, err := keepaliveExecutor.HandlePing(body)
	if err != nil {
		switch err.(type) {
		case errors.NotFound:
			common.WriteResponse(w, http.StatusNotFound, common.MapToJsonByte(resp))
		default:
			common.WriteError(w, err)
		}
		return
	}

	common.WriteResponse(w, http.StatusOK, common.MapToJsonByte(resp))
}
