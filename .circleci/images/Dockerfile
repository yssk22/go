FROM circleci/golang:1.12.5 AS go

FROM circleci/node:8.16 AS node

FROM google/cloud-sdk:254.0.0 AS cloud-sdk

COPY --from=go /usr/local/go /usr/local/go
COPY --from=node /usr/local/bin/node /usr/local/bin/node
COPY --from=node /usr/local/lib/node_modules/ /usr/local/lib/node_modules/

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin:/usr/local/node_modules/.bin/
RUN go get github.com/jstemmer/go-junit-report
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update && apt-get install -y yarn make wget curl sudo
RUN yarn global add prettier

## from https://github.com/circleci/circleci-images/blob/staging/shared/images/Dockerfile-basic.template
RUN mkdir -p /etc/sudoers.d/
RUN groupadd --gid 3434 circleci \
    && useradd --uid 3434 --gid circleci --shell /bin/bash --create-home circleci \
    && echo 'circleci ALL=NOPASSWD: ALL' >> /etc/sudoers.d/50-circleci \
    && echo 'Defaults    env_keep += "DEBIAN_FRONTEND"' >> /etc/sudoers.d/env_keep

RUN mkdir -p /go
RUN chown -R circleci:circleci /go

USER circleci

CMD ["/bin/sh"]