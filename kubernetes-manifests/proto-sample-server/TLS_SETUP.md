# Setup TLS

* [Private Certificate Authority](#private-certificate-authority)
* [Public Certificate Authority](#public-certificate-authority)

The instructions here provide two examples to setup TLS for the NGINX ingress: private and public certificate authorities. Both options require the [Prerequisites](#prerequisites) below

## Prerequisites

Enable the minikube Docker for Mac port-forward for ingress controller (in a separate window)

    minikube service ingress-nginx-controller -n ingress-nginx

# Private Certificate Authority

Create the RSA public/private keypairs with SAN

     openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=grpctest.dev.mydomain.com" -addext "subjectAltName = DNS:grpctest.dev.mydomain.com"


Create the TLS secret in kubernetes

    kubectl create secret tls wildcard.dev.mydomain.com --key tls.key --cert tls.crt


Test `grpcurl` with `-authority` header (if DNS not configured) and CA root with not RSA public key

    grpcurl -vv -cacert tls.crt -authority=grpctest.dev.mydomain.com 127.0.0.1:49537 list


Test `grpcurl` with `-authority` header (if DNS not configured)

    grpcurl -vv -insecure -authority=grpctest.dev.mydomain.com 127.0.0.1:49537 list


# Public Certificate Authority

TO-DO
