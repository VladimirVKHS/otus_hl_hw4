# Мониторинг (homework 12)

## Zabbix

- zabbix агент установлен на оба инстанса сервиса диалогов (внедрен в докер образ сервиса диалогов)
- в docker-compose добавлены:
 
     - zabbix-db   - сервер БД для Zabbix
     - zabbix-server-mysql   - сервер Zabbix
     - zabbix-server-mysql   - сервер Zabbix
     - zabbix-web-nginx-mysql - web-интерфйес Zabbix
  
- в web-интерфйесе (http://localhost:9009/) добавлены хосты обоих инстансов для мониторинга за состоянием серверов

## Prometheus

- в сервис диалогов добавлен prometheus-экспортер (эндпоинт "/metrics") c RED-метриками работы сервиса:

     - dialog_requests_count: Counter - число запросов
     - dialog_errors_count: Counter - число ошибок
     - dialog_request_time: Gauge - время выполнения запроса, мс

- в docker-compose добавлены:

     - prometheus - сервер Prometheus. Конфигурация: [prometheus.yml](./docker/prometheus/prometheus.yml)
  
- в web-интерфейсе (http://localhost:9090/) протестированы запросы на получение данных:
  
    - число запросов: rate(dialog_requests_count[1m])
    - число ошибок: rate(dialog_errors_count[1m])
    - время запроса: avg_over_time(dialog_request_time[1m])

## Grafana

- в docker-compose добавлены:

    - grafana - сервер Grafana.

- в web-интерфейсе (http://localhost:9010/):

    - установлен плагин Zabbix
    - добавлены источники данных Prometheus и Zabbix
    - добавлен дашборд со следующими данными:
        
        - число запросов, 1/с
        - число ошибок, 1/с
        - среднее время запроса, мс
        - CPU utilization, %
        - Available memory, %
        - Space utilization, %
        - Список проблем с серверами сервиса диалогов
      

[Скриншоты](https://docs.google.com/document/d/1tW-rdikGJgokTzVcz8e6744tEpoOWR5HGkzq763ldDo/edit?usp=sharing)