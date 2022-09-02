FROM golang:1.16-alpine AS builder
ADD ./ /src
WORKDIR /src
RUN go build .

FROM alpine:latest AS runner
WORKDIR /bin
COPY --from=builder /src/EPGSWatcher /bin/

ENTRYPOINT ["/bin/EPGSWatcher"]