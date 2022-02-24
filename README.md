# ExchangeDiary_Server

> ExchangeDiary backend server

<div align="center">
  <img width="350" height="450" src="https://user-images.githubusercontent.com/37536298/153554715-f821d0f8-8f51-4f4c-b9e6-a19e02ecb5c2.png" />
</div>

- <strike>[V1 API](./docs/api.md)</strike>
- <strike>[Features](./docs/features.md)</strike>
- [Fixed Policy](./docs/fixed_policy.md)

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

![voda v1 erd](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/ExchangeDiary/ExchangeDiary_Server/main/docs/erd.puml)

## Room flow

> 다이어리방 관련된 플로우

### CRUD

![room crud api](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/ExchangeDiary/ExchangeDiary_Server/main/docs/rooms-crud.puml)

### ETC

![room etc](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/ExchangeDiary/ExchangeDiary_Server/main/docs/rooms-etc.puml)

## Diary flow

> 다이어리 관련된 플로우

![diary api](http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/ExchangeDiary/ExchangeDiary_Server/main/docs/diaries.puml)
