FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

ADD server ./server
ADD config ./config
ADD sessions ./sessions
ADD telemetry ./telemetry
ADD users ./users

COPY *.go ./

RUN go build -o /auth_api

EXPOSE 8840

CMD [ "/auth_api" ]