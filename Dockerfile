FROM debian:latest

RUN apt-get update
RUN apt-get install -y libltdl-dev && rm -rf /var/lib/apt/lists/*

RUN mkdir -p /var/ocms
RUN mkdir -p /var/log/ocms
COPY ocms /var/ocms/ocms
COPY ocms_prod.toml /var/ocms/ocms.toml
ADD ./fixtures /var/ocms/fixtures

# Set binary as entrypoint
ENTRYPOINT cd /var/ocms && ./ocms

# Expose port (8000)
EXPOSE 8000