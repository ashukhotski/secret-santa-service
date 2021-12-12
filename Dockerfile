FROM golang:alpine AS build

RUN mkdir -p /go/src/github.com/ashukhotski/secret-santa-service
WORKDIR /go/src/github.com/ashukhotski/secret-santa-service

COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/secret-santa-service

FROM scratch
COPY --from=build /go/bin/secret-santa-service /go/bin/secret-santa-service
EXPOSE 8080 8080
ENTRYPOINT ["/go/bin/secret-santa-service"]