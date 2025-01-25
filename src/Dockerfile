FROM golang:1.23-alpine AS gobuild

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build
ADD . /build

RUN go get -d -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w -extldflags "-static"' -o ./prometheus ./cmd/prometheus/main.go

RUN chmod +x ./prometheus

FROM alpine:3.20.3

ARG TARGETOS
ARG TARGETARCH
# Install required packages and manually update the local certificates
RUN apk add --no-cache tzdata bash ca-certificates && update-ca-certificates
# Copy executable from build
COPY --from=gobuild /build/prometheus /prometheus
# Expose port 8080 and 9001 by default
EXPOSE 8080
EXPOSE 9001
# Set entrypoint to run executable
ENTRYPOINT [ "/prometheus" ]
CMD [ "agent" ]