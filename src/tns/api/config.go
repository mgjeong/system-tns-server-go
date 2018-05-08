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

package api

import (
	"github.com/BurntSushi/toml"
	"os"
	"tns/commons/logger"
)

type Config struct {
	Server struct {
		Ip                string
		Port              uint
		KeepAliveInterval uint
	}
	Database struct {
		Ip         string
		Port       uint
		Name       string
		Collection string
	}
}

// Read and parse the configuration file
func (c *Config) Read(filePath string) error {
	logger.Logging(logger.DEBUG, "File path: "+filePath)

	if _, err := os.Stat(filePath); err != nil {
		logger.Logging(logger.ERROR, "Cannot Open the file: ", err.Error())
		return err
	}

	if _, err := toml.DecodeFile(filePath, &c); err != nil {
		logger.Logging(logger.ERROR, "DecodeFile failed: "+err.Error())
		return err
	}

	return nil
}
