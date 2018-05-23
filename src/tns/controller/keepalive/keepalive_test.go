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
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
	"time"
	"tns/commons/errors"
	topicDbMock "tns/db/topic/mocks"
)

var Handler Command

func init() {
	Handler = Executor{}
}

func TestCallInitKeepAlive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicDbMockObj := topicDbMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicDbExecutor = topicDbMockObj

	var dummyInterval uint = 10

	testCases := []struct {
		name        string
		dummyTopics []map[string]interface{}
		dummyError  error
	}{
		{"Success", []map[string]interface{}{{"name": "/a", "endpoint": "0.0.0.0:1234", "datamodel": "test_0.0.1"}}, nil},
		{"DbFailed", nil, errors.Unknown{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gomock.InOrder(
				topicDbMockObj.EXPECT().ReadTopicAll().Return(tc.dummyTopics, tc.dummyError),
			)

			err := Handler.InitKeepAlive(dummyInterval)
			if err != tc.dummyError {
				t.Fail()
			}
			if err == nil {
				dummyName := tc.dummyTopics[0]["name"].(string)
				if _, exist := kaInfo.table[dummyName]; !exist {
					t.Fail()
				}
				if kaInfo.interval != dummyInterval {
					t.Fail()
				}
			}
		})
	}
}

func TestCallAddTopicAndDeleteTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicDbMockObj := topicDbMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicDbExecutor = topicDbMockObj

	dummyTopicName := "/a"

	Handler.AddTopic(dummyTopicName)
	if _, exist := kaInfo.table[dummyTopicName]; !exist {
		t.Errorf("Topic does not exist: %s", dummyTopicName)
	}

	Handler.DeleteTopic(dummyTopicName)
	if _, exist := kaInfo.table[dummyTopicName]; exist {
		t.Errorf("Topic exists: %s", dummyTopicName)
	}
}

func TestCallHandlePing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicDbMockObj := topicDbMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicDbExecutor = topicDbMockObj

	dummyBodyString := `{"topic_names":["/a"]}`

	Handler.AddTopic("/a")

	_, err := Handler.HandlePing(dummyBodyString)
	if err != nil {
		t.Errorf("HandlePing returned an error: %s", err.Error())
	}
}

func TestCallHandlePingExpectedError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicDbMockObj := topicDbMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicDbExecutor = topicDbMockObj

	testCases := []struct {
		name            string
		dummyBodyString string
		expectedResp    map[string]interface{}
		expectedError   error
	}{
		{"InvalidJsonFormat", `invalid_json[`, nil, errors.InvalidJSON{}},
		{"InvalidKey", `{"invalid_key":["/a","/b"]}`, nil, errors.InvalidParam{}},
		{"InvalidValue", `{"topic_names":["/a",1]}`, nil, errors.InvalidParam{}},
		{"TopicNameNotFound", `{"topic_names":["/a","/b"]}`, map[string]interface{}{"topic_names": []string{"/b"}}, errors.NotFound{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			Handler.AddTopic("/a") // add "/a"

			resp, err := Handler.HandlePing(tc.dummyBodyString)

			if reflect.TypeOf(err) != reflect.TypeOf(tc.expectedError) {
				t.Errorf("Expected Error: %s, Actual: %s", tc.expectedError, err)
			}

			if isEqual := reflect.DeepEqual(resp, tc.expectedResp); !isEqual {
				t.Errorf("Expected Resp: %s, Actual: %s", tc.expectedResp, resp)
			}
		})
	}

}

func TestCallGetInterval(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicDbMockObj := topicDbMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicDbExecutor = topicDbMockObj

	var dummyInterval uint = 100
	expectedRetVal := dummyInterval / kaPingFrequency

	topicDbMockObj.EXPECT().ReadTopicAll()

	Handler.InitKeepAlive(dummyInterval)

	interval := Handler.GetInterval()
	if interval != expectedRetVal {
		t.Errorf("Expected val: %d, Actual: %d", expectedRetVal, interval)
	}
}

func TestKeepAliveTimerLoopCalled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicDbMockObj := topicDbMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicDbExecutor = topicDbMockObj

	var dummyTopics = []map[string]interface{}{{"name": "/a", "endpoint": "0.0.0.0:1234", "datamodel": "test_0.0.1"}}
	var dummyInterval uint = 1
	const waitingTimeForTopicExpired = 2

	gomock.InOrder(
		topicDbMockObj.EXPECT().ReadTopicAll().Return(dummyTopics, nil),
		topicDbMockObj.EXPECT().DeleteTopic("/a").Return(nil),
		topicDbMockObj.EXPECT().DeleteTopic("/tmp").Return(errors.Unknown{}),
	)

	// add "/a"
	Handler.InitKeepAlive(dummyInterval)

	// "/a" will be deleted after dummyInterval seconds
	time.Sleep(waitingTimeForTopicExpired * time.Second)

	// add "/tmp"
	Handler.AddTopic("/tmp")

	// topicDbMock will return error
	time.Sleep(waitingTimeForTopicExpired * time.Second)
}
