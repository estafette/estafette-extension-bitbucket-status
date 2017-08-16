FROM scratch

MAINTAINER estafette.io

COPY ca-certificates.crt /etc/ssl/certs/
COPY estafette-extension-bitbucket-status /

ENTRYPOINT ["/estafette-extension-bitbucket-status"]