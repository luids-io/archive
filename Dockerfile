FROM golang:alpine as build-env
ARG arch=amd64

# Install git and certificates
RUN apk update && apk add --no-cache git make ca-certificates && update-ca-certificates

# create user for service
RUN adduser -D -g '' luarchive

WORKDIR /app
## dependences
COPY go.mod .
COPY go.sum .
RUN go mod download

## build
COPY . .
RUN make binaries SYSTEM="GOOS=linux GOARCH=${arch}"

## create docker
FROM scratch

LABEL maintainer="Luis Guillén Civera <luisguillenc@gmail.com>"

# Import the user and group files from the builder.
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /etc/passwd /etc/passwd

COPY --from=build-env /app/bin/luarchive /bin/
COPY --from=build-env /app/configs/docker/luarchive.toml /etc/luids/archive/

USER luarchive

EXPOSE 5821
VOLUME [ "/etc/luids/archive" ]
CMD [ "/bin/luarchive" ]
