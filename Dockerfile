# syntax=docker/dockerfile:1

FROM golang:1.18-alpine
# install golabg-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

ENV CGO_ENABLED=0
#ENV APP_ENV="dev"

# Set destination for COPY
WORKDIR /app

COPY internal/env/.env ./

COPY . ./


# Download Go modules
RUN go mod download

# Build
RUN go build -o main ./cmd/short-access/main.go

EXPOSE 8080

# Run
CMD ["./main"]