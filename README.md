
# **Wallet Service**

### 1. Overview

The **Wallet Service** is a microservice that manages wallet-related operations, including wallet creation, deposits, and withdrawals. It integrates with external payment gateways to process transactions.

Key capabilities include:
- Wallet creation and management.
- Handling deposits and withdrawals.
- Integration with multiple payment gateways.
- Asynchronous transaction processing to ensure responsiveness and fault tolerance.
- Extensible architecture, enabling easy integration with additional payment gateways.

### 2. Project Structure

The project is organized into various directories, each with a clear responsibility, making it easy to maintain and extend. Below is the directory structure:

```bash
├── cmd/                         # Command entry points for the services
│   └── api/
│       ├── mock/                # Mock payment gateway services
│       │   ├── gateway-a/       # Mock for Gateway A (JSON)
│       │   └── gateway-b/       # Mock for Gateway B (XML)
│       └── wallet/              # Main wallet service entry point
├── config/                      # Service configuration (env variables, etc.)
│   └── config.go
├── docker/                      # Docker setup for building and running services
│   ├── dockerfile.mocks
│   └── dockerfile.wallet
├── docker-compose.yml           # Docker Compose file to run all services
├── docs/                        # OpenAPI/Swagger documentation
│   ├── docs.go
│   └── swagger.yaml
├── internal/                    # Core business logic and internal packages
│   ├── db/                      # Database migration scripts
│   ├── dto/                     # Data Transfer Objects for requests/responses
│   ├── handlers/                # HTTP/gRpc layer handlers
│   ├── models/                  # Entities (e.g., wallet, transaction)
│   ├── payment/                 # Payment processing package for gateways
│   ├── repos/                   # Repositories for interacting with data source - in this service we have only postgres
│   ├── services/                # Business logic / Usecase services 
│   └── web/                     # Middleware, routing, and logging utilities
├── pkg/                         # Reusable packages and utilities
│   ├── database/                # Database connection helpers
│   ├── errs/                    # Custom error system handling package with custom system codes
│   ├── logger/                  # Logging package
│   ├── queue/                   # Simple in-memory queue for async processing
│   ├── rest/                    # HTTP client with retry functionality
│   └── web/                     # Web response helpers
```
### 3. Features

- **Wallet Management**: Create and list wallets.
- **Transaction Handling**: Supports deposits and withdrawals using external payment gateways.
- **Async Processing**: Transactions are processed asynchronously via an queue and background worder for improved performance and non-blocking execution.
- **Transaction Tracking**: Transactions can be tracked via API, allowing the status of transaction to be monitored.
- **Mock Payment Gateways**: Includes two mock payment gateways (A for JSON, B for XML) to simulate payment flows.
- **Extensibility**: Easily add new payment gateways by implementing the `PaymentGateway` interface.


### 4. Architecture Overview

The **Wallet Service** is built with a modular and extensible architecture to handle wallet and transaction operations and integrate with external payment gateways. The design separates business logic from payment processing, making the system easy to extend and scale.

#### **Core Components**:

- **Wallet Service**:
  - Manages wallet creation, deposits, and withdrawals.
  - Handles asynchronous transaction processing using queue with background worder.
  - Stores and updates transactions in the database.

- **Payment Package**:
  - Provides a unified interface for interacting with various payment gateways.
  - Includes functions for deposit, withdrawal, and callback handling.
  - Supports retry and circuit breaker mechanisms for gateway communication.

The architecture allows for easy integration of additional payment gateways and ensures that the system is scalable and robust.

### 5. Technical Design and Decisions

- **Asynchronous Processing**: Keeps the API responsive by processing transactions in the background using an in-memory queue, ensuring non-blocking operations.
  
- **In-Memory Queue**: A simple and efficient mechanism for managing transaction processing within the assessment’s scope. Service Queue workers handle transactions asynchronously to offload processing.

- **Payment Gateway Abstraction**: The `PaymentGateway` interface enables easy integration of new payment gateways without modifying core business logic, ensuring high extensibility.

- **Retry and Circuit Breaker**: Implemented to ensure robustness when interacting with external gateways, preventing gateway failures from impacting the system. Both are customizable for each payment gateway client, allowing fine-tuned control over retry logic and circuit-breaking behavior depending on the gateway's characteristics and requirements.

- **Error Handling**: A custom `errs` package is used to create system-level error categorization. This approach will help in long-term system scalability by allowing easy error categorization and handling.

