
# **Wallet Service**

This project implements a **wallet management system** integrated with external payment gateways. The service supports creating wallets, handling deposits and withdrawals, and processing transactions asynchronously through a queue.

## **Project Structure**

```bash
├── cmd/                         # Command entry points for the services
│   └── api/
│       ├── mock/                # Mock payment gateway services
│       │   ├── gateway-a/       # Mock for Gateway A (JSON)
│       │   └── gateway-b/       # Mock for Gateway B (XML)
│       └── wallet/              # Main wallet service entry point
├── config/                      # Service configuration (env variables, etc.)
│   └── config.go
├── docker/                      # Docker setup for building and running the services
│   ├── dockerfile.mocks
│   └── dockerfile.wallet
├── docker-compose.yml           # Docker Compose file to run all services
├── docs/                        # OpenAPI/Swagger documentation
│   ├── docs.go
│   └── swagger.yaml
├── internal/                    # Core business logic and internal packages
│   ├── db/                      # Database migration scripts
│   ├── dto/                     # Data Transfer Objects for requests/responses
│   ├── handlers/                # HTTP handlers for API endpoints
│   ├── models/                  # Database models (e.g., wallet, transaction)
│   ├── payment/                 # Payment processing package for gateways
│   ├── repos/                   # Repositories for interacting with databases
│   ├── services/                # Business logic services (wallets, transactions)
│   └── web/                     # Middleware, routing, and logging utilities
├── pkg/                         # Reusable packages and utilities
│   ├── database/                # Database connection helpers
│   ├── errs/                    # Custom error handling package
│   ├── logger/                  # Logging package
│   ├── queue/                   # In-memory queue for async processing
│   ├── rest/                    # HTTP client with retry functionality
│   └── web/                     # Web response helpers
```

## **Features**

- **Wallet Management**: Create wallets, and list wallets.
- **Transaction Handling**: Process deposits and withdrawals through external payment gateways.
- **Async Transaction Processing**: Offload transaction processing to an in-memory queue for improved performance.
- **Extensible Payment Gateways**: Easily add new payment gateways by implementing the `PaymentGateway` interface.
- **Mock Payment Gateways**: Includes two mock gateways (A and B) to simulate payment flows.

## **Architecture Overview**

The service is designed with modularity and extensibility in mind, separating business logic (wallet and transactions) from payment gateway integrations. This design enables smooth scaling and future enhancements.

### **Core Components**:

1. **Wallet Service**:  
   - Manages wallets, creates deposits/withdrawals, and interacts with the payment system.
   - Handles transaction processing asynchronously using an internal queue.
   - Transactions are pushed into the service queue for background processing.
   - Queue workers process transactions asynchronously, interacting with Payment package

2. **Payment Package**:  
   - Defines a unified interface for interacting with various payment gateways.
   - Provides functions for deposit, withdrawal, and callback handling.
   - Includes gateways client and each client have retry and circuit breaker functionality for robustness.

### **Flow of Deposit/Withdrawal**:

1. **Client Request**:  
   The user submits a **deposit** or **withdrawal** request through the API. The request is validated for correctness, including payment details validation.

2. **Transaction Creation**:  
   After validation, a new **transaction** is created with a status of `created` and stored in the database. This transaction is then placed into the in-memory queue for async processing.

3. **Asynchronous Processing**:  
   The queue worker picks up the transaction and triggers the appropriate Payment package function and the payment package with will interact with The payment gateway client (mock or real) processes the request and returns a status (e.g., `success`, `failed`, or `pending`).

4. **Callback Handling**:  
   In the case of asynchronous payment gateways, a callback URL is provided during the initial request. When the gateway confirms the transaction, the callback will triggerd. wallet service will pass it to payment package which will verify it with the payment gateway.

5. **Final Status**:  
   The transaction's final status is updated in the system and can be retrieved via the wallet API (e.g., `completed`, `failed`, or `pending`).

### **Extensibility**:
Adding a new payment gateway is seamless. Simply implement the `PaymentGateway` interface and plug it with its key into the `Payment` package in cmd.

```
// Payment gateway setup
paymentGateways := map[models.PaymentGateway]payment.PaymentGateway{
	models.PaymentGatewayA: gatewaya.New(cfg.PaymentGatewayConfig.GatewayA),
	models.PaymentGatewayB: gatewayb.New(cfg.PaymentGatewayConfig.GatewayB),
}

// Payment handler setup
paymentHandler := payment.New(paymentGateways)
```
The system can support multiple payment gateways without modifying core business logic.

```go
type PaymentGateway interface {
    Deposit(ctx context.Context, req *Request) (*Response, error)
    Withdraw(ctx context.Context, req *Request) (*Response, error)
    VerifyCallback(ctx context.Context, refID string, data []byte) (*Response, error)
    VerifyMethod(typ models.TransactionType, method models.PaymentMethod) error
}
```

### **Tables**:
#### Wallets Table
Holds wallet-relate, only ID in purpose if the assessment.
| Column      | Type       | Description                              |
|-------------|------------|------------------------------------------|
| id        | uuid     | Primary key (UUID) for the wallet         |
| created_at| timestamp| The timestamp when the wallet was created |
| updated_at| timestamp| The timestamp when the wallet was last updated |

#### Transactions Table
Stores all transaction details such as deposits and withdrawals, including their status and reference to the related wallet.
| Column                | Type                | Description                                                        |
|-----------------------|---------------------|--------------------------------------------------------------------|
| id                  | uuid              | Primary key (UUID) for the transaction                             |
| wallet_id           | uuid              | Foreign key referring to the wallet (UUID)                         |
| amount              | decimal(18, 2)    | The amount involved in the transaction (up to 2 decimal points)    |
| type                | transaction_type  | Enum indicating transaction type (deposit, withdrawal)         |
| status              | transaction_status| Enum indicating the status (created, pending, completed, failed) |
| payment_gateway      | payment_gateway   | Enum for the payment gateway used (gateway_a, gateway_b)        |
| payment_method_details| JSONB           | JSON containing masked details of the payment method (e.g., card details) |
| payment_method       | varchar(255)      | The method used for the payment (e.g., credit card, bank transfer) |
| reference_id         | varchar(255)      | Unique reference ID for the transaction (optional)                 |
| created_at           | timestamp         | When the transaction was created                                   |
| updated_at           | timestamp         | When the transaction was last updated                              |

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

   The wallet service will be available at `http://localhost:8080`.

3. **API Documentation**:
   - Access the OpenAPI (Swagger) documentation at `http://localhost:8080/swagger`.
   - Postman collection `wallet.postman_collection.json` 

## **Configuration**

The service configuration can be adjusted via the `config/config.go` and `.env` file. these include setting up the database, payment gateways, and other environment variables.

---

## **Tests**

- Unit tests are provided for the **payment package** (integration with mock gateways) due to time constraints.
- Run tests with:
   ```bash
   go test ./...
   ```

---

## **To Extend the Service**

To add new payment gateway just 2 steps follow these steps:
1. **Creat client that implement the `PaymentGateway` Interface** for new gateways.
2. **Add New Client with its key to paymentGateways in cmd** in the `cmd/api/wallet/main.go`

---

This **Wallet Service** is designed with **scalability**, **extensibility**, and **performance** in mind, enabling easy integration with additional gateways and providing a reliable platform for wallet and transaction management.

