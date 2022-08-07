FROM golang:alpine AS development
WORKDIR /app/goFileService
COPY . /app/goFileService
ENV GO111MODULE=auto
ENV GOPROXY https://goproxy.cn,direct
RUN go version
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /app/goFileService/cmd/fileService/service /app/goFileService/cmd/fileService/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /app/goFileService/cmd/migrate/dbcli /app/goFileService/cmd/migrate/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /app/goFileService/cmd/fscli /app/goFileService/cmd/main.go

FROM alpine:latest AS production
WORKDIR /app
COPY --from=development /app/goFileService/cmd/fileService/service .
COPY --from=development /app/goFileService/cmd/migrate/dbcli .
COPY --from=development /app/goFileService/cmd/fscli .
CMD ["/app/service"]
