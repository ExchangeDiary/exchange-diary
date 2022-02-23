# exchange-diary

> `exchange-diary` backend server

<div align="center">
  <img width="350" height="450" src="https://user-images.githubusercontent.com/37536298/153554715-f821d0f8-8f51-4f4c-b9e6-a19e02ecb5c2.png" />
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

- application: application layer
  - 원래는 `interface`라고 명칭을 가져가야 하지만, 코드의 interface와 명칭이 중복되어, application 영역으로 명시함.
  - `hexagonal`에서 application service layer + interface layer의 코드가 들어있음
- domain: domain layer
- infrastructure: infra layer

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

## Branching strategy with phase

> 아직 proposal

### Phase

- sandbox
  - heroku + gcs + db(free spec db를 찾아야해서 postgreSQL 또는 Google Cloud SQL)
- prod
  - static server: [google cloud run](https://cloud.google.com/run/docs/quickstarts/build-and-deploy/deploy-go-service) + gcs
  - voda-api server: [google cloud platform](https://cloud.google.com/gcp/?hl=ko&utm_source=google&utm_medium=cpc&utm_campaign=japac-KR-all-ko-dr-bkws-all-all-trial-e-dr-1009882&utm_content=text-ad-none-none-DEV_c-CRE_540744488055-ADGP_Hybrid%20%7C%20BKWS%20-%20EXA%20%7C%20Txt%20~%20GCP%20~%20General_cloud%20-%20platform-KWID_43700061085499317-kwd-87853815&userloc_1009893-network_g&utm_term=KW_gcp&gclid=CjwKCAiAx8KQBhAGEiwAD3EiPwkav_JRSG7KTofwsR--hAIxPecczxpKkym85b6z7IwENYDQxz-K7xoC8FIQAvD_BwE&gclsrc=aw.ds) + db(google cloud sql)

### Branching

[github-flow]() 어떨까?

> sandbox 페이즈는 feature가 main에 머지될떄 자동 배포되며, production은 sandbox가 정상 동작하면 수동배포하는 방식

- main: prod
- feature/<칸반보드이슈번호>

## Erd

![voda v1 erd](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/ExchangeDiary/exchange-diary/main/docs/erd.puml)

## Room flow

> 다이어리방 관련된 플로우

![room api](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/ExchangeDiary/exchange-diary/main/docs/rooms.puml)

## Diary flow

> 다이어리 관련된 플로우

![diary api](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/ExchangeDiary/exchange-diary/main/docs/diaries.puml)
