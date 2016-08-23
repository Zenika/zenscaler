FROM docker/compose:1.8.0
MAINTAINER Maximilien Richer <maximilien.richer@gmail.com>
COPY ./build/zscaler /app/zscaler
EXPOSE 3000
WORKDIR /app/config
ENTRYPOINT ["/app/zscaler"]

