FROM golang:1.23-alpine

WORKDIR /app

RUN apk add --no-cache \
    git \
    bash \
    tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV PORT=8080
ENV TZ=Asia/Jakarta

EXPOSE 8080

CMD ["go", "run", "main.go"]
