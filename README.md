# Ticket Allocation Coding Test

## Project Overview

This is my attempt at the [ticket allocation coding test](https://github.com/Fatsoma/ticket_allocation_coding_test/).

The API has been built using Go (v1.24) and PostgreSQL (v17). It allows ticket options to be created and read, and allows 
purchases to be made from a ticket option while ensuring that the ticket allocation does not drop below zero when concurrent
purchase requests are made.

### Concurrency
I considered two options for concurrency control:
- Explicit row-level locking
- Specifying a stricter transaction isolation level

In this case I opted for explicit locking. 

Each ticket option corresponds to a single row in the `ticket_options` table, and so I thought it was conceptually simple to lock that row when processing a purchase. If multiple rows, or entire tables, needed to be locked then setting the transaction isolation level to `REPEATABLE READ` or `SERIALIZABLE` would probably be preferable, to increase throughput. (Though even in the case of single-row locking, `SERIALIZABLE` transactions might have better performance because of the overhead of explicit lock management).

A benefit of explicit locking is that each transaction will either succeed or fail on its first attempt, and so we don't need any extra logic for catching `serialization failure` errors and retrying the transactions.

## Getting Started

### Local Development Setup

1. Clone the repository
2. Ensure Docker, Docker Compose and Go are installed

#### Starting the Application

```bash
# Start the PostgreSQL databases
docker compose up -d

# Run the application
make run

# Remove the database containers
docker compose down
```

The API will be available at `http://localhost:3000`

#### Running Tests

```bash
# Start the PostgreSQL databases if not done already
docker compose up -d

# Run database tests
make test

# Remove the database containers
docker compose down
```

## Project Structure
- `main.go`: Application entry point
- `routes.go`: HTTP route definitions
- `handlers.go`: Request handlers
- `db.go`: Database interaction methods
- `models/`: Data models and request/response structures
- `docker-compose.yml`: Database container configuration
- `Makefile`: Convenience commands for running and testing

## Notes
- The code layout is inspired what's done in [Let's go](https://lets-go.alexedwards.net/), which I quickly read through before starting this project. 
- If I were doing this again I'd want to try a more functional approach. Attaching methods to structs feels more OOP-like than I'm used to. A lot of the ideas in [this blog post](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/) sounded useful.
- I'd also want to try out Interfaces properly.

## Future Improvements
- Lots more tests, I've only focussed on testing the functions that interact with the database (`db_test.go`), to confirm that the concurrency control works as expected. None of the handlers or util functions are tested at the moment.
- Improve input validation. There's a lot of duplication in the handler functions related to validation the JSON input.
- In finance, [double-entry accounting](https://beancount.github.io/docs/the_double_entry_counting_method.html#introduction) is used. Is that something that could also be applied to ticket sales to improve auditability, error detection & reporting? 
