FROM scratch
COPY ./zscaler /
EXPOSE 3000
ENTRYPOINT ["/zscaler"]

