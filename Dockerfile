FROM golang:latest as build
RUN mkdir /app
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -v ./cmd/vk-bot

FROM alpine
COPY --from=build /app/vk-bot /app/
WORKDIR /app
EXPOSE ${CALLBACK_PORT}
CMD ["/app/vk-bot"]