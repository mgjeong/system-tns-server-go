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
	"reflect"
	"regexp"
	"testing"
	"tns/commons/errors"
	mgo "tns/db/wrapper"
	mgoMock "tns/db/wrapper/mocks"

	"github.com/golang/mock/gomock"
	"gopkg.in/mgo.v2/bson"
)

var Handler Command

func init() {
	Handler = Executor{}
}

func TestCallConnect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mgoConnectionMockObj := mgoMock.NewMockConnection(ctrl)
	mgoSessionMockObj := mgoMock.NewMockSession(ctrl)
	mgoDatabaseMockObj := mgoMock.NewMockDatabase(ctrl)

	// pass mockObj to a real object.
	mgoDial = mgoConnectionMockObj

	name := "topic"

	testCases := []struct {
		name           string
		mockRetSession mgoMock.MockSession
		mockRetError   error
		expectedError  error
	}{
		{"Success", *mgoSessionMockObj, nil, nil},
		{"DialFailed", *mgoSessionMockObj, errors.Unknown{}, errors.Unknown{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			callFist := mgoConnectionMockObj.EXPECT().Dial(DB_URL).Return(mgoSessionMockObj, tc.mockRetError)

			if tc.name == "Success" {
				callSecond := mgoSessionMockObj.EXPECT().DB(name).Return(mgoDatabaseMockObj).After(callFist)
				mgoDatabaseMockObj.EXPECT().C(TOPIC_COLLECTION).After(callSecond)
			}

			err := Handler.Connect(name)
			if reflect.TypeOf(err) != reflect.TypeOf(tc.expectedError) {
				t.Errorf("Expected Error: %s, Actual: %s", tc.expectedError, err)
			}
		})
	}
}

func TestCallClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mgoSessionMockObj := mgoMock.NewMockSession(ctrl)

	// pass mockObj to a real object.
	mgoSession = mgoSessionMockObj

	gomock.InOrder(
		mgoSessionMockObj.EXPECT().Close(),
	)

	Handler.Close()
}

func TestCallCreateTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mgoCollectionMockObj := mgoMock.NewMockCollection(ctrl)
	mgoQueryMockObj := mgoMock.NewMockQuery(ctrl)

	// pass mockObj to a real object.
	mgoTopicCollection = mgoCollectionMockObj

	dummyProperties := map[string]interface{}{"name": "/a", "endpoint": "0.0.0.0:1234", "datamodel": "test_0.0.1"}
	dummyQuery := bson.M{"name": "/a"}
	dummpyTopic := Topic{
		Name:      "/a",
		Endpoint:  "0.0.0.0:1234",
		Datamodel: "test_0.0.1",
	}

	gomock.InOrder(
		mgoCollectionMockObj.EXPECT().Find(dummyQuery).Return(mgoQueryMockObj),
		mgoQueryMockObj.EXPECT().Count().Return(0, nil),
		mgoCollectionMockObj.EXPECT().Insert(dummpyTopic).Return(nil),
	)

	err := Handler.CreateTopic(dummyProperties)
	if err != nil {
		t.Errorf("CreateTopic returned an error: %s", err.Error())
	}
}

func TestCallCreateTopicWithInvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mgoCollectionMockObj := mgoMock.NewMockCollection(ctrl)
	mgoQueryMockObj := mgoMock.NewMockQuery(ctrl)

	// pass mockObj to a real object.
	mgoTopicCollection = mgoCollectionMockObj

	testCases := []struct {
		name            string
		dummyProperties map[string]interface{}
		dummyQuery      bson.M
		dummpyTopic     Topic
		expectedError   error
	}{
		{"InvalidParam_name", map[string]interface{}{"endpoint": "0.0.0.0:1234", "datamodel": "test_0.0.1"}, nil, Topic{}, errors.InvalidParam{}},
		{"InvalidParam_endpoint", map[string]interface{}{"name": "/a", "datamodel": "test_0.0.1"}, nil, Topic{}, errors.InvalidParam{}},
		{"InvalidParam_datamodel", map[string]interface{}{"name": "/a", "endpoint": "0.0.0.0:1234"}, nil, Topic{}, errors.InvalidParam{}},
		{"DbFailed_Find", map[string]interface{}{"name": "/a", "endpoint": "0.0.0.0:1234", "datamodel": "test_0.0.1"}, bson.M{"name": "/a"}, Topic{}, errors.InternalServerError{}},
		{"TopicAlreadyExists", map[string]interface{}{"name": "/a", "endpoint": "0.0.0.0:1234", "datamodel": "test_0.0.1"}, bson.M{"name": "/a"}, Topic{}, errors.Conflict{}},
		{"DbFailed_Insert", map[string]interface{}{"name": "/a", "endpoint": "0.0.0.0:1234", "datamodel": "test_0.0.1"}, bson.M{"name": "/a"}, Topic{Name: "/a", Endpoint: "0.0.0.0:1234", Datamodel: "test_0.0.1"}, errors.InternalServerError{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			if tc.name == "DbFailed_Find" {
				gomock.InOrder(
					mgoCollectionMockObj.EXPECT().Find(tc.dummyQuery).Return(mgoQueryMockObj),
					mgoQueryMockObj.EXPECT().Count().Return(0, errors.Unknown{}),
				)
			} else if tc.name == "TopicAlreadyExists" {
				gomock.InOrder(
					mgoCollectionMockObj.EXPECT().Find(tc.dummyQuery).Return(mgoQueryMockObj),
					mgoQueryMockObj.EXPECT().Count().Return(1, nil),
				)
			} else if tc.name == "DbFailed_Insert" {
				gomock.InOrder(
					mgoCollectionMockObj.EXPECT().Find(tc.dummyQuery).Return(mgoQueryMockObj),
					mgoQueryMockObj.EXPECT().Count().Return(0, nil),
					mgoCollectionMockObj.EXPECT().Insert(tc.dummpyTopic).Return(errors.Unknown{}),
				)
			}

			err := Handler.CreateTopic(tc.dummyProperties)
			if reflect.TypeOf(err) != reflect.TypeOf(tc.expectedError) {
				t.Errorf("Expected Error: %s, Actual: %s", tc.expectedError, err)
			}
		})
	}
}

func TestCallDeleteTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mgoCollectionMockObj := mgoMock.NewMockCollection(ctrl)

	// pass mockObj to a real object.
	mgoTopicCollection = mgoCollectionMockObj

	dummyName := "/a"
	dummyQuery := bson.M{"name": dummyName}

	testCases := []struct {
		name          string
		mockRetError  error
		expectedError error
	}{
		{"Success", nil, nil},
		{"TopicNotFound", mgo.ErrNotFound, errors.NotFound{}},
		{"DbFailed", errors.Unknown{}, errors.InternalServerError{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gomock.InOrder(
				mgoCollectionMockObj.EXPECT().Remove(dummyQuery).Return(tc.mockRetError),
			)

			err := Handler.DeleteTopic(dummyName)
			if reflect.TypeOf(err) != reflect.TypeOf(tc.expectedError) {
				t.Errorf("Expected Error: %s, Actual: %s", tc.expectedError, err)
			}
		})
	}
}

func TestCallReadTopicAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mgoCollectionMockObj := mgoMock.NewMockCollection(ctrl)
	mgoQueryMockObj := mgoMock.NewMockQuery(ctrl)

	// pass mockObj to a real object.
	mgoTopicCollection = mgoCollectionMockObj

	testCases := []struct {
		name          string
		mockRetError  error
		expectedError error
	}{
		{"Success", nil, nil},
		{"DbFailed", errors.Unknown{}, errors.InternalServerError{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			outTopics := []Topic{{Name: "/a", Endpoint: "0.0.0.0:1234", Datamodel: "test_0.0.1"}}

			gomock.InOrder(
				mgoCollectionMockObj.EXPECT().Find(nil).Return(mgoQueryMockObj), // nil query to read all
				mgoQueryMockObj.EXPECT().All(gomock.Any()).SetArg(0, outTopics).Return(tc.mockRetError),
			)

			_, err := Handler.ReadTopicAll()
			if reflect.TypeOf(err) != reflect.TypeOf(tc.expectedError) {
				t.Errorf("Expected Error: %s, Actual: %s", tc.expectedError, err)
			}
		})
	}
}

func TestCallReadTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mgoCollectionMockObj := mgoMock.NewMockCollection(ctrl)
	mgoQueryMockObj := mgoMock.NewMockQuery(ctrl)

	// pass mockObj to a real object.
	mgoTopicCollection = mgoCollectionMockObj

	testCases := []struct {
		name          string
		topicName     string
		hierarchical  bool
		mockRetError  error
		expectedError error
	}{
		{"Success", "/a", false, nil, nil},
		{"Success_Hierarchical", "/a", true, nil, nil},
		{"Success_Wildcard", "/a/*", false, nil, errors.InvalidQuery{}}, // @TODO: Since Wildcard feature is not implemented yet, the error is expteced.
		{"InvalidQuery_HierarchicalAndWildcard", "/a/*", true, nil, errors.InvalidQuery{}},
		{"DbFailed", "/a", false, errors.Unknown{}, errors.InternalServerError{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name != "InvalidQuery_HierarchicalAndWildcard" && tc.name != "Success_Wildcard" {
				dummyQuery := bson.M{}
				if tc.hierarchical {
					name := tc.topicName
					len := len(name)
					for i := 0; i < len; i++ {
						matched, _ := regexp.MatchString("[^a-zA-Z0-9]", string(name[i]))
						if matched {
							name = name[:i] + "\\" + name[i:]
							len++
							i++
						}
					}

					pattern := "^" + name + "$|^" + name + "/"
					dummyQuery = bson.M{"name": bson.RegEx{Pattern: pattern}}
				} else {
					dummyQuery = bson.M{"name": tc.topicName}
				}

				outTopics := []Topic{{Name: "/a", Endpoint: "0.0.0.0:1234", Datamodel: "test_0.0.1"}}

				gomock.InOrder(
					mgoCollectionMockObj.EXPECT().Find(dummyQuery).Return(mgoQueryMockObj),
					mgoQueryMockObj.EXPECT().All(gomock.Any()).SetArg(0, outTopics).Return(tc.mockRetError),
				)
			}
			_, err := Handler.ReadTopic(tc.topicName, tc.hierarchical)
			if reflect.TypeOf(err) != reflect.TypeOf(tc.expectedError) {
				t.Errorf("Expected Error: %s, Actual: %s", tc.expectedError, err)
			}
		})
	}
}
