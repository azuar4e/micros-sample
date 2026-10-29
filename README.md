# Micros-Sample (gRPC branch)

![Go](https://img.shields.io/badge/Go-1.24-blue)
![Kubernetes](https://img.shields.io/badge/Kubernetes-1.34-326CE5?logo=kubernetes&logoColor=white)
![Cilium](https://img.shields.io/badge/Cilium-1.20-F8C517?logo=cilium&logoColor=black)
![gRPC](https://img.shields.io/badge/gRPC-Protocol%20Buffers-244c5a?logo=grpc&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=white)

This repository contains the source code for a sample mini e-commerce application. The API Gateway is the entry point to the application, while the Orders and Products services implement the business logic.

This is the `grpc` branch. It is the version used in my article [Cilium + EBPF: Observability and Security in Kubernetes](https://azuar4e.github.io/en/posts/cilium/), where the microservices run on a Cilium cluster. The original version, with Istio and Keycloak, is explained in [Zero Trust in Kubernetes](https://azuar4e.github.io/en/posts/micros-sec-istio/).

## Overview

The application consists of three Go microservices:

- **API Gateway**: Entry point of the application. Forwards the requests over REST to the appropriate backend service.
- **Products Service**: Exposes a simple in-memory product catalog, over both REST (for the API Gateway) and gRPC (for the Orders service).
- **Orders Service**: Provides endpoints to create and retrieve orders. When creating an order, it validates the product by calling the Products service over gRPC.

The repository includes the Kubernetes manifests required to deploy the application components, the Gateway API resources and the Cilium network policies.

## Differences with the original branch

Compared to the original version of the project, this branch changes the following:

- **Products and Orders now talk gRPC.** The contract is defined in `proto/product.proto` (`ProductService` with `GetProducts` and `GetProductById`). The Products service serves gRPC on port `8077` and keeps REST on `8081`, since the API Gateway still uses REST. Orders is the gRPC client.
- **New `proto/` module** (`github.com/azuar4e/proto`), shared by both services. Each service references it with a `replace` directive in its `go.mod`, so the Dockerfiles of Products and Orders were updated and their images must be built from the root of the repository.
- **Authentication was removed.** There is no Keycloak, PostgreSQL or OAuth 2.0 flow in this branch, to keep the focus on the protocols and the network.
- **No Istio.** Cilium replaces it as CNI. Traffic between nodes is encrypted with WireGuard.
- **HTTPS ingress with the Kubernetes Gateway API** (`shop-gw`, host `shop.local`) instead of an Istio Gateway.
- **Network policies are `CiliumNetworkPolicy`** with L7 rules. For example, Products only accepts `GetProductById` over gRPC from Orders, and `GET /products` and `GET /products/{id}` over REST from the API Gateway.

## Project structure

```text
.
├── api/              # API Gateway
├── order-service/    # Orders service (gRPC client)
├── product-service/  # Products service (REST + gRPC server)
├── proto/            # Protocol Buffers contract and its Go module
└── manifests/        # Deployments, Services, Gateway, network policies
```

## Requirements

- A Kubernetes cluster with Cilium as CNI and Gateway API enabled (tested with kind, Kubernetes 1.34 and Cilium 1.20). The article explains how to set it up.
- Docker
- kubectl
- OpenSSL (to generate the self-signed certificate)
- Go 1.24+ and `protoc` with the `protoc-gen-go` and `protoc-gen-go-grpc` plugins, only if you want to modify the `.proto` file and regenerate the code.

## Deployment

To deploy the microservices locally, first build their Docker images and then load them into the kind cluster. Run these commands from the root of the repository, since the Products and Orders images need the `proto/` folder as part of the build context:

```bash
docker build -t azuar4e/apigw-ecom:v2 .
docker build -f ./product-service/Dockerfile -t azuar4e/products-svc:v3 .
docker build -f ./order-service/Dockerfile -t azuar4e/orders-svc:v3 .
```

> The API Gateway image is built from its own directory (`api/`), as it does not use the proto code.

Then load them into the cluster:

```bash
kind load docker-image azuar4e/products-svc:v3 azuar4e/orders-svc:v3 azuar4e/apigw-ecom:v2 --name kind
```

> Notice that I'm using kind for my local cluster. If you are using another one, like Minikube, look for the proper command.

To deploy the manifests:

1. Change the images specified in the Deployments to match the images you have built.
2. Create a self-signed certificate and a TLS secret for the Gateway (`shop-tls`, referenced by `shop-gw`):

```bash
openssl req -x509 -newkey rsa:2048 -nodes -days 365 \
  -keyout tls.key -out tls.crt \
  -subj "/CN=shop.local" -addext "subjectAltName=DNS:shop.local"

kubectl create secret tls shop-tls --cert=tls.crt --key=tls.key
```

3. Apply the manifests:

```bash
kubectl apply -f manifests/
```