- **Logging**: A logging middleware is used to track incoming and outgoing requests. Due to time constraints, logs are focused on request/response cycles and within the transaction worker to track processing status.

- **Database Models**: PostgreSQL is used to store wallet and transaction data. The models are designed to be simple yet extensible, ensuring that wallet and transaction details are reliably stored.
  
- **Extensibility**: The system allows for easy integration of new gateways by implementing the `PaymentGateway` interface, enabling seamless extension of functionality.


### 6. Flow of Deposit/Withdrawal

The deposit and withdrawal flow is designed to handle transactions efficiently and asynchronously while maintaining a quick response time for clients.

#### **1. Client Request**:
The flow begins when the client submits a deposit or withdrawal request via the API. This request contains the amount and payment details.

#### **2. Validation and Transaction Creation**:
Once the request is received, the wallet service validates the amount and payment details (via payment package). If everything is valid:
- A new transaction is created in the database with a status of `created`.
- Sensitive payment details, such as credit card information, are stored in a masked format in the database.

#### **3. Immediate Response**:
After creating the transaction, the system pushes the transaction ID and the unmasked payment details into the queue for further processing. At this point, the API returns an immediate response to the client with the transaction ID and a status of `created`, allowing the client to track the progress of the transaction and also do action in case of payment require action.

#### **4. Asynchronous Processing**:
- The transaction is processed asynchronously in the background by a worker. The worker retrieves transaction ID and unmasked payment details and triggers the appropriate payment gateway interaction (either `Deposit` or `Withdraw`).
- The payment gateway (real or mock) processes the request and returns a status (`success`, `failed`, or `pending`).
- update the transaction status to `Pending`
  
#### **5. Callback Handling**:
- For asynchronous gateways, a callback URL is provided.
- When triggered, the callback is passed to the wallet service.
- The wallet service uses the payment package's `VerifyCallback` to select the appropriate gateway and verify the request.
- The transaction status is updated based on the gateway’s response to `completed`, `failed`.

Here’s the flow diagram representing the process:

