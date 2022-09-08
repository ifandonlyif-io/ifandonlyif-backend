# Build stage
FROM golang:1.18-alpine3.16 AS builder
ENV CGO_ENABLED 0
WORKDIR /app
COPY . .

RUN go build -o main main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY db/migration ./db/migration
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
RUN ["chmod", "+x", "/app/wait-for.sh"]
RUN ["chmod", "+x", "/app/start.sh"]


EXPOSE 8080
EXPOSE 5432
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]