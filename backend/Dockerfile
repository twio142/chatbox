FROM golang:1.22-alpine3.20
RUN apk add bash ca-certificates git gcc g++ libc-dev
RUN mkdir /app
ADD . /app/
WORKDIR /app
ARG APP_ENV=development
ENV APP_ENV=$APP_ENV
ARG PORT=8080
ENV PORT=$PORT
EXPOSE $PORT
RUN go mod download
RUN go build -v
CMD ["/app/realtime-chat-go-react"]
