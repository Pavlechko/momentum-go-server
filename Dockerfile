FROM golang:1.21.1-bullseye
# RUN go mod download
# RUN mkdir -p /app
WORKDIR /usr/src/app
# COPY go.mod .
# COPY go.sum .
# COPY /cmd/momentum/. .
# COPY /internal/. .
COPY . .
# RUN go get -d -v ./...
# RUN go install -v ./...
RUN go build -o bin ./cmd/momentum


EXPOSE 8080

# ENTRYPOINT [ "/go/cmd/momentum/app/bin" ]
CMD [ "./bin" ]