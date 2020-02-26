FROM node:13-stretch AS parcel

RUN npm install -g parcel-bundler
RUN mkdir /tmp/site
WORKDIR /tmp/site
ADD . /tmp/site

RUN npm install \
    && sed -i "s/gitSHA: ''/gitSHA: '$(git rev-parse --short HEAD)'/" frontend/js/simpic.js \
    && sed -i "s/gitTag: ''/gitTag: '$(git describe --tags --always)'/" frontend/js/simpic.js \
    && parcel build frontend/index.html --no-source-maps


FROM golang:1.14.0 AS build
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go install -ldflags "-X github.com/simpicapp/simpic/internal.GitTag=$(git describe --tags --always) -X github.com/simpicapp/simpic/internal.GitSHA=$(git rev-parse --short HEAD)" github.com/simpicapp/simpic/cmd/serve

FROM scratch
COPY --from=parcel /tmp/site/dist /dist
COPY --from=build /go/bin/serve /serve
COPY migrations /migrations
VOLUME /data
ENTRYPOINT ["/serve"]
