# --- Build Stage ---
    FROM golang:1.26.8-alpine3.24 AS build

    WORKDIR /app
    
    # Nur go.mod | go.sum zuerst → besseres Caching
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Danach den Rest
    COPY . .
    
    # Statisches Linux-Binary bauen
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dummycontentsigner ./main.go
    
    # --- Runtime Stage ---
    FROM alpine:3.20 AS runner
    
    WORKDIR /opt
    
    # Minimal notwendige Runtime-Pakete
    RUN apk --no-cache add ca-certificates
    
    COPY --from=build /app/dummycontentsigner /opt/dummycontentsigner
    
    EXPOSE 8080
    
    CMD ["./dummycontentsigner"]
    
