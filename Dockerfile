FROM golang:1.21.1-alpine AS build

# create a directory named `app` inside the container
# and also tells docker to use this directory as a default destination for all the subsequent commands
WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o auth


FROM alpine:latest AS prod

WORKDIR /app
COPY --from=build /app/auth .

EXPOSE 1323
CMD ["./auth"]
