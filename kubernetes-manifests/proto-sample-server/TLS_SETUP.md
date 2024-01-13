# Setup TLS

* [Private Certificate Authority](#private-certificate-authority)
* [Public Certificate Authority](#public-certificate-authority)
* [Troubleshooting](#troubleshooting)

The instructions here provide two examples to setup TLS for the NGINX ingress: private and public certificate authorities. Both options require the [Prerequisites](#prerequisites) below

## Prerequisites

Enable the minikube Docker for Mac port-forward for ingress controller (in a separate window). Capture the second IP:port from the URLS that are exposed from the minikube reverse-proxy (e.g http://127.0.0.1:56836 --> 127.0.0.1:56836). Since there are two (http:80 and https:443), we will use the second one for TLS.

    minikube service ingress-nginx-controller -n ingress-nginx


# Private Certificate Authority

Create the RSA public/private keypairs with SAN

     openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=*.dev.abc.com" -addext "subjectAltName = DNS:proto-sample.dev.abc.com"


Create the TLS secret in kubernetes

    kubectl create secret tls dev.abc.com --key tls.key --cert tls.crt


Test `grpcurl` with `-authority` header (if DNS not configured) and CA root with not RSA public key

    grpcurl -vv -cacert tls.crt -authority=proto-sample.dev.abc.com 127.0.0.1:49537 list


Test `grpcurl` with `-authority` header (if DNS not configured)

    grpcurl -vv -insecure -authority=proto-sample.dev.abc.com 127.0.0.1:49537 list


View the Ingress Controller logs

     kubectl logs -f deployment.apps/ingress-nginx-controller -n ingress-nginx

# Public Certificate Authority


# Troubleshooting

Here are some common scenarios your may encounter when configuring and testing

**Insecure Port**

Issue: Clients try to connect to Ingress and receive the following error(s)

    Failed to dial target host "127.0.0.1:57361": tls: first record does not look like a TLS handshake

Details: This typically occurs if Ingress is only configured for HTTPS but exposes both HTTP & HTTPS ports and does not auto-redirect HTTP -> HTTPS.

Resolve: Check the Ingress configuration and have clients change to the HTTPS port to resolve

**Incorrect Certificate or Client Host**

Issue: Clients try to connect to Ingress and receive the following error(s)

    Failed to dial target host "127.0.0.1:57362": tls: failed to verify certificate: x509: certificate is valid for ingress.local, not 123.com

Details: This typically occurs if Ingress is configured for a specific domain FQDN or wildcard FQDN and the client does not have the correct Host or X-Forwarded-Host domain in the header.

Resolve: Check the Ingress configuration and certificate SAN names and have clients verify and retry

**Private Certificate**

Issue: Clients try to connect to Ingress and receive the following error(s)

    Failed to dial target host "127.0.0.1:57362": tls: failed to verify certificate: x509: certificate signed by unknown authority

Details: This typically occurs when the Ingress certificate is configured using a private certificate authority (CA) and the client does not have the public certificate to verify.

Resolve: Check the Ingress configuration and certificate CA and have the clients verify their TLS configuration



# TO-DO
