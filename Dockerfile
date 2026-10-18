FROM golang:1.26 AS build
WORKDIR /app
COPY . .
ENV GOPROXY=https://goproxy.io,direct
RUN CGO_ENABLED=0 go build -o screening .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /app/screening /screening
# main.go seeds Redis from this file at startup and exits if it's missing
COPY --from=build /app/internal/watchlist.json /watchlist.json
ENV WATCHLIST_PATH=/watchlist.json
USER nonroot:nonroot
ENTRYPOINT ["/screening"]
