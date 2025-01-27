FROM debian:bullseye-slim

WORKDIR /app

COPY seacow-phone /app/

RUN chmod +x /app/seacow-phone

EXPOSE 80

CMD ["/app/seacow-phone"]