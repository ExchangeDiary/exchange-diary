# exchange-diary

> `exchange-diary` backend server

<div align="center">
  <img width="250" height="350" src="https://user-images.githubusercontent.com/37536298/153554715-f821d0f8-8f51-4f4c-b9e6-a19e02ecb5c2.png" />
</div>

- <strike>[V1 API](./docs/api.md)</strike>
- <strike>[Features](./docs/features.md)</strike>
- [Fixed Policy](./docs/fixed_policy.md)

## Structure

> Domain Driven Design (Hexagonal architecture)

![](https://github.com/Sairyss/domain-driven-hexagon/blob/master/assets/images/DomainDrivenHexagon.png)

```bash
├── application         // same as interface layer in hexagonal architecture
│   ├── cmd             // command line interface
│   ├── controller      // http controller
│   ├── middleware      // middleware that handles requests
│   └── route           // http route, which delegate impl to controller
├── domain   // domain layer in hexagonal architecture, never have any external dependencies
│   ├── entity  // entity in domain layer
│   ├── repository  // interface of persistence layer
│   └── service     // domain service layer
└── infrastructure  // handle external dependencies
    ├── configs     // every configs include gin framework
    └── persistence // impl of persistence layer
```

> **GOLDEN_RULE: domain/ 에는 외부 dependency가 들어오지 않는다.**

- `application/`: application layer
  - 원래는 `interface`라고 명칭을 가져가야 하지만, 코드의 interface와 명칭이 중복되어, application 영역으로 명시함.
  - `hexagonal`에서 application service layer + interface layer의 코드가 들어있음
- `domain/`: domain layer
- `infrastructure/`: infra layer

## precommit-hook

> [refs](https://tutorialedge.net/golang/improving-go-workflow-with-git-hooks/)

### .zshrc or .bashrc

- go mod를 사용할 경우

```sh
... 중략 ...
export PATH="$PATH:$HOME/go/bin"
export GO111MODULE=on
```

```bash
$ go get golang.org/x/tools/cmd/goimports
$ go get -u golang.org/x/lint/golint

$ cp pre-commit.example .git/hooks/pre-commit
$ chmod +x .git/hooks/pre-commit
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

## Deploy env

- api server: `google cloud run`
- static server: `google cloud storage FUSE`
- db: `google cloud sql`
- ci / cd: `Cloud Code` & `Cloud Build`
- devops
  - `Cloud Monitoring`
  - `Cloud Logging`
  - `Cloud Trace`

## Erd

![voda v1 erd](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/ExchangeDiary/exchange-diary/main/docs/erd.puml)

## Room flow

### CRUD

> 다이어리방 생성 / 읽기 / 업데이트 / 나가기 관련 플로우

![room crud api](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/ExchangeDiary/exchange-diary/main/docs/rooms-crud.puml)

### ETC

> crud를 제외한 나머지 api 플로우

![room etc](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/ExchangeDiary/exchange-diary/main/docs/rooms-etc.puml)

## Diary flow

> 다이어리 관련된 플로우

![diary api](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/ExchangeDiary/exchange-diary/main/docs/diaries.puml)
