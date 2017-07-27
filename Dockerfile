FROM scratch
ADD ca-certificates.crt /etc/ssl/certs/
ADD main /
EXPOSE 80
CMD ["/main"]