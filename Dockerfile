FROM migrate/migrate as migrate

WORKDIR /migrate
COPY ./migraitons .

ARG HOST
ARG PORT
ARG USER
ARG PASSWORD
ARG PATH


RUN migrate 


# FROM golang:1.22.1

# WORKDIR /app

# COPY . .
# RUN go mod download
# RUN go mod verify

# RUN GOOS=linux go build -o bin ./cmd/main.go