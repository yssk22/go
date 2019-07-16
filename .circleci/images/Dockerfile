FROM golang:1.12.5 AS go

FROM node:8.16 AS node

FROM google/cloud-sdk:224.0.0 AS cloud-sdk

COPY --from=go /usr/local/go /usr/local/go
COPY --from=node /usr/local/bin/node /usr/local/bin/node
COPY --from=node /usr/local/lib/node_modules/ /usr/local/lib/node_modules/

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin:/usr/local/node_modules/.bin/
RUN go get github.com/jstemmer/go-junit-report
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update && apt-get install -y yarn make wget curl
RUN yarn global add prettier