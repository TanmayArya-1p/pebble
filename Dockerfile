FROM golang:latest
ENV PORT=5000
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 5000
CMD ["/app/main"]
