FROM scratch
ADD bin/main /
ENTRYPOINT ["/main"]
