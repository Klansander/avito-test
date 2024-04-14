FROM golang:1.22-alpine3.19 AS builder
ADD . /opt
COPY ./app/cmd/app/config.local.yaml /opt
WORKDIR /opt

RUN  go mod download && go mod verify
RUN GOOS=linux go build -o /backend ./app/cmd/app/main.go


FROM alpine:3.18.0
WORKDIR /
COPY --from=builder /backend /backend
COPY --from=builder opt/config.local.yaml /config.local.yaml

CMD ["./backend"]