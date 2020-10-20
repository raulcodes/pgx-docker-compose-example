FROM golang:1.14-alpine AS build

WORKDIR /src/
COPY main.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/pgx-docker-compose

FROM scratch
COPY --from=build /bin/pgx-docker-compose /bin/pgx-docker-compose
ENTRYPOINT ["/bin/pgx-docker-compose"]