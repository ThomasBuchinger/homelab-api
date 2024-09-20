FROM cgr.dev/chainguard/node:latest AS ui

ADD ./ui/package.json ./ui/package-lock.json ./
RUN npm ci

ADD ./ui/public ./public
ADD ./ui/app ./app
ADD ./ui/tests ./tests
ADD ./ui/.eslintrc.json ./ui/next.config.js ./ui/tsconfig.json ./ui/jest.config.js ./
ADD ./README.md ./
RUN npm run test
RUN npm run build

FROM golang:latest AS api
ENV CGO_ENABLED=0

WORKDIR /go/github.com/thomasbuchinger/homelab-api
ADD ./go.mod ./go.sum ./

ADD /geoip ./geoip
ADD ./cmd cmd/
ADD ./pkg pkg/
ADD ./README.md ./
RUN go get ./...
RUN go build ./cmd/homelab-api
RUN go build ./cmd/public-api
RUN go build ./cmd/syncthing-helper
RUN go build ./cmd/copy-geoip-for-envoy
RUN go test ./...

FROM scratch AS app
ENV MODE=dev

WORKDIR /app
ENV GEOIP_DATABASE=/geoip/GeoLite2-City.mmdb

COPY /geoip /geoip
COPY --from=ui /app/out/ ui/out
COPY --from=api /go/github.com/thomasbuchinger/homelab-api/homelab-api /app
COPY --from=api /go/github.com/thomasbuchinger/homelab-api/public-api /app
COPY --from=api /go/github.com/thomasbuchinger/homelab-api/syncthing-helper /app
COPY --from=api /go/github.com/thomasbuchinger/homelab-api/copy-geoip-for-envoy /app

ENTRYPOINT ["/app/homelab-api"]