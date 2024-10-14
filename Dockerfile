FROM golang:1.22-bullseye

WORKDIR /wd

COPY go.mod .
COPY go.sum .

RUN go mod download -x

COPY . .

CMD ["go", "run", "main.go"]
