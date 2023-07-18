FROM golang:1.19-bookworm

ENV APP_HOME /opt/docker/scripts
WORKDIR $APP_HOME

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./snowflake-was

ENTRYPOINT ["./snowflake-was"]
CMD ["was"]
