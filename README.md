# VODA 📙

> VODA is an abbreviation for voice of diary

<div align="center">
  <img width="150" height="250" src="https://user-images.githubusercontent.com/37536298/153554715-f821d0f8-8f51-4f4c-b9e6-a19e02ecb5c2.png" />
</div>

`VODA` is an app to create an exchange diary.
You can write your diary by attaching audio, text or images, and you can share it with your friends. For voice recording, you can also modulate your voice to record it if you wish.

`VODA`는 교환일기를 작성하는 앱입니다. 음성 또는 텍스트 또는 이미지들을 첨부하여 당신의 일기를 작성할 수 있으며, 이를 친구들과 공유할 수 있습니다. 음성 녹음의 경우 원한다면 목소리를 변조하여 기록할 수도 있습니다.

더 자세한 정책은 다음에서 확인가능합니다. [about voda policies](./docs/fixed_policy.md)

## Terminology (domain)

- `room` = 교환일기방
  - `roomMaster` = 교환일기방을 최초로 생성했거나, 양도받아 roomMember에서 승격된 존재
  - `roomMember` = 특정 교환일기방에 참여하고 있는 멤버
- `diary` = 교환일기방에 생성되는 교환일기
- `member` = VODA의 회원 체계

- `alarm` = 교환일기방 알림
- `task`(event) = 알림에 필요한 이벤트
- `file` = 교환일기에 사용되는 static file (image / audio)
- `terms` = 회원가입 시 동의가 필요한 약관

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

## Precommit-hook

> [refs](https://tutorialedge.net/golang/improving-go-workflow-with-git-hooks/)

### .zshrc or .bashrc

- go mod를 사용할 경우

```sh
... 중략 ...
export PATH="$PATH:$HOME/go/bin"
export GO111MODULE=on
```

```bash
$ go install golang.org/x/tools/cmd/goimports
$ go install golang.org/x/lint/golint

$ cp pre-commit.example .git/hooks/pre-commit
$ chmod +x .git/hooks/pre-commit
```

## Cmd

### local

- go run + local mysql db

```sh
$ make run
$ make build
$ make docker
$ make clean

$ ./bin/exchange-diray -phase=${phase}
```

### sandbox

- local docker api server image + google cloud sql

```sh
$ make down && make up
```

### prod

- google cloud run + google cloud sql (same as sandbox db)
- trigger (cloud build)
  - **push to /main branch**

## Deploy env

- api server: `google cloud run`
- static server: `google cloud storage FUSE`
- db: `google cloud sql`
- ci / cd: `Cloud Code` & `Cloud Build`
- devops
  - `Cloud Monitoring`
  - `Cloud Logging`
  - `Cloud Trace`

## Phase

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
