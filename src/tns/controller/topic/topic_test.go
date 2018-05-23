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
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
	"tns/commons/errors"
	kaControllerMock "tns/controller/keepalive/mocks"
	topicDbMock "tns/db/topic/mocks"
)

var Handler Command

func init() {
	Handler = Executor{}
}

func TestCallCreateTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicDbMockObj := topicDbMock.NewMockCommand(ctrl)
	kaControllerMockObj := kaControllerMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicDbExecutor = topicDbMockObj
	keepaliveExecutor = kaControllerMockObj

	dummyBodyString := `{"topic":{"name":"/a","endpoint":"0.0.0.0:1234","datamodel":"test_0.0.1"}}`
	dummyTopic := map[string]interface{}{"name": "/a", "endpoint": "0.0.0.0:1234", "datamodel": "test_0.0.1"}
	interval := uint(10)
	expectedResp := map[string]interface{}{"ka_interval": interval}

	gomock.InOrder(
		topicDbMockObj.EXPECT().CreateTopic(dummyTopic).Return(nil),
		kaControllerMockObj.EXPECT().AddTopic("/a"),
		kaControllerMockObj.EXPECT().GetInterval().Return(interval),
	)

	resp, err := Handler.CreateTopic(dummyBodyString)
	if err != nil {
		t.Errorf("CreateTopic returned an error: %s", err.Error())
	}
	if isEqual := reflect.DeepEqual(resp, expectedResp); !isEqual {
		t.Errorf("Expected Resp: %s, Actual: %s", expectedResp, resp)
	}
}

func TestCallCreateTopicWithInvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicDbMockObj := topicDbMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicDbExecutor = topicDbMockObj

	testCases := []struct {
		name            string
		dummyBodyString string
		expectedError   error
	}{
		{"InvalidJson", `{invalidJson[}`, errors.InvalidJSON{}},
		{"InvalidParam_topic", `{"invalid":{"name":"/a","endpoint":"0.0.0.0:1234","datamodel":"test_0.0.1"}}`, errors.InvalidParam{}},
		{"InvalidParam_name", `{"topic":{"datamodel":"test_0.0.1","endpoint":"0.0.0.0:1234"}}`, errors.InvalidParam{}},
		{"Conflict", `{"topic":{"name":"/a","endpoint":"0.0.0.0:1234","datamodel":"test_0.0.1"}}`, errors.Conflict{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// mock will be called only for the conflict error case.
			if tc.name == "Conflict" {
				dummyTopic := map[string]interface{}{"name": "/a", "endpoint": "0.0.0.0:1234", "datamodel": "test_0.0.1"}
				topicDbMockObj.EXPECT().CreateTopic(dummyTopic).Return(errors.Conflict{})
			}

			_, err := Handler.CreateTopic(tc.dummyBodyString)
			if reflect.TypeOf(err) != reflect.TypeOf(tc.expectedError) {
				t.Errorf("Expected Error: %s, Actual: %s", tc.expectedError, err)
			}
		})
	}
}

func TestCallReadTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicDbMockObj := topicDbMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicDbExecutor = topicDbMockObj

	topics := []map[string]interface{}{{"name": "/a"}}
	successResp := map[string]interface{}{"topics": topics}
	hierarchical := false

	testCases := []struct {
		name          string
		topicName     string
		mockRetTopics []map[string]interface{}
		mockRetError  error
		expectedResp  map[string]interface{}
		expectedError error
	}{
		{"Success_Single", "/a", topics, nil, successResp, nil},
		{"Success_All", "", topics, nil, successResp, nil},
		{"DbFailed", "/a", nil, errors.InternalServerError{}, nil, errors.InternalServerError{}},
		{"NotFound", "", nil, nil, nil, errors.NotFound{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.topicName == "" {
				topicDbMockObj.EXPECT().ReadTopicAll().Return(tc.mockRetTopics, tc.mockRetError)
			} else {
				topicDbMockObj.EXPECT().ReadTopic(tc.topicName, hierarchical).Return(tc.mockRetTopics, tc.mockRetError)
			}

			resp, err := Handler.ReadTopic(tc.topicName, hierarchical)
			if reflect.TypeOf(err) != reflect.TypeOf(tc.expectedError) {
				t.Errorf("Expected Error: %s, Actual: %s", tc.expectedError, err)
			}
			if isEqual := reflect.DeepEqual(resp, tc.expectedResp); !isEqual {
				t.Errorf("Expected Resp: %s, Actual: %s", tc.expectedResp, resp)
			}
		})
	}
}

func TestCallDeleteTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicDbMockObj := topicDbMock.NewMockCommand(ctrl)
	kaControllerMockObj := kaControllerMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicDbExecutor = topicDbMockObj
	keepaliveExecutor = kaControllerMockObj

	topicName := "/a"

	testCases := []struct {
		name          string
		mockRetError  error
		expectedError error
	}{
		{"Success", nil, nil},
		{"DbFailed", errors.NotFound{}, errors.NotFound{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			topicDbMockObj.EXPECT().DeleteTopic(topicName).Return(tc.mockRetError)

			// kaMock will be called only for the success case.
			if tc.name == "Success" {
				kaControllerMockObj.EXPECT().DeleteTopic(topicName)
			}

			err := Handler.DeleteTopic(topicName)
			if reflect.TypeOf(err) != reflect.TypeOf(tc.expectedError) {
				t.Errorf("Expected Error: %s, Actual: %s", tc.expectedError, err)
			}
		})
	}
}
