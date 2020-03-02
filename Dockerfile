FROM node:13-stretch AS frontend

RUN mkdir /tmp/site
WORKDIR /tmp/site
ADD . /tmp/site

RUN cd frontend \
    && npm install \
    && sed -i "s/gitSHA: ''/gitSHA: '$(git rev-parse --short HEAD)'/" src/main.ts \
    && npm run build


FROM golang:1.14.0-buster

ARG DEBIAN_FRONTEND=noninteractive
ARG IMAGEMAGICK_VERSION=7.0.9-27

RUN apt-get update \
    && apt-get --no-install-recommends -y install -y wget build-essential pkg-config \
    && apt-get --no-install-recommends -y install libjpeg-dev libpng-dev libtiff-dev libgif-dev libraw-dev libwebp-dev \
    && wget https://github.com/ImageMagick/ImageMagick/archive/${IMAGEMAGICK_VERSION}.tar.gz \
    && tar xvzf ${IMAGEMAGICK_VERSION}.tar.gz \
    && cd ImageMagick* \
    && ./configure \
        --without-magick-plus-plus \
        --without-perl \
        --disable-openmp \
        --with-gvc=no \
        --disable-docs \
    && make -j$(nproc) \
    && make install \
    && ldconfig /usr/local/lib \
    && rm -rf /var/lib/apt/lists/* \
    && cd .. \
    && rm -rf ImageMagick*

WORKDIR /go/src/app
COPY . .

RUN go install -ldflags "-X github.com/simpicapp/simpic/internal.GitSHA=$(git rev-parse --short HEAD)" github.com/simpicapp/simpic/cmd/serve

WORKDIR /
COPY --from=frontend /tmp/site/frontend/dist /dist
COPY migrations /migrations
VOLUME /data
ENTRYPOINT ["/go/bin/serve"]
