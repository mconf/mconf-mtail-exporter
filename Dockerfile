FROM golang

COPY ./main.go ./exporter/main.go

WORKDIR ./exporter

RUN go mod init mtailexporter && \
    go get k8s.io/client-go@v0.17.0 && \
    go mod tidy

RUN go build -o exporter main.go

CMD ["./exporter"]