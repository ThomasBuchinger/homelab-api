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

ADD ./pkg pkg/
ADD ./README.md ./main.go ./
RUN go get
RUN go build .

FROM scratch AS app
ENV MODE=dev
ENV GIN_MODE=debug
# ENV GIN_MODE=release

WORKDIR /app
COPY --from=ui /app/out/ ui/out
COPY --from=api /go/github.com/thomasbuchinger/homelab-api/homelab-api /app

ENTRYPOINT ["/app/homelab-api"]