![Flow Diagram](https://www.planttext.com/api/plantuml/png/ZPH1Zzem48Nl_XKZUjazLBtdK6s0BKBgmPHsSpRE0DOcSTQUGFFlEquS892LwcFUyyppDtPUF2b7JLa8eJHP1ul2O4MWF2p4plw5MLhNXT6AZArWhlGxLlaClhnsIm2lcWiORMh5ssQfNCDFrQARXHBfeo5JHO44MtGdex5pPTj7crHj6N98xgWElK_AHz-cmQPNDu_YKf7QAT_hoxdWwC1d4fETLehmhDg-qqg81Npz3ca2ssPN6e8SQ-iDVJiREkPEdS7XHuEUH1fysJQ17zQTbSilGhOTb3TLc9pBWofj4-1oZZgsBP6EDe_cvJo1XSDW9QSgpoC9s9zuIDJu17IdvS_Hlab0DluuyfA5Zy0aMXO9_49gN3KohPSJDSLco9jPzuuEQcSrUWzq7CM9bQNah3nCM4OoMIGZfEpqLGBhYj3nBWZKu8wqaAkXJepOHuhxGv3uVM3bq3S5tR3wK-VthFeQ0OFaSPlgm1Ux84XzM-akxuvlL7TL-lOyuI7NeSy5lvqv7FZy-jRzwPY3U4Va3PtPj-Dc5oQzEChSqHcwFvaz60-LfUuiMF08dcy2t_0wXLB3sunmhirk04uPcOxu7vINXy1FpRLJXYiQ97sSSbpRS21dy8HZ9RqatPjA5PyCL7U_9Y5UE3h_iVu1)

This flow ensures that transactions are processed efficiently in the background while the client gets a quick response. It is designed to handle high traffic, ensuring scalability and performance.


### 7. Extensibility:
Adding a new payment gateway is seamless. Simply implement the `PaymentGateway` interface and plug it with its key into the `Payment` package in cmd:

- Create client for the new payment gateway and implement the `PaymentGateway`:

```go
type PaymentGateway interface {
    Deposit(ctx context.Context, req *Request) (*Response, error)
    Withdraw(ctx context.Context, req *Request) (*Response, error)
    VerifyCallback(ctx context.Context, refID string, data []byte) (*Response, error)
    VerifyMethod(typ models.TransactionType, method models.PaymentMethod) error
}
```
-  update Payment gateway setup in cmd/api/wallet/main.go like:
```
// Payment gateway setup
paymentGateways := map[models.PaymentGateway]payment.PaymentGateway{
	models.PaymentGatewayA: gatewaya.New(cfg.PaymentGatewayConfig.GatewayA),
	models.PaymentGatewayB: gatewayb.New(cfg.PaymentGatewayConfig.GatewayB),
	models.NewGateway: newgateway.New(cfg.PaymentGatewayConfig.NewGateway)
}

// Payment handler setup
paymentHandler := payment.New(paymentGateways)
```
The system can support multiple payment gateways without modifying core business logic

### 8. Entities:
Here’s the Entities diagram representing the tables:
![Flow Diagram](https://www.planttext.com/api/plantuml/png/fPDDQuD048Rl_eh5a_qmD850yL1YgaqjQH8IGNgImPqa4dSZkZQ4ql_UrLZSc8yUkWV1TnxddPaT1xc0J1GiqRGKeWsiaEWk5x68CTV9bqRagHvH0dbE0aWI5BLUdhkOMgGeOjeeKOOWa8OWB29YXjA1fKsuIEc5yBUcEFaPy1mY4M_vTRjTLLBu6o36SdFJH85j2owTA4OnWyJeFjwJdX8N-nGjWhnWXcWSmr9MA5cZAF8pt26Wa2di6N8HhcIFEzZNdxJKCpn3iTxIaAA0E95ERulfP7W9iyWdPD4QCgFNxol9CbnYXZp2QXhdcV_UJjaFQOzAlI77dKqNdjy8WUU_EdCxiCTNynn6gMPwdhksxpgDC7CdZZSPASJqVJPsZvWNsnlNxwfJwmPKcv4q2UoFq3wLXcgUUlVrhavCa-WFdSwjVhIc5bb3Ng6gQffFf_EIkvhZtsmzaojqkwyQbIKFaDFol_u1)

---
## **How to Run the Project**

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-repo/wallet-service.git
   cd wallet-service
   ```

2. **Run the Project**:
   Using Docker Compose:
   ```bash
   docker-compose up --build
   ```

   The wallet service will be available at `http://localhost:8080/api/v1`.

3. **API Documentation**:
   Mock services both provide `/deposit` and `/withdraw` endpoints:
   - GatewayA `JSON` running on `http://localhost:8090`.
   - GatewayA `SOAP/XML` running on `http://localhost:8091`.
  
5. **Mock Gateways**
   - Access the OpenAPI (Swagger) documentation at `http://localhost:8080/swagger`.
   - Postman collection `wallet.postman_collection.json`

## **Configuration**

The service configuration can be adjusted via the `config/config.go` and `.env` file. these include setting up the database, payment gateways, and other environment variables.

---

## **Tests**

- Tests are provided **only for the payment package** due to time constraints.
- Tests follow the table-driven test pattern for flexibility and maintainability.
- Run tests with:
   ```bash
   go test ./...
   ```
---
## **Technologies Used**

- **Go**: Version 1.22
- **Postgres**: Database for wallet and transaction data.
- **Docker**: Containerization of services.
- **Docker Compose**: Manage multi-container applications.
- **dbmate**: Database migration tool.
- **Swagger (OpenAPI)**: API documentation.
- **Postman**: API testing and collection.
---
## Improvements and Future Work

Due to time constraints, the following improvements were not implemented but are planned:

- **Standalone Payment Service**: Separate the payment package into a standalone service focused exclusively on handling interactions with payment gateways.

- **Event-Driven Architecture**: Implement an event-driven architecture using Kafka or similar systems for asynchronous communication between services like wallet and payment in future.

- **Comprehensive Logging**: Implement full logging coverage using the ELK stack (Fluent, Logstash, Kibana) to better monitor and query logs.

- **Monitoring and Distributed Tracing**: Use a tool like New Relic, along with OpenTelemetry, to provide comprehensive monitoring and tracing from the middleware to the database layer. This will track the entire request flow across services and detect performance bottlenecks.

- **Advanced Testing**: Increase test coverage beyond the payment package, covering the wallet service and API layers as well.

- **Data Security**: Add encryption and hashing to secure sensitive data, such as transaction amounts, in the database.

- **External Queue System**: Replace the in-memory queue with a distributed message broker like RabbitMQ or Kafka for handling large-scale transactions.

These improvements will enhance the system's scalability, observability, and security in the future.

---

This **Wallet Service** is designed with **scalability**, **extensibility**, and **performance** in mind, enabling easy integration with additional gateways and providing a reliable platform for wallet and transaction management.
