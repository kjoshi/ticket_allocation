# Ticket Sales API

## Project Overview

TODO Update this:
This is a lightweight ticket sales API built with Go, utilizing PostgreSQL as the database. The application allows for creating ticket options and purchasing tickets with allocation management.

## Getting Started

### Local Development Setup

1. Clone the repository
2. Ensure Docker and Docker Compose are installed

#### Starting the Application

```bash
# Start the PostgreSQL database and API
docker-compose up -d

# Run the application
make run
```

The API will be available at `http://localhost:3000`

#### Running Tests

```bash
# Start the PostgreSQL database if not done already
docker-compose up -d
# Run database tests
make test
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
- Something about Let's Go **TODO 
- Something about row-level locking vs Serializable transactions.
- Environment variables control database and server configuration
- Transactions ensure data consistency during ticket purchases
- Allocation is managed at the database level to prevent overselling

## Future Improvements
- Lots more tests.
- Experiment with transaction isolation level of `Serializable`
- Try a more functional approach, rather than methods on structs
- Double-entry accounting?
- Ticket buckets?
- Implement more comprehensive input validation
