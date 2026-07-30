# Micros-Sample

![Go](https://img.shields.io/badge/Go-1.24-blue)
![Kubernetes](https://img.shields.io/badge/Kubernetes-1.33-326CE5?logo=kubernetes&logoColor=white)
![Istio](https://img.shields.io/badge/Istio-1.30.2-466BB0?logo=istio&logoColor=white)
![Keycloak](https://img.shields.io/badge/Keycloak-26-4D4D4D?logo=keycloak&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=white)

This repository contains the source code for a sample mini e-commerce application. The API Gateway is the entry point to the application, while the Orders and Products services implement the business logic. 

The complete explanation of the project can be found in my [Zero Trust in Kubernetes](https://azuar4e.github.io/en/posts/micros-sec-istio/) article on my blog.

## Overview

This project implements a simple microservices-based e-commerce application designed to demonstrate modern authentication and secure service-to-service communication in Kubernetes.

The application consists of three Go microservices:

- **API Gateway**: Entry point of the application. Handles the OAuth 2.0 Authorization Code + PKCE flow with Keycloak and forwards requests to the appropriate backend service.
- **Products Service**: Exposes a simple in-memory product catalog.
- **Orders Service**: Provides endpoints to create and retrieve orders.

The repository includes the Kubernetes manifests required to deploy the application components (API Gateway, services, secrets, Keycloak and PostgreSQL).

The accompanying blog post explains how Istio can be integrated with the application to provide:

- OAuth 2.0 / OpenID Connect authentication with Keycloak
- mTLS between microservices
- HTTPS ingress through an Istio Gateway
- AuthorizationPolicy and PeerAuthentication resources

The project is intended to run entirely on a local Kubernetes cluster (Kind), allowing the complete authentication and networking stack to be tested without relying on a cloud provider.

## Deployment

To deploy the microservices locally, first build their Docker images and then load them into the Kind cluster.

For example, for the API Gateway we will do.
````Bash
docker build -t azuar4e/apigw-ecom:v14 .
````

And then:

````Bash
kind load docker-image azuar4e/products-svc:local azuar4e/orders-svc:v2 azuar4e/apigw-ecom:v14 --name kind
````

>Notice that i'm using Kind for my local cluster. If you are using other type like Minikube, look for the proper command.

To deploy the infrastructure manifests:

1. Change the images specified in the Deployments to match the images you have built.
2. As explained in the blog post, create a self-signed certificate and private key, and create a TLS secret for the Istio Gateway (`credentialName: api-gateway-tls-secret`).

To create those files you can use the following command:

```bash
openssl req -x509 -newkey rsa:2048 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=localhost"
```

And then:

```bash
kubectl create secret tls api-gateway-tls-secret --key=key.pem --cert=cert.pem -n istio-system
```

## Requirements

- Go 1.24+
- Docker
- Kind
- kubectl
- Istio (`istioctl`)
- OpenSSL (to generate self-signed certificates)
