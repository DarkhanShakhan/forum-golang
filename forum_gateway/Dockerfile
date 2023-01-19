FROM golang:alpine3.16 AS build
LABEL stage=build
WORKDIR /app
COPY . ./
RUN apk add build-base
RUN go build cmd/main.go
#copy all needed files into second container
FROM alpine:3.16 AS runner
WORKDIR /app
LABEL authors="@Subudei, @DarkhanShakhan"
COPY --from=build /app/main /app/main
COPY /templates /app/templates
COPY /web /app/web
CMD ["/app/main"]
