###############################################################################
# Copyright 2018 Samsung Electronics All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
###############################################################################
#!/bin/bash

EXECUTABLE_FILE_NAME="tns-server"

export GOPATH=$PWD

usage() {
    echo "  -c                          :  Remove all downloaded/output files"
    echo "  -h / --help                 :  Display help and exit"
}

function clean(){
    rm -rf $GOPATH/pkg
    rm -rf $GOPATH/src/github.com
    rm -rf $GOPATH/src/golang.org
    rm -rf $GOPATH/src/gopkg.in
    rm -f coverall.html
    rm -f ${EXECUTABLE_FILE_NAME}
    echo -e "Finished Cleaning"
}

function build(){
    CGO_ENABLED=0 GOOS=linux go build -o ${EXECUTABLE_FILE_NAME} -a -ldflags '-extldflags "-static"' src/main/main.go
    if [ $? -ne 0 ]; then
        echo -e "\n\033[31m"build fail"\033[0m"
        func_cleanup
        exit 1
    fi
}

function download_pkgs(){
    pkg_list=(
        "github.com/BurntSushi/toml"
        "gopkg.in/mgo.v2"
    )

    idx=1
    for pkg in "${pkg_list[@]}"; do
        echo -ne "(${idx}/${#pkg_list[@]}) go get $pkg"
        go get $pkg
        if [ $? -ne 0 ]; then
            echo -e "\n\033[31m"download fail"\033[0m"
            clean
            exit 1
        fi
        echo ": Done"
        idx=$((idx+1))
    done
}

process_cmd_args() {
    while [ "$#" -gt 0  ]; do
        case "$1" in
            -c)
                clean
                shift 1; exit 0
                ;;
            -h)
                usage; exit 0
                ;;
            --help)
                usage; exit 0
                ;;
            -*)
                echo -e "${RED}"
                echo "unknown option: $1" >&2;
                echo -e "${NO_COLOUR}"
                usage; exit 1
                ;;
            *)
                echo -e "${RED}"
                echo "unknown option: $1" >&2;
                echo -e "${NO_COLOUR}"
                usage; exit 1
                ;;
        esac
    done
}

process_cmd_args "$@"

echo -e "\n\033[33m"Start building of Topic Name Service"\033[0m"

echo -e "\nDownload dependent go-pkgs"
download_pkgs

echo -ne "\nMaking executable file of Topic Name Service"
build
echo ": Done"

echo -e "\n\033[33m"Succeed build of Topic Name Service"\033[0m"
