FROM node:16-stretch AS frontend

RUN mkdir /tmp/site
WORKDIR /tmp/site
ADD . /tmp/site

RUN cd frontend \
    && npm install \
    && sed -i "s/gitSHA: ''/gitSHA: '$(git rev-parse --short HEAD)'/" src/main.ts \
    && npm run build


FROM simpicapp/golang-imagemagick

WORKDIR /go/src/app
COPY . .

RUN go install -ldflags "-X github.com/simpicapp/simpic/internal.GitSHA=$(git rev-parse --short HEAD)" github.com/simpicapp/simpic/cmd/serve

# Generate licence information
RUN go get github.com/google/go-licenses && go-licenses save ./... --save_path=/notices

WORKDIR /
COPY --from=frontend /tmp/site/frontend/dist /dist
COPY migrations /migrations
VOLUME /data
ENTRYPOINT ["/go/bin/serve"]
