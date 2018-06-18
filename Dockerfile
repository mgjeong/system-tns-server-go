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
# Docker image for "tns-server"
FROM alpine:3.6

# environment variables
ENV APP_DIR=/tns
ENV APP=tns-server
ENV APP_PORT=48323

# install MongoDB
RUN apk add --no-cache mongodb && \
    rm -rf /var/cache/apk/*

# make mongodb volume
RUN mkdir -p /data/db
VOLUME /data/db

# copy files
COPY $APP run.sh $APP_DIR/
COPY ./config $APP_DIR/config

# expose tns-server rest api port
EXPOSE $APP_PORT

# set the working directory
WORKDIR $APP_DIR

# kick off the tns-server container
CMD ["sh", "run.sh"]
