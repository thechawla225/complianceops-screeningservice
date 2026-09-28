FROM golang:1.26 AS build
WORKDIR /app
COPY . .
ENV GOPROXY=https://goproxy.io,direct
RUN CGO_ENABLED=0 go build -o screening .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /app/screening /screening
USER nonroot:nonroot
ENTRYPOINT ["/screening"]