FROM golang:1.16.6-alpine
WORKDIR /app
COPY ./ ./
RUN cd src && go mod download && go build -o reaper
CMD /app/reaper
