# BigBang Project

BigBang is a Go-based application consisting of two main servers:
1. **gRPC Server** - Distributes configurations to Envoy instances using go-control-plane.
2. **REST Server** - Provides API endpoints for the [Elchi](https://github.com/sefaphlvn/elchi), handling CRUD operations with MongoDB.

## Table of Contents
- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Running the Servers](#running-the-servers)
  - [gRPC Server](#grpc-server)
  - [REST Server](#rest-server)
- [Project Structure](#project-structure)
- [Environment Variables](#environment-variables)
- [Troubleshooting](#troubleshooting)
- [License](#license)

## Getting Started

These instructions will help you set up the BigBang project on your local machine for development and testing purposes.

### Prerequisites

Before you begin, ensure you have met the following requirements:
- Go 1.18 or later installed on your machine
- MongoDB running locally or accessible via network
- Docker (optional) if running MongoDB or Envoy instances locally

### Installation

Clone the repository to your local machine:

```bash
git clone https://github.com/sefaphlvn/bigbang.git
cd bigbang
```

Install the required Go modules:

```bash
go mod tidy
```

## Running the Servers

### gRPC Server

The gRPC server is responsible for distributing configurations to Envoy instances using go-control-plane. It serves as the control plane in an Envoy xDS setup.

To run the gRPC server:

```bash
go run main.go server-grpc
```

#### Configuration

The gRPC server uses a snapshot-based approach to distribute configurations. It supports Delta gRPC and utilizes `go-control-plane` for managing Envoy's xDS resources like CDS, EDS, LDS, and RDS.

**Important Note:** If you encounter initial fetch timeouts for EDS, you can enable the runtime flag `envoy.restart_features.use_eds_cache_for_ads` in your Envoy bootstrap configuration.

### REST Server

The REST server provides API endpoints for the Elchi frontend application, handling CRUD operations with MongoDB.

To run the REST server:

```bash
go run main.go server-rest
```

#### API Endpoints

The REST server provides a variety of endpoints for managing resources in MongoDB, including but not limited to:

- **CRUD operations**: Create, read, update, and delete resources like configurations, secrets, and certificates.
- **Resource Management**: Interact with MongoDB to manage configurations required by Elchi's frontend.

### Troubleshooting

- **EDS Initial Fetch Timeout**: If you see warnings about EDS initial fetch timeouts, ensure that the runtime flag \`envoy.restart_features.use_eds_cache_for_ads\` is enabled in your Envoy configuration. This flag is enabled by default starting from Envoy v1.31.0.
- **MongoDB Connection Issues**: If you encounter errors related to MongoDB connections, check that your \`MONGO_URI\` is correct and that MongoDB is accessible.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.