FROM oryd/kratos:v25.4.0

COPY config/kratos/ /etc/config/kratos/

CMD ["serve", "-c", "/etc/config/kratos/kratos.prod.yml", "--watch-courier"]
