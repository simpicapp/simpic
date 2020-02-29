FROM node:13-stretch AS frontend

RUN mkdir /tmp/site
WORKDIR /tmp/site
ADD . /tmp/site

RUN cd frontend \
    && npm install \
    && sed -i "s/gitSHA: ''/gitSHA: '$(git rev-parse --short HEAD)'/" src/main.ts \
    && npm run build


FROM golang:1.14.0 AS backend
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go install -ldflags "-X github.com/simpicapp/simpic/internal.GitSHA=$(git rev-parse --short HEAD)" github.com/simpicapp/simpic/cmd/serve


FROM scratch
COPY --from=frontend /tmp/site/frontend/dist /dist
COPY --from=backend /go/bin/serve /serve
COPY migrations /migrations
VOLUME /data
ENTRYPOINT ["/serve"]
