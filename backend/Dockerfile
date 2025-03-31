FROM golang:1.24-alpine3.20
RUN apk add --no-cache bash ca-certificates git gcc g++ libc-dev

RUN addgroup -S app && adduser -S app -G app && \
    mkdir -p /app && \
    chown -R app:app /app

WORKDIR /app
COPY --chown=app:app . .

USER app

ARG APP_ENV=development
ENV APP_ENV=$APP_ENV
ARG PORT=8080
ENV PORT=$PORT
EXPOSE $PORT

RUN go mod download && go build -v

CMD ["/app/realtime-chat-go-react"]
