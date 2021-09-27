FROM alpine:latest
WORKDIR /app
COPY reaper ./
RUN addgroup -S reaper && adduser -S reaper -G reaper
USER reaper
ENTRYPOINT ["./reaper"]
