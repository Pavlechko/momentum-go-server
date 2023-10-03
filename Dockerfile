FROM golang:1.21.1-bullseye
# RUN go mod download
# RUN mkdir -p /app
WORKDIR /go/cmd/momentum/app
COPY . .
COPY /cmd/momentum/. .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o bin .

EXPOSE 8080

ENTRYPOINT [ "/go/cmd/momentum/app/bin" ]
