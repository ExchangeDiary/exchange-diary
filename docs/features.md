# Features

> 도메인별 필요한 기능들을 리스트업합니다.

TODO: [fixed policy](./fixed_policy.md)에 맞게 문서 업데이트 필요합니다.

## 1. Account

- 회원가입을 한다.
  - profile image
  - nickname
  - select oauth sns
- 로그인 상태를 유지한다.
  - jwt
    - access token
    - refresh token
- 로그인을 한다.
- 로그아웃을 한다.
- 회원 탈퇴를 한다. (Destroy)
  - validation
  - invitation delete (pending -> delete)
    - from invitation
    - to invitation
  - DairyOrder.order update
  - DiaryMember delete
    - 속해 있는 diary들을 탈퇴한다. (DiaryMember delete, )
  - 받을 alarm delete
  - account delete
- 프로필 > 알림 설정 수정 (Update)
  - 작성 차례 알림 on/off
  - (새 글) 활동 알림 on/off
- 프로필 > 이미지 수정 (Update)
- 프로필 정보 보여주기 (Read)

## 2. DiaryRoom

- 교환일기방 생성 양식 전달 (Default)
  - theme dropdown list 전달
  - 작성 주기 dropdown list 전달 (서버에서 관리해야함)
- 교환일기방 생성 (Create)
  - validate()
    - !참여코드 / 힌트도 생성시점에 전달 받아야함
  - create()
  - (optional) 참여코드 encrypt
- 가입된 다이어리 리스트 보여주기 without pagination (List)
- 교환일기방 내부화면 (Read)
  - authorization: 방장 / member only
  - 교환일기 list 요청
- 방장 설정화면 (= 교환일기방 수정)
  - authorization: 방장 only
  - 작성 순서 변경 (update DiaryOrder)
  - 작성기한 변경
  - 비밀번호 변경

## 4. DiaryStory

- LIST
  - 교환일기방 내부화면 진입시 필요
  - api 엔드포인트는 없어도 될듯
  - authorization: 방장 / 멤버 only
- CREATE
  - authorization: 방장 / 멤버 && 현재 Order인 경우
- DELETE

## 5. DiaryOrder

- CREATE
  - Diary 생성 시점
- UPDATE
  - case 1: 새로운 유저가 들어왔을 경우
  - case 2: 방장이 순서를 변경했을 경우
  - case 3: 멤버 또는 방장이 회원 탈퇴 할 경우
- DELETE
  - cascade delete with DiaryRoom
- READ
  - authorization: 방장 / 멤버 only
  - case 1: 방장이 순서 수정 페이지로 접근 시
  - case 2: 교환일기 READ가 호출 될 때

## 6. Alarm (TODO)

> 이게 제일 어려울 듯

- server push to client (http2.0)
- crontab golang worker

## 7. Notice (TODO)

> admin 기능
