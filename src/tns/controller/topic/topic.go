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
	"tns/commons/errors"
	"tns/commons/logger"
	"tns/commons/util"
	keepaliveController "tns/controller/keepalive"
	topicDB "tns/db/topic"
)

type Command interface {
	CreateTopic(body string) (map[string]interface{}, error)
	ReadTopic(name string, hierarchical bool) (map[string]interface{}, error)
	DelteTopic(name string) error
}

// Executor implements the Command interface.
type Executor struct{}

var topicDbExecutor topicDB.Command
var keepaliveExecutor keepaliveController.Command

func init() {
	topicDbExecutor = topicDB.Executor{}
	keepaliveExecutor = keepaliveController.Executor{}
}

func (Executor) CreateTopic(body string) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	bodyMap, err := util.ConvertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, "ConvertJsonToMap failed: "+err.Error())
		return nil, err
	}

	topic, exists := bodyMap["topic"].(map[string]interface{})
	if !exists {
		logger.Logging(logger.DEBUG, "'topic' does not present in body")
		return nil, errors.InvalidJSON{"'topic' field is required"}
	}

	name, exists := topic["name"].(string)
	if !exists {
		logger.Logging(logger.DEBUG, "'name' does not present in body")
		return nil, errors.InvalidJSON{"'name' field is required"}
	}

	err = topicDbExecutor.CreateTopic(topic)
	if err != nil {
		logger.Logging(logger.DEBUG, "CreateTopic failed: "+err.Error())
		return nil, err
	}

	keepaliveExecutor.AddTopic(name)

	resp := make(map[string]interface{})
	resp["ka_interval"] = keepaliveExecutor.GetInterval()

	return resp, nil
}

func (Executor) ReadTopic(name string, hierarchical bool) (map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	var topics []map[string]interface{}
	var err error

	if name == "" {
		topics, err = topicDbExecutor.ReadTopicAll()
	} else {
		topics, err = topicDbExecutor.ReadTopic(name, hierarchical)
	}

	if err != nil {
		return nil, err
	} else if len(topics) == 0 {
		logger.Logging(logger.DEBUG, "Nothing found")
		if name == "" {
			name = "topic is empty"
		}
		return nil, errors.NotFound{name}
	}

	resp := make(map[string]interface{})
	resp["topics"] = topics

	return resp, nil
}

func (Executor) DelteTopic(name string) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	err := topicDbExecutor.DeleteTopic(name)
	if err != nil {
		logger.Logging(logger.DEBUG, "DeleteTopic failed")
		return err
	}

	keepaliveExecutor.DeleteTopic(name)

	return nil
}
