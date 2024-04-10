FROM golang:1.22-alpine3.19 AS builder
ADD . /opt
COPY ./app/cmd/app/config.local.yaml /opt
WORKDIR /opt
#COPY wait-for-it.sh /wait-for-it.sh
#RUN chmod +x /wait-for-it.sh
#CMD ["./wait-for-it.sh", "database:5432"]
RUN  go mod download && go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -o /backend ./app/cmd/app/main.go

FROM alpine:3.18.0
WORKDIR /
COPY --from=builder /backend /backend
COPY --from=builder opt/config.local.yaml /config.local.yaml

#CMD ["./wait-for-it.sh", "database:5432"]
CMD ["./backend"]