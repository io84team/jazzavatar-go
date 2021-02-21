FROM golang:alpine as go
WORKDIR /src
COPY ./ /src
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

###
FROM alpine:latest
WORKDIR /app
COPY --from=go /src/ /app

CMD [ "/app/main" ]