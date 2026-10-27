FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /api ./cmd/api
FROM debian:bookworm-slim
COPY --from=build /api /usr/local/bin/api
USER 65534:65534
EXPOSE 8080
CMD ["api"]
