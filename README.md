# Error Simulator

A Go HTTP service that **intentionally triggers various runtime errors** (panics, nil dereferences, deadlocks, etc.) for testing error handling, monitoring pipelines, and downstream consumers (e.g. an AI debugger that reads error events from Kafka).

## Do We Need This?

**Yes**, if you want to:

- **Test recovery/monitoring**: Verify that your panic recovery middleware, logging, or APM correctly captures and reports errors.
- **Test Kafka error pipeline**: Produce real error events (stack traces, service name, repo) to a topic consumed by another service (e.g. ai-debugger).
- **Demo or chaos testing**: Reproducibly trigger specific failure modes (nil pointer, index OOB, type assertion, etc.) in a controlled way.

If you do **not** use Kafka or an ai-debugger consumer, you can still run the server and hit the endpoints to see recovered panic responses (JSON) and stack traces in logs; Kafka is optional.

---

## Requirements

- **Go 1.21+**
- **Kafka** (optional): only needed if you want events published to a topic. If `KAFKA_BOOTSTRAP_SERVERS` is unset or Kafka is unavailable, the app still runs and logs that events were not published.

---

## Quick Start

```bash
# Install dependencies
go mod download

# Run (default port 8080, Kafka disabled if not configured)
go run .

# Or build and run
go build -o error-simulator .
./error-simulator
```

With Kafka (e.g. local broker):

```bash
export KAFKA_BOOTSTRAP_SERVERS=localhost:9092
export KAFKA_TOPIC=service.errors
go run .
```

Trigger an error:

```bash
curl http://localhost:8080/error/nil-pointer
```

---

## Configuration

| Variable                  | Description                    | Default                                                |
| ------------------------- | ------------------------------ | ------------------------------------------------------ |
| `SERVER_PORT`             | HTTP server port               | `8080`                                                 |
| `KAFKA_BOOTSTRAP_SERVERS` | Kafka broker list              | `localhost:9092` (if empty, Kafka publish is disabled) |
| `KAFKA_TOPIC`             | Topic for error events         | `service.errors`                                       |
| `GITHUB_REPOSITORY`       | Repository name sent in events | `error-simulator`                                      |

---

## API Endpoints

All endpoints are **GET**. Each one triggers a specific kind of error; the **recovery middleware** catches panics, publishes an event to Kafka (if configured), and returns a JSON error response with `500`. Exceptions: **deadlock** and **stack-overflow** are not recoverable and will terminate the process.

| Endpoint                | Error type         | Description                                                                                                                      |
| ----------------------- | ------------------ | -------------------------------------------------------------------------------------------------------------------------------- |
| `/error/nil-pointer`    | Nil pointer        | `OrderService.ProcessOrder` dereferences `order.Patient` when it is nil.                                                         |
| `/error/db`             | DB error           | `UserRepository.GetUserByID` uses a nil `*sql.DB`, causing a panic.                                                              |
| `/error/panic`          | Explicit panic     | `PaymentService.ProcessPayment` panics when amount exceeds limit.                                                                |
| `/error/index-oob`      | Index out of range | `ReportGenerator.GetTopProducts` accesses index 5 on a slice of length 3.                                                        |
| `/error/type-assertion` | Type assertion     | `ConfigLoader.GetDatabaseConfig` asserts `config["database"]` as `map[string]interface{}` but it is a string.                    |
| `/error/division-zero`  | Division by zero   | `MetricsService.CalculateConversionRate` divides by `totalVisits` when it is 0.                                                  |
| `/error/deadlock`       | Deadlock           | `CacheManager` uses two mutexes in opposite order in different methods; concurrent calls deadlock. **Fatal:** process exits.     |
| `/error/stack-overflow` | Stack overflow     | `TreeNode.CalculateDepth` recurses without a nil base case and is called with a self-referential node. **Fatal:** process exits. |

---

## Recovered Response Format

For recoverable panics, the server returns **500** with JSON like:

```json
{
    "error": "panic recovered",
    "error_message": "...",
    "error_type": "NilPointer",
    "kafka_sent": true,
    "timestamp": "2025-03-13T12:00:00Z"
}
```

`error_type` is one of: `NilPointer`, `DBError`, `Panic`, `IndexOOB`, `TypeAssertion`, `DivisionZero`, `Deadlock`, `StackOverflow`, or `Unknown`.

---

## Kafka Event Schema

Events published to the configured topic match the schema consumed by **ai-debugger**:

| Field           | Type   | Description                            |
| --------------- | ------ | -------------------------------------- |
| `service`       | string | Always `"error-simulator"`             |
| `repository`    | string | From `GITHUB_REPOSITORY`               |
| `branch`        | string | `"main"`                               |
| `error_message` | string | Panic message or string representation |
| `stack_trace`   | string | Full Go stack trace                    |
| `timestamp`     | string | RFC3339 UTC                            |
| `environment`   | string | `"development"`                        |

---

## Project Layout

```
.
├── main.go                 # Server setup, routes, graceful shutdown
├── config/
│   └── config.go          # Env-based config (Kafka, port, repo name)
├── middleware/
│   └── recovery.go        # Panic recovery, Kafka publish, JSON error response
├── handlers/              # One file per error scenario
│   ├── nil_pointer.go     # OrderService nil Patient
│   ├── db_error.go        # UserRepository nil DB
│   ├── panic_recovery.go  # PaymentService explicit panic
│   ├── index_oob.go       # ReportGenerator index out of range
│   ├── type_assertion.go  # ConfigLoader bad type assertion
│   ├── division_zero.go   # MetricsService divide by zero
│   ├── deadlock.go        # CacheManager mutex deadlock
│   └── stack_overflow.go  # TreeNode infinite recursion
├── kafka/
│   └── producer.go        # Singleton producer, PublishErrorEvent
└── models/
    └── models.go          # Order, User, Product, ErrorEvent, etc.
```

---

## License

Use and modify as needed for your organization.
