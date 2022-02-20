# Exchange-diary

## Build & Run

*Step 1. Build application using Makefile.*
```sh
make build
```
*Step 2. Run application with specific flag.*
```sh
./bin/exchange-diray -phase=${phase}
```
- ${phase} would be "dev", "production", or "sandbox".
- Example:
```sh
./bin/exchange-diray -phase=dev
```
