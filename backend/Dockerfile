# Build stage
FROM golang:1.21-alpine AS build
WORKDIR /src

COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download

COPY backend/ ./backend/
COPY frontend/ ./frontend/

RUN cd backend && CGO_ENABLED=0 GOOS=linux go build -o /out/server .

# Runtime stage
FROM alpine:3.20
RUN apk add --no-cache ca-certificates

WORKDIR /app/backend

COPY --from=build /out/server /app/backend/server
COPY --from=build /src/backend/templates /app/backend/templates
COPY --from=build /src/backend/static /app/backend/static
COPY --from=build /src/frontend /app/frontend

EXPOSE 8080
CMD ["/app/backend/server"]
