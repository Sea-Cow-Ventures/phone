FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates libc6 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY phone /app/
COPY web /app/web

RUN chmod +x /app/phone
EXPOSE 443

CMD ["/app/phone"]