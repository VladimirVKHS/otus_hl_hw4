FROM golang:latest as builder
RUN mkdir /app
ADD . /app/
WORKDIR /app
EXPOSE 7000
RUN GOOS=linux GOARCH=amd64 go build -o /app/main /app/cmd/app/main.go

FROM ubuntu:18.04

ADD ./docker/app/zabbix_agentd /usr/sbin/zabbix_agentd
ADD ./docker/app/zabbix_agentd.conf /etc/zabbix/zabbix_agentd.conf
COPY --from=builder /app/main /usr/bin/main
COPY --from=builder /app/logs /usr/bin/logs
COPY --from=builder /app/docker/app/env /usr/bin/.env
WORKDIR /usr/bin/
EXPOSE 7000
EXPOSE 10050
ADD ./docker/app/entrypoint.sh /bin/entrypoint.sh
CMD /bin/entrypoint.sh