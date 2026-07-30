# Build dos dois binários do backend a partir de um único Dockerfile.
# Use `target: api` ou `target: ws` para escolher o serviço.
#
# GO_VERSION vem do .env (via docker-compose) ou de --build-arg. O default abaixo
# é só um fallback: ajuste-o para a versão que `go version` reporta na sua máquina.

ARG GO_VERSION=1.24

# ---------- build ----------
FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /src

# Camada de dependências separada: só refaz o download quando go.mod/go.sum mudam
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api \
 && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/ws  ./cmd/ws

# ---------- base de runtime ----------
FROM alpine:3 AS runtime-base
RUN apk add --no-cache ca-certificates tzdata \
 && adduser -D -u 10001 squadr
USER squadr

# ---------- API REST ----------
FROM runtime-base AS api
COPY --from=build /out/api /usr/local/bin/squadr-api
EXPOSE 8080
ENTRYPOINT ["squadr-api"]

# ---------- serviço WebSocket ----------
FROM runtime-base AS ws
COPY --from=build /out/ws /usr/local/bin/squadr-ws
EXPOSE 8081
ENTRYPOINT ["squadr-ws"]
