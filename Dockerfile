# syntax=docker/dockerfile:1
FROM golang:1.20-alpine AS build
WORKDIR /app/caddy-json-ui
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -tags=nomsgpack -o /caddy-json-ui .

FROM scratch
WORKDIR /
COPY --from=build /caddy-json-ui /caddy-json-ui
COPY ./public /public
EXPOSE 8045
# HEALTHCHECK --interval=10s --timeout=60s CMD curl --fail http://127.0.0.1:8045/api/healthy || exit 1
ENTRYPOINT ["/caddy-json-ui"]