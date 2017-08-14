FROM golang:1.8.0 AS go

FROM google/cloud-sdk:164.0.0 AS cloud-sdk
COPY --from=go /usr/local/go /usr/local/go

ENV GOROOT /usr/local/go
ENV PATH $PATH:/usr/local/go/bin
RUN go get github.com/jstemmer/go-junit-report
