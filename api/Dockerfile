FROM golang:1.26 AS dev

RUN go install github.com/air-verse/air@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/riverqueue/river/cmd/river@latest

COPY ./.local/docker/watch.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENV PATH="$PATH:/app/bin"

# This is where you should mount your code
VOLUME /app

ENTRYPOINT ["/entrypoint.sh"]
