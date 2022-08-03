FROM golang:alpine AS development
WORKDIR /app
COPY . /app
ENV GOPROXY https://goproxy.cn
RUN go version
WORKDIR /app/goFileService
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o cmd/fileService/service cmd/fileService/amin.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o cmd/migrate/dbcli cmd/migrate/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o cmd/fscli cmd/main.go

FROM alpine:latest AS production
WORKDIR /app
COPY --from=development /app/goFileService/cmd/fileService/service .
COPY --from=development /app/goFileService/cmd/migrate/dbcli .
COPY --from=development /app/goFileService/cmd/fscli .
CMD ["/app/goFileService/service"]
