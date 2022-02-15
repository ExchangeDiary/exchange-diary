# ExchangeDiary_Server

> ExchangeDiary backend server

<div align="center">
  <img width="350" height="450" src="https://user-images.githubusercontent.com/37536298/153554715-f821d0f8-8f51-4f4c-b9e6-a19e02ecb5c2.png" />
</div>

- [API 스펙](./api.md)
- [ERD](https://lucid.app/lucidchart/a4191542-ece6-416c-879b-028973e51a7e/edit?invitationId=inv_ab6c32ee-9ef4-40be-b426-d2929fff1463)

## 1. Fixed Policy

> 논의된 정책을 리스트업합니다.

## 2. Need Policy

> 아직 정해지지 않은 정책 내용들을 적어봅니다.

### 2.1. Account

- 회원 탈퇴 시, 해당 회원이 DiaryRoom의 (마스터이다 && room에 다른 member가 존재), 탈퇴가 불가능해야 할까? 또는 마스터 권한 양도 기능이 있어야 할까?
- 계정 삭제 시, 실제 계정 데이터를 삭제해야 할까? (탈퇴 정책)
- 계정이 삭제되었다면 계정이 활동했던 데이터들은 어떻게 처리해야할까?

  - if master
    - DiaryRoom : 양도
    - Invitation : 삭제
      추가로 아래 member 고민 사항 적용
  - if member
    - DiaryOrder.order: 수정
    - DiaryMember: 제거
    - 받을 Alarm 제거
    - DiaryStory: 프로필 정보는 default를 내려주도록 설계 필요

### 2.2. DiaryRoom(교환일기방)

- DiaryRoom 삭제 기능은 없나요?
- DiaryRoom 페이지네이션 (아마 많아도 3~4개정도 만들지 않을까 싶어서 아직은 필요없을 듯)

### 2.3. DiaryMember (교환일기 멤버)

- 초대받은 멤버(친구) 내쫒기 기능은 없나요? (만약 생긴다면 관련 히스토리들은 어떻게 할지 논의 필요)

### 2.4. DiaryStory (교환일기)

- 변조된 음성 녹음 데이터는 따로 저장되는 서버가 있나요?
- DiaryStory는 수정/삭제 불가인가요? (음성 녹음 포함)
- _댓글 /like 기능은 크게 필요해보이진 않는다._

## 3. Features

> 도메인별 필요한 기능들을 리스트업합니다.

### 3.1. Account

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

### 3.2. DiaryRoom

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

### 3.3. Invitation (Deprecated)

권한은 code로 보면 되기때문에, sharable link를 생성해주고, 들어오는 모든 사용자들을 허용한다. **즉 초대장을 row로 관리할 필요가 없어 보인다.**

- 초대장 전송 (CREATE)
  - kakao / tms / sharable link
  - authorization: 방장 only
- 초대 받기 (DELETE)
  - code verify()
  - 초대를 받을 경우 delete 한다.
- (optional) 일정 기간 지난 초대장들 제거
  - 어차피 crontab worker가 필요하니 초대장

### 3.4. DiaryStory

- LIST
  - 교환일기방 내부화면 진입시 필요
  - api 엔드포인트는 없어도 될듯
  - authorization: 방장 / 멤버 only
- CREATE
  - authorization: 방장 / 멤버 && 현재 Order인 경우
- DELETE

### 3.5. DiaryOrder

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

### 3.6. Alarm (TODO)

> 이게 제일 어려울 듯

- server push to client (http2.0)
- crontab golang worker

### 3.7. Notice (TODO)

> admin 기능
