FROM golang:1.11.1-alpine3.8
RUN mkdir /app
ADD . /app/
WORKDIR /app
ARG APP_ENV=development
ENV APP_ENV $APP_ENV
ENV PORT $PORT
RUN go mod download
RUN go build -o main ./...
CMD ["/app/main"]
