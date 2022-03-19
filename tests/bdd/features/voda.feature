Feature: VODA
  VODA is an abbreviation for voice of diary

  VODA is an app to create an exchange diary.
  You can write your diary by attaching audio, text or images, and you can share it with your friends.
  For voice recording, you can also modulate your voice to record it if you wish.

  VODA는 교환일기를 작성하는 앱입니다.
  음성 또는 텍스트 또는 이미지들을 첨부하여 당신의 일기를 작성할 수 있으며, 이를 친구들과 공유할 수 있습니다.
  음성 녹음의 경우 원한다면 목소리를 변조하여 기록할 수도 있습니다.

  > tl;dr
  room = 교환일기방
  diary = 교환일기방에 생성되는 교환일기
  member = VODA의 회원 체계

  roomMaster = 교환일기방을 최초로 생성했거나, 양도받아 roomMember에서 승격된 존재
  roomMember = 특정 교환일기방에 참여하고 있는 멤버

  alarm = 교환일기방 알림
  task(event) = 알림에 필요한 이벤트
  file = 교환일기에 사용되는 static file (image / audio)
  terms = 회원가입 시 동의가 필요한 약관


  @room @roomMaster @task
  Scenario: Create a new diary room
  교환일기방을 새롭게 생성합니다.
  @room @roomMaster
  Scenario: Changing a diary room's hint / password
  교환일기방 참여 시 필요한 code / hint를 변경합니다.
  @room @roomMaster
  Scenario: Changing the order of writing exchange diaries between members
  멤버들 끼리 교환일기방을 작성하는 순서를 변경합니다.
  @room @roomMaster
  Scenario: Changing the period(cycle) of keeping a diary.
  멤버들 끼리 교환일기방을 작성하는 주기를 변경합니다.
  @room @roomMember
  Scenario: Join a diary room
  만들어진 교환일기방에 참여합니다.
  @room @roomMaster @roomMember
  Scenario: Leave a diary room
  생성했거나 참여했던 교환일기방에서 나갑니다.
  @room @roomMaster @roomMember
  Scenario: Check the list of diary-rooms I have joined.
  생성했거나 참여했던 교환일기방의 목록을 확인합니다.

  @diary @task
  Scenario: Create a diary (without audio / image)
  음성과 이미지가 없는 교환일기를 작성합니다.
  @diary @file @task
  Scenario: Create a diary (with unmodulated audio)
  변환되지 않은 오디오를 첨부한 교환일기를 작성합니다.
  @diary @file @task
  Scenario: Create a diary (with modulated audio)
  변환된 오디오를 첨부한 교환일기를 작성합니다.
  @diary @file @task
  Scenario: Create a diary (with image)
  이미지를 첨부한 교환일기를 작성합니다.
  @diary @file @task
  Scenario: Create a diary (with audio / image)
  오디오와 이미지를 첨부한 교환일기를 작성합니다.

  @member @terms
  Scenario: Sign up to voda app
  voda에 회원가입합니다.
  (닉네임 검사, oauth 어러개인경우 체크)
  @member
  Scenario: Log in to voda app
  voda에 로그인합니다.
  @member
  Scenario: Log out to voda app
  voda에 로그아웃합니다.
  @member @terms
  Scenario: Log in to voda app
  voda를 탈퇴합니다.

  @alarm @task
  Scenario: Time to write a diary.
  교환일기 작성해야 하는 순서가 되었습니다.
  @alarm @task
  Scenario: 1hour before I submit my exchange diary.
  교환일기를 제출하기 1시간 전입니다.
  @alarm @task
  Scenario: 4hour before I submit my exchange diary.
  교환일기를 제출하기 4시간 전입니다.
  @alarm @task @room
  Scenario: Scheduled exchange diary writing cycle (cycle) has arrived
예약해두었던 교환일기 작성주기(사이클)가 되었습니다.