FROM golang:1.16.6-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY main.go ./
RUN go mod download && go build -o reaper
CMD /app/reaper
