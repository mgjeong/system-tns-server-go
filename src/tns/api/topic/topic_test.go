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
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"tns/commons/errors"
	topicControllerMock "tns/controller/topic/mocks"
)

const topicUrl = "/api/v1/tns/topic"

var (
	testBodyString = `{"topic":{"datamodel":"test_0.0.1","endpoint":"0.0.0.0:1234","name":"/a"}}`

	testBody = map[string]interface{}{
		"topic": map[string]interface{}{
			"name": "/a", "endpoint": "0.0.0.0:1234", "datamodel": "test_0.0.1"}}
)

var Handler Command

func init() {
	Handler = RequestHandler{}
}

func TestCallHandleWithInvalidRequest(t *testing.T) {
	// Mock is not necessary for this test

	testCases := []struct {
		name         string
		method       string
		url          string
		expectedCode int
	}{
		{"InvalidUrl", "POST", topicUrl + "/invalid", http.StatusNotFound},
		{"InvalidMethod_Put", "PUT", topicUrl, http.StatusBadRequest},
		{"EmptyParameter_Post", "POST", topicUrl, http.StatusBadRequest},
		{"InvalidQuery_Get_MultiValue", "GET", topicUrl + "?name=a&name=b", http.StatusBadRequest},
		{"InvalidQuery_Get_InvalidValue", "GET", topicUrl + "?hierarchical=invalid", http.StatusBadRequest},
		{"InvalidQuery_Get_InvalidQuery", "GET", topicUrl + "?key=value", http.StatusBadRequest},
		{"InvalidQuery_Delete_MultiValue", "DELETE", topicUrl + "?name=a&name=b", http.StatusBadRequest},
		{"InvalidQuery_Delete_Hierarchical", "DELETE", topicUrl + "?hierarchical=yes", http.StatusBadRequest},
		{"InvalidQuery_Delete_InvalidQuery", "DELETE", topicUrl + "?key=value", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, nil)
			w := httptest.NewRecorder()

			Handler.Handle(w, req)

			if w.Code != tc.expectedCode {
				t.Errorf("Expected Code: %s, Actual: %s", http.StatusText(tc.expectedCode), http.StatusText(w.Code))
			}
		})
	}
}

func TestCallHandlePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicCtrlrMockObj := topicControllerMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicExecutor = topicCtrlrMockObj

	expectedResp := map[string]interface{}{"ka_interval": 200}
	expectedRespByte, _ := json.Marshal(expectedResp)

	gomock.InOrder(
		topicCtrlrMockObj.EXPECT().CreateTopic(testBodyString).Return(expectedResp, nil),
	)

	body, _ := json.Marshal(testBody)
	req := httptest.NewRequest("POST", topicUrl, bytes.NewReader(body))
	w := httptest.NewRecorder()

	Handler.Handle(w, req)

	expectedCode := http.StatusCreated
	if w.Code != expectedCode {
		t.Errorf("Expected Code: %s, Actual: %s", http.StatusText(expectedCode), http.StatusText(w.Code))
	}
	if 0 != bytes.Compare(w.Body.Bytes(), expectedRespByte) {
		t.Errorf("Expected body: %s, Actual: %s", body, w.Body.Bytes())
	}
}

func TestCallHandlePostWithEmptyBody(t *testing.T) {
	req := httptest.NewRequest("POST", topicUrl, nil)
	w := httptest.NewRecorder()

	Handler.Handle(w, req)

	expectedCode := http.StatusBadRequest
	if w.Code != expectedCode {
		t.Errorf("Expected Code: %s, Actual: %s", http.StatusText(expectedCode), http.StatusText(w.Code))
	}
}

func TestCallHandlePostWithInvalidPayload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicCtrlrMockObj := topicControllerMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicExecutor = topicCtrlrMockObj

	gomock.InOrder(
		topicCtrlrMockObj.EXPECT().CreateTopic(testBodyString).Return(nil, errors.InvalidParam{}),
	)

	body, _ := json.Marshal(testBody)
	req := httptest.NewRequest("POST", topicUrl, bytes.NewReader(body))
	w := httptest.NewRecorder()

	Handler.Handle(w, req)

	expectedCode := http.StatusBadRequest
	if w.Code != expectedCode {
		t.Errorf("Expected Code: %s, Actual: %s", http.StatusText(expectedCode), http.StatusText(w.Code))
	}
}

func TestCallHandleGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicCtrlrMockObj := topicControllerMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicExecutor = topicCtrlrMockObj

	expectedResp := map[string]interface{}{"topics": []map[string]interface{}{{"name": "/a", "endpoint": "0.0.0.0:1234", "datamodel": "test_0.0.1"}}}
	expectedRespByte, _ := json.Marshal(expectedResp)

	name := "/a"
	hierarchical := "yes"

	gomock.InOrder(
		topicCtrlrMockObj.EXPECT().ReadTopic(name, true).Return(expectedResp, nil),
	)

	req := httptest.NewRequest("GET", topicUrl+"?name="+name+"&hierarchical="+hierarchical, nil)
	w := httptest.NewRecorder()

	Handler.Handle(w, req)

	expectedCode := http.StatusOK
	if w.Code != expectedCode {
		t.Errorf("Expected Code: %s, Actual: %s", http.StatusText(expectedCode), http.StatusText(w.Code))
	}
	body, _ := json.Marshal(expectedResp)
	if 0 != bytes.Compare(w.Body.Bytes(), expectedRespByte) {
		t.Errorf("Expected body: %s, Actual: %s", body, w.Body.Bytes())
	}
}

func TestCallHandleGetWithNonExistTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicCtrlrMockObj := topicControllerMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicExecutor = topicCtrlrMockObj

	name := "/a"
	hierarchical := "no"

	gomock.InOrder(
		topicCtrlrMockObj.EXPECT().ReadTopic(name, false).Return(nil, errors.NotFound{}),
	)

	req := httptest.NewRequest("GET", topicUrl+"?name="+name+"&hierarchical="+hierarchical, nil)
	w := httptest.NewRecorder()

	Handler.Handle(w, req)

	expectedCode := http.StatusNotFound
	if w.Code != expectedCode {
		t.Errorf("Expected Code: %s, Actual: %s", http.StatusText(expectedCode), http.StatusText(w.Code))
	}
}

func TestCallHandleDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicCtrlrMockObj := topicControllerMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicExecutor = topicCtrlrMockObj

	name := "/a"

	gomock.InOrder(
		topicCtrlrMockObj.EXPECT().DelteTopic(name).Return(nil),
	)

	req := httptest.NewRequest("DELETE", topicUrl+"?name="+name, nil)
	w := httptest.NewRecorder()

	Handler.Handle(w, req)

	expectedCode := http.StatusOK
	if w.Code != expectedCode {
		t.Errorf("Expected Code: %s, Actual: %s", http.StatusText(expectedCode), http.StatusText(w.Code))
	}
}

func TestCallHandleDeleteWithNonExistTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicCtrlrMockObj := topicControllerMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicExecutor = topicCtrlrMockObj

	name := "/a"

	gomock.InOrder(
		topicCtrlrMockObj.EXPECT().DelteTopic(name).Return(errors.NotFound{}),
	)

	req := httptest.NewRequest("DELETE", topicUrl+"?name="+name, nil)
	w := httptest.NewRecorder()

	Handler.Handle(w, req)

	expectedCode := http.StatusNotFound
	if w.Code != expectedCode {
		t.Errorf("Expected Code: %s, Actual: %s", http.StatusText(expectedCode), http.StatusText(w.Code))
	}
}
