# Exchange-diary

## precommit hook

> [refs](https://tutorialedge.net/golang/improving-go-workflow-with-git-hooks/)

```bash
$ cp pre-commit.example ../.git/hooks/pre-commit
$ chmod +x ../.git/hooks/pre-commit
```

## Build & Run

_Step 1. Build application using Makefile._

```sh
make build
```

_Step 2. Run application with specific flag._

```sh
./bin/exchange-diray -phase=${phase}
```

- ${phase} would be "dev", "production", or "sandbox".
- Example:

```sh
./bin/exchange-diray -phase=dev
```
