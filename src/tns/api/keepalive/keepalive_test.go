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
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"tns/commons/errors"
	keepaliveControllerMock "tns/controller/keepalive/mocks"
)

const (
	testBodyString = `{"topic_names":["/a"]}`
)

var testBody = map[string]interface{}{
	"topic_names": []string{"/a"},
}

var Handler Command

func init() {
	Handler = RequestHandler{}
}

func TestCallHandle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kaCtrlrMockObj := keepaliveControllerMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	keepaliveExecutor = kaCtrlrMockObj

	gomock.InOrder(
		kaCtrlrMockObj.EXPECT().HandlePing(testBodyString),
	)

	body, _ := json.Marshal(testBody)
	req := httptest.NewRequest("POST", "/api/v1/tns/keepalive", bytes.NewReader(body))
	w := httptest.NewRecorder()

	Handler.Handle(w, req)
}

func TestCallHandleWithInvalidRequest(t *testing.T) {
	// Mock is not necessary for this test

	testCases := []struct {
		name         string
		method       string
		url          string
		expectedCode int
	}{
		{"InvalidUrl", "POST", "/api/v1/tns/keepalive/invalid", http.StatusNotFound},
		{"InvalidMethod_Get", "GET", "/api/v1/tns/keepalive", http.StatusBadRequest},
		{"InvalidMethod_Put", "PUT", "/api/v1/tns/keepalive", http.StatusBadRequest},
		{"InvalidMethod_Delete", "DELETE", "/api/v1/tns/keepalive", http.StatusBadRequest},
		{"EmptyParameter", "POST", "/api/v1/tns/keepalive", http.StatusBadRequest},
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

func TestCallHandleWithNonExistTopicName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kaCtrlrMockObj := keepaliveControllerMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	keepaliveExecutor = kaCtrlrMockObj

	gomock.InOrder(
		kaCtrlrMockObj.EXPECT().HandlePing(testBodyString).Return(testBody, errors.NotFound{}),
	)

	body, _ := json.Marshal(testBody)
	req := httptest.NewRequest("POST", "/api/v1/tns/keepalive", bytes.NewReader(body))
	w := httptest.NewRecorder()

	Handler.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected Code: %s, Actual: %s", http.StatusText(http.StatusBadRequest), http.StatusText(w.Code))
	}
	if 0 != bytes.Compare(w.Body.Bytes(), body) {
		t.Errorf("Expected body: %s, Actual: %s", body, w.Body.Bytes())
	}
}

func TestCallHandleExpectAnyError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kaCtrlrMockObj := keepaliveControllerMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	keepaliveExecutor = kaCtrlrMockObj

	gomock.InOrder(
		kaCtrlrMockObj.EXPECT().HandlePing(testBodyString).Return(nil, errors.InternalServerError{}),
	)

	body, _ := json.Marshal(testBody)
	req := httptest.NewRequest("POST", "/api/v1/tns/keepalive", bytes.NewReader(body))
	w := httptest.NewRecorder()

	Handler.Handle(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected Code: %s, Actual: %s", http.StatusText(http.StatusInternalServerError), http.StatusText(w.Code))
	}
}
