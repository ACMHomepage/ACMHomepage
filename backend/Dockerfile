FROM golang:1.19

WORKDIR /usr/src/app

# Download the modules.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Build.
COPY . .
RUN go build -v -o /usr/local/bin/app .

# Meta infomation.
EXPOSE 8080

# Run.
CMD ["app"]
