package kakao

const AuthType = "kakao"

type KakaoUser struct {
	Account KakaoAccount `json:"kakao_account"`
	Profile KakaoProfile `json:"properties"`
}

type KakaoProfile struct {
	NickName     string `json:"nickname,omitempty"`
	ProfileImage string `json:"profile_image,omitempty"`
}

type KakaoAccount struct {
	Email string `json:"email"`
}
