FROM debian:bullseye-slim

WORKDIR /app

COPY phone /app/
COPY web /app/web

RUN chmod +x /app/phone
EXPOSE4 443

CMD ["/app/phone"]