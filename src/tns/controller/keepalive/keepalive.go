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
	"sync"
	"time"
	"tns/commons/errors"
	"tns/commons/logger"
	"tns/commons/util"
	topicDB "tns/db/topic"
)

type Command interface {
	InitKeepAlive(interval uint) error
	AddTopic(name string)
	DeleteTopic(name string)
	HandlePing(body string) (map[string]interface{}, error)
	GetInterval() uint
}

// Executor implements the Command interface.
type Executor struct{}

type kaTableType map[string]time.Time // "topic":"timestamp"

type keepAliveInfo struct {
	sync.Mutex
	table    kaTableType
	interval uint
}

const kaPingFrequency = 3

var topicDbExecutor topicDB.Command
var kaInfo keepAliveInfo

func init() {
	topicDbExecutor = topicDB.Executor{}
}

func (Executor) InitKeepAlive(interval uint) error {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// Read Topics from DB
	topics, err := topicDbExecutor.ReadTopicAll()
	if err != nil {
		logger.Logging(logger.ERROR, "ReadTopicAll failed")
		return err
	}

	// Init Keepalive Table
	logger.Logging(logger.DEBUG, "Initialize Keep-alive Table")
	kaInfo.table = make(kaTableType)
	currTime := time.Now()
	for _, topic := range topics {
		logger.Logging(logger.DEBUG, topic["name"].(string))
		kaInfo.table[topic["name"].(string)] = currTime
	}

	kaInfo.interval = interval

	// Start Timer loop
	go keepAliveTimerLoop(interval)

	return nil
}

func (Executor) AddTopic(name string) {
	currTime := time.Now()

	kaInfo.Lock()
	kaInfo.table[name] = currTime
	kaInfo.Unlock()

	logger.Logging(logger.DEBUG, "Topic added: "+name)
}

func (Executor) DeleteTopic(name string) {
	kaInfo.Lock()
	delete(kaInfo.table, name)
	kaInfo.Unlock()

	logger.Logging(logger.DEBUG, "Topic deleted: "+name)
}

func (Executor) HandlePing(body string) (map[string]interface{}, error) {
	bodyMap, err := util.ConvertJsonToMap(body)
	if err != nil {
		logger.Logging(logger.ERROR, "ConvertJsonToMap failed: "+err.Error())
		return nil, err
	}

	topicNamesInterface, exists := bodyMap["topic_names"].([]interface{})
	if !exists {
		logger.Logging(logger.DEBUG, "'topic_names' does not present in body")
		return nil, errors.InvalidJSON{"'topic_names' field is required"}
	}

	topicNames := make([]string, len(topicNamesInterface))
	for i, v := range topicNamesInterface {
		name, exists := v.(string)
		if !exists {
			return nil, errors.InvalidParam{"topic_names"}
		}
		topicNames[i] = name
	}

	var notFound []string
	currTime := time.Now()

	kaInfo.Lock()
	for _, name := range topicNames {
		_, exists := kaInfo.table[name]
		if exists {
			// Update timestamp
			kaInfo.table[name] = currTime
		} else {
			notFound = append(notFound, name)
		}
	}
	kaInfo.Unlock()

	if len(notFound) != 0 {
		resp := make(map[string]interface{})
		resp["topic_names"] = notFound

		return resp, errors.NotFound{}
	}

	return nil, nil
}

func (Executor) GetInterval() uint {
	return kaInfo.interval / kaPingFrequency
}

func keepAliveTimerLoop(interval uint) {
	logger.Logging(logger.DEBUG, "Start KeepAlive Timer loop")
	defer logger.Logging(logger.DEBUG, "KeepAlive Timer loop Finished")

	timeDurationSec := time.Duration(interval) * time.Second
	ticker := time.NewTicker(timeDurationSec)

	for _ = range ticker.C {
		kaInfo.Lock()
		for topic, timestamp := range kaInfo.table {
			// Remove expired topics
			if time.Since(timestamp) > timeDurationSec {
				logger.Logging(logger.DEBUG, "KeepAlive time expired: "+topic)
				// Delete topic from DB
				if err := topicDbExecutor.DeleteTopic(topic); err != nil {
					logger.Logging(logger.ERROR, "DeleteTopic failed")
				}
				// Delete from KA table
				delete(kaInfo.table, topic)
				logger.Logging(logger.DEBUG, "Topic deleted: "+topic)
			}
		}
		kaInfo.Unlock()
	}
}
