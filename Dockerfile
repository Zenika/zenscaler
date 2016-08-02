FROM scratch
COPY ./build/zscaler /
EXPOSE 3000
ENTRYPOINT ["/zscaler"]

