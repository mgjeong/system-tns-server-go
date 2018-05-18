/*******************************************************************************
 * Copyright 2017 Samsung Electronics All Rights Reserved.
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

package api

import (
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	kaApiMock "tns/api/keepalive/mocks"
	topicApiMock "tns/api/topic/mocks"
)

func TestCallServeHTTPWithInvalidUrl(t *testing.T) {
	// Mock is not necessary for this test

	testCases := []struct {
		name string
		url  string
	}{
		{"WithoutBaseUrl", "/tns/topic"},
		{"NonExistUrl", "/api/v1/tns/invalid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", tc.url, nil)
			w := httptest.NewRecorder()

			Handler.ServeHTTP(w, req)

			expectedCode := http.StatusNotFound
			if w.Code != expectedCode {
				t.Errorf("Expected Code: %s, Actual: %s", http.StatusText(expectedCode), http.StatusText(w.Code))
			}
		})
	}
}

func TestCallServeHTTPWithTopicUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicApiMockObj := topicApiMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	topicHandler = topicApiMockObj

	req := httptest.NewRequest("POST", "/api/v1/tns/topic", nil)
	w := httptest.NewRecorder()

	gomock.InOrder(
		topicApiMockObj.EXPECT().Handle(w, req),
	)

	Handler.ServeHTTP(w, req)
}

func TestCallServeHTTPWithKeepAliveUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kaApiMockObj := kaApiMock.NewMockCommand(ctrl)

	// pass mockObj to a real object.
	keepAliveHandler = kaApiMockObj

	req := httptest.NewRequest("POST", "/api/v1/tns/keepalive", nil)
	w := httptest.NewRecorder()

	gomock.InOrder(
		kaApiMockObj.EXPECT().Handle(w, req),
	)

	Handler.ServeHTTP(w, req)
}

func TestCallRead(t *testing.T) {
	tomlFile, err := os.Create("test.toml")
	if err != nil {
		t.Error("Create failed")
	}
	defer os.Remove(tomlFile.Name())

	_, err = tomlFile.Write([]byte("[server]\nip = \"0.0.0.0\""))
	if err != nil {
		t.Error("Write failed")
	}

	config = Config{}
	config.Read(tomlFile.Name())
}

func TestCallRead_OpenFailed(t *testing.T) {
	config = Config{}
	err := config.Read("nonExistsFile")
	if err == nil {
		t.Error("Read did not return an error")
	}
}

func TestCallRead_DecodeFailed(t *testing.T) {
	tomlFile, err := os.Create("test.toml")
	if err != nil {
		t.Error("Create failed")
	}
	defer os.Remove(tomlFile.Name())

	_, err = tomlFile.Write([]byte("[server]\nip = invalid toml"))
	if err != nil {
		t.Error("Write failed")
	}

	config = Config{}
	err = config.Read(tomlFile.Name())
	if err == nil {
		t.Error("Read did not return an error")
	}
}
