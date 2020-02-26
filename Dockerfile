FROM node:13-stretch AS parcel

RUN npm install -g parcel-bundler
RUN mkdir /tmp/site
ADD . /tmp/site
RUN cd /tmp/site && npm install && parcel build frontend/index.html


FROM golang:1.14.0 AS build
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GO111MODULE=on go install github.com/simpicapp/simpic/cmd/serve

FROM scratch
COPY --from=parcel /tmp/site/dist /dist
COPY --from=build /go/bin/serve /serve
COPY migrations /migrations
VOLUME /data
ENTRYPOINT ["/serve"]
