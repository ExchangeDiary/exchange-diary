# <strike>`API` (v1) deprecated</strike>

- 클라 통신에 필요한 api 스펙들을 정리합니다.
- API 서버의 base_url은 **"https://voda-api.com"** 라고 가정합니다.
- [voda figma](https://www.figma.com/file/NcqZIf2XzxeufnIyX8iMPu/%EA%B3%A0%EC%98%81%ED%9D%AC%EB%AF%B8%EB%A7%8C%EB%8B%A4%EA%BE%B8%EB%9F%AC?node-id=6%3A2)를 기준으로 기능들을 분류하였습니다.

## 1. 로그인 / 프로필 설정 (가입/탈퇴 약관 필요)

> 약관 서빙 / 약관 동의 post / email은 unique해야하며 email: oauth = 1:1? no 다중 로그인 가능하게 구현

1. 구글/apple/카카오 클릭
2. 닉네임 / 프로필 사진 설정
3. 완료 클릭

## 2. 알림 리스트

> TODO

## 3. 알림 보내기

> > TODO

## 4. MY(알림 설정)

1. 작성 차례 알림 on/off (1시간 / 4시간)
   - PATCH /v1/users/<:user_id>
   ```
   header: {Authorization: JWT_ACCESS_TOKEN}
   body:{
       "turnAlarmFlag": <bool>
   }
   ```
2. 활동 알림 on/off
   - PATCH /v1/users/<:user_id>
   ```
   header: {Authorization: JWT_ACCESS_TOKEN}
   body:{
       "activityAlarmFlag": <bool>
   }
   ```

## 5. MY(공지사항)

1. (어드민) 공지 사항 생성
   - TODO
2. (어드민) 공지 사항 수정
   - TODO
3. (어드민) 공지 사항 삭제
   - TODO
4. 공지 사항 리스트

   - 1달안에 만들어진 공지사항일 경우 (New)를 둔다.
   - 최근 생성된 순으로 정렬해서 보내준다.
   - 페이지네이션 없다.
   - GET /v1/notices

   ```
    {
        "notices": [
            {
                "noticeId": <int>,
                "title": <str>,
                "body": <text>,
                "createdAt": <datetime>
            },
            {},
            {},
            ...
        ]
    }
   ```

5. 공지 사항 detail

- list에서 보내준 데이터로 보여주기 때문에, 따로 api 없다.

## 6. MY(로그 아웃)

## 7. MY(회원 탈퇴)

## 8. Home

## 9. 교환일기방 생성 (다이어리 만들기)

- POST /v1/rooms/diaries
- request

```
header: {Authorization: JWT_ACCESS_TOKEN}
body:
{
    "theme": <THEME_ID>, // 1,2,3
    "name": <str>,
    "period": <int>, // 1,2,3,4,5 .... 00주기
    "code": <str>,   // ex) 레오제이1!
    "hint": <>       // 요즘핫한 뷰티유튜버는?
}
```

- response

```
{
    "roomId": <int>
}
```

## 10. 교환일기방 상세화면

- Q? api camelCase로 주면 될까요?
- Q? 해시태그는 어디서 나온 친구?
- Q? 파랑색 테두리는 본인일 경우인건가요? 아니면 현재 작성해야하는 사람인가요?
- Q? "아직*0턴", "우린이제*시작하는사이"는 어디서 나오는 정보?
- Q? "고영희\_자랑하는 일기" vs "고영희\n미만 다꾸러"의 차이는?
- Q? \_ 언더바가 꼭필요한걸까?
- Q? "000님이 이야기를 쓰고있어요!"는 현재 작성해야 하는 사람을 말하는 걸까?
- Q? 탈퇴한 경우 사용자 프로필 사진default 이며 기존에 작성했던 diary그대로 유지하는게 어떨까요?
- Q? "방구석에서"라는 문구는 장소를 보여주면 되는거죠?
- Q? 다이어리 회색 라인들은 다이어리의 "글 내용" 요약인가요? 아니면 그냥 저렇게 UI?
- GET /v1/rooms/<:room_id>
- request

```
header: {Authorization: JWT_ACCESS_TOKEN}
```

- response

```
{
    // 참여된 순으로 정렬 (방장 -> 멤버1, 멤법2...)
    // 클라에서는 accountId를 비교해서 현재 로그인한 사용자 확인가능
    "members": [
        {
            "accountId": <int>,         // account_id이며, 본인이면 파랑색
            "profileUrl": <url> // https://voda-api.com/profiles/leoo.png
        },
        {
            "accountId": <int>,
            "profileUrl": <url> // https://voda-api.com/profiles/leoo.png
        },
        ...
    ],
    "turnAccountName": <str>, //i.g "고영희", ("고영희님이 이야기를 쓰고있어요!"일때 사용)
    "tagName": <str>, //ig) #고영희_자랑하는 일기,
    // 최근 생성된 순으로 정렬
    "diaries": [
        {
            "place"     :<str>,         // 장소 ig) "방구석에서"
            // 지금은 필요없지만, 사용자 프로필 클릭했을 경우, 유저 detail 보여질수도 있음
            "accountId" :<int>,
            "accountName" :<str>,         // 생성한 accountId
            "profileUrl"   :<url>,         // 사용자 프로필 이미지 url
            "createdAt" :<datetime>,    //생성 일자
        },
        {},
        {},
        ...
    ]
}
```

## 11. 초대

> 모바일 접근 / pc접근에 따른

- Q? 정책: 모든 다이어리방 접근 시, 멤버가 아니라면(서버에서 에러제공), 클라에서는 방의 /invite api를 쏴주어야한다.
- Q? 초대 링크를 외부에서 클릭했을 때, 회원가입이 유도 flow 필요.

- 초대 링크 보내기(공유하기)

```
https://voda-api.com + /v1/rooms/<room_id>
url: https://voda-api.com/v1/rooms/19
```

- 초대 링크를 외부에서 클릭했을 때

  1. (페이지 필요)회원이 아닌경우
  2. 회원가입후 해당 방으로 redirection
  3. GET https://voda-api.com/v1/rooms/<:room_id>

  - 멤버일 경우: 교환일기방 상세화면 렌더링
  - 멤버가 아닐경우: **서버에서 401 에러 return**

  4. 401일 경우, 클라에서 참여코드 렌더링을 위해 필요한 api (힌트와 code)

     - GET **https://voda-api.com/v1/rooms/<:room_id>/invite**
     - request: header: {Authorization: JWT_ACCESS_TOKEN}
     - response

     ```
     {
        "tagName": <str>, //ig) #고영희_자랑하는 일기,
        "hint": <str>, // 요즘핫한 뷰티유튜버는?
        "code": <str>, // i.g) 레오제이1! (hash는 나중에 하자.)
        "totalMemberCount": <int>, // 현재참여인원 extra 숫자 보여주기 용
        // 최대 4명까지 제공, 참여된 순으로 정렬
        "members": [
            {
                "accountId": <int>,         // account_id이며, 본인이면 파랑색
                "profileUrl": <url> // https://voda-api.com/profiles/leoo.png
            },
            {
                "accountId": <int>,
                "profileUrl": <url> // https://voda-api.com/profiles/leoo.png
            },
            ...
        ]
     }
     ```

  5. 참여코드 입력 페이지 렌더링
  6. 클라에서 서버에서 전달해준 code와 사용자 입력 code를 비교해서 validation 진행
  7. 만약 validation ok이면, 서버로 멤버 초대 요청 필요

     - request

       ```
       POST https://voda-api.com/v1/rooms/<:room_id>/

       header: {Authorization: JWT_ACCESS_TOKEN}
       body: {} //account_id는 jwt에서 얻으면 된다.
       ```

## 12. 교환일기 작성

> 이미지와 오디오와 diary를 어떻게 효율적으로 저장할 수 있을까?

1. 현재 글 작성할 사용자인지 확인 api
   - GET https://voda-api.com/v1/rooms/<:room_id>/members/status
   - 로그인한 사용자에 대해서만 확인 가능
   ```
   header: {Authorization: JWT_ACCESS_TOKEN}
     {
        "isMaster": <bool>, // 마스터여부
        "isTurn": <bool>, //현재 turn 여부
     }
   ```
2. "turn" true 여부 확인
3. 글 폼 작성 / validation
4. POST
   - 클라에서 diaryUUID 만들어서 전달해줘야함
   - 이미지 데이터: multipart/form-data; (imageID return)
   - audio: POST multipart/form-data;
   - diaryUUID와 imageID, audioID전달해줘야함
   - redis에 저장도 가능할까?

- 작성중 페이지를 벗어나게되면 cache destroy 필요. /api/v1/

## 13. 교환일기 상세화면

- Q? 교환일기 상세 페이지가 여러 버전들이 있는 것 같은데, 음성/글+사진/글/
- Q? 교환일기 상세화면 > 음성변조선택 페이지는 어떻게 들어가야 할까요?
- Q? 템플릿 때문에 글씨 / 사진 겹쳐서 영역 침범 안하나용?

GET https://voda-api.com/v1/rooms/<:room_id>/diaries/<:diary_id>

```
   header: {Authorization: JWT_ACCESS_TOKEN}
   body: {
        "title": <str>,
        "place": <str>,
        // 지금은 필요없지만, 사용자 프로필 클릭했을 경우, 유저 detail 보여질수도 있음
        "accountId" :<int>,     // 생성한 accountId
        "accountName" :<str>,   // 생성한 account name      (탈퇴 경우 "OOO"로 대체)
        "profileUrl"   :<url>,  // 생성한 사용자 프로필 이미지 url (탈퇴 경우 default url로 대체)

        "createdAt": <datetime>,
        "body"   :<text>, // 글 body
        "photos" :[str, str...], // imageUrl(string) 리스트
        "audio" : {
            "id": ,
            "fileName": <str>,
            "length": <int>, //second 기준
            "content": []byte //실제 오디오 파일 (TODO: 이게 가능한지 확인해봐야함.)
        }
   }
```

## 14. 방장 설정 화면

들어가기 앞서

1. 방장인지 확인한다. GET https://voda-api.com/v1/rooms/<:room_id>/members/status

```
   header: {Authorization: JWT_ACCESS_TOKEN}
     {
        "isMaster": <bool>, // 마스터여부
        "isTurn": <bool>, //현재 turn 여부
        "state": <ACTIVE | INACTIVE>, // 멤버 활성화여부
     }
```

2. 필요한 데이터(작성기한) 가져온다. **GET /v1/rooms/<:room_id>/period**

```
{
    "period" : <int>
}
```

### 14.1. 비밀번호 변경(참여코드 재설정 화면)

PATCH **https://voda-api.com/v1/rooms/<:room_id>**

```
{
    "code": <str>,
    "hint": <str>,
}
```

### 14.2. 작성 순서 변경

- 현재 순서 가져오기

  - GET **https://voda-api.com/v1/rooms/<:room_id>/order**

  ```
  {
      "orders": [19, 18, 15...] // <[int]>타입 accountID가 순서에 맞춰서 제공됨
      // accountID: 이름 / 프로필이미지 url
      // 주의! json은 key가 항상 string이어야 한다.
      "accounts": {
          "19": {
            "name": <str>,      // leoo.j
            "profileUrl": <url> // https://voda-api.com/profiles/leoo.png
          },
          "18": {},
          "15": {},
        ...
      }

  }

  ```

- "완료" 클릭시

  - PUT **https://voda-api.com/v1/rooms/<:room_id>/order**

  ```
  {
     "orders": [19, 18, 15...] // accountID가 순서에 맞춰서 제공됨
  }
  ```

### 14.3. 작성 기한 변경

PATCH **https://voda-api.com/v1/rooms/<:room_id>**

```
{
    "period": <int>, // 1,2,3,4,5 .... 00주기
}
```

## 15. 음성녹음 화면

> TODO

- Q? 음성녹음 화면은 어떻게 접근이 가능한가요?
  - 교환일기 작성해주는 시점에 "음성파일 추가하기"누르면 진행되는 건가요?
- CRUD 필요
