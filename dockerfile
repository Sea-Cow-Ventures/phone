FROM debian:bullseye-slim

# Add bullseye-backports repository to get newer versions of packages like libc6
RUN echo "deb http://deb.debian.org/debian bullseye-backports main" > /etc/apt/sources.list.d/bullseye-backports.list

# Install ca-certificates, libc6 (from backports), and update
RUN apt-get update && \
    apt-get install -y ca-certificates libc6/experimental && \
    rm -rf /var/lib/apt/lists/*
	
WORKDIR /app

COPY phone /app/
COPY web /app/web

RUN chmod +x /app/phone
EXPOSE 443

CMD ["/app/phone"]