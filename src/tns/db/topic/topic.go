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
	"gopkg.in/mgo.v2/bson"
	"strings"
	"tns/commons/errors"
	"tns/commons/logger"
	mgo "tns/db/wrapper"
)

const (
	WILDCARD         = "*"
	DB_URL           = "127.0.0.1:27017"
	TOPIC_COLLECTION = "TOPIC"
)

type Command interface {
	Connect(name string) error
	Close()
	CreateTopic(map[string]interface{}) error
	ReadTopicAll() ([]map[string]interface{}, error)
	ReadTopic(name string, hierarchical bool) ([]map[string]interface{}, error)
	DeleteTopic(name string) error
}

type Executor struct{}

type Topic struct {
	//ID            bson.ObjectId    `bson:"_id,omitempty"`
	Name      string `bson:"name"`
	Endpoint  string `bson:"endpoint"`
	Datamodel string `bson:"datamodel"`
	Secured   bool   `bson:"secured"`
}

var (
	mgoDial            mgo.Connection
	mgoSession         mgo.Session
	mgoTopicCollection mgo.Collection
)

func init() {
	mgoDial = mgo.MongoDial{}
}

func (topic Topic) convertToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":      topic.Name,
		"endpoint":  topic.Endpoint,
		"datamodel": topic.Datamodel,
		"secured":   topic.Secured,
	}
}

func (m Executor) Connect(name string) error {
	session, err := mgoDial.Dial(DB_URL)
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return err
	}

	mgoSession = session
	mgoTopicCollection = mgoSession.DB(name).C(TOPIC_COLLECTION)

	logger.Logging(logger.DEBUG, "DB connected: "+DB_URL)

	return nil
}

func (m Executor) Close() {
	mgoSession.Close()
}

func (m Executor) CreateTopic(properties map[string]interface{}) error {
	name, exists := properties["name"].(string)
	if !exists {
		return errors.InvalidParam{"'name' field is required"}
	}

	endpoint, exists := properties["endpoint"].(string)
	if !exists {
		return errors.InvalidParam{"'endpoint' field is required"}
	}

	datamodel, exists := properties["datamodel"].(string)
	if !exists {
		return errors.InvalidParam{"'datamodel' field is required"}
	}

	secured, exists := properties["secured"].(bool)
	if !exists {
		secured = false
	}

	exists, err := m.isTopicNameExists(name)
	if err != nil {
		logger.Logging(logger.ERROR, "isTopicNameExists failed")
		return err
	}
	if exists {
		logger.Logging(logger.DEBUG, "Topic already exists: "+name)
		return errors.Conflict{name}
	}

	topic := Topic{
		//ID:            bson.NewObjectId(),
		Name:      name,
		Endpoint:  endpoint,
		Datamodel: datamodel,
		Secured:   secured,
	}

	if err := mgoTopicCollection.Insert(topic); err != nil {
		return errors.InternalServerError{"Database Insert Failed"}
	}

	return nil
}

func (m Executor) DeleteTopic(name string) error {
	pattern := "^" + name + "$"
	query := bson.M{"name": bson.RegEx{Pattern: pattern}}

	err := mgoTopicCollection.Remove(query)
	if err != nil {
		if err == mgo.ErrNotFound {
			logger.Logging(logger.DEBUG, "Not found on mongoDb: "+name)
			return errors.NotFound{name}
		}
		logger.Logging(logger.ERROR, "Failed to Remove on mongoDb: "+name)
		return errors.InternalServerError{"Database Remove Failed"}
	}

	return nil
}

func (m Executor) ReadTopicAll() ([]map[string]interface{}, error) {
	pattern := "."
	topics, err := m.readTopicFromDB(pattern)
	if err != nil {
		logger.Logging(logger.ERROR, "readTopicFromDB failed")
		return nil, err
	}

	return topics, nil
}

func (m Executor) ReadTopic(name string, hierarchical bool) ([]map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	if strings.Contains(name, WILDCARD) {
		if hierarchical {
			return nil, errors.InvalidQuery{"wildcard is not available while hierarchical is yes"}
		}
		return m.readTopicWildcard(name)
	}

	pattern := "^" + name + "$" // One exactly matched
	if hierarchical {
		pattern += "|^" + name + "/" // All started with 'name'
	}

	topics, err := m.readTopicFromDB(pattern)
	if err != nil {
		logger.Logging(logger.ERROR, "readTopicFromDB failed")
		return nil, err
	}

	return topics, nil
}

func (m Executor) readTopicFromDB(pattern string) ([]map[string]interface{}, error) {
	query := bson.M{"name": bson.RegEx{Pattern: pattern}}
	topics := []Topic{}

	err := mgoTopicCollection.Find(query).All(&topics)
	if err != nil {
		logger.Logging(logger.ERROR, "Failed to Find All on mongoDB: "+err.Error())
		return nil, errors.InternalServerError{"Database Query Failed"}
	}

	topicsInterface := make([]map[string]interface{}, len(topics))
	for i, topic := range topics {
		topicsInterface[i] = topic.convertToMap()
	}

	return topicsInterface, nil
}

func (m Executor) readTopicWildcard(name string) ([]map[string]interface{}, error) {
	logger.Logging(logger.DEBUG, "IN")
	defer logger.Logging(logger.DEBUG, "OUT")

	// @TODO not supported yet
	return nil, errors.InvalidQuery{"wildcard is not supported yet"}
}

func (m Executor) isTopicNameExists(name string) (bool, error) {
	pattern := "^" + name + "$"
	query := bson.M{"name": bson.RegEx{Pattern: pattern}}

	hit, err := mgoTopicCollection.Find(query).Count()
	if err != nil {
		logger.Logging(logger.ERROR, err.Error())
		return true, errors.InternalServerError{"Database Query Failed"}
	}
	return (hit != 0), nil
}
