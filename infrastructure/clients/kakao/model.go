package kakao

// AuthType value for kakao
const AuthType = "kakao"

// User struct for kakao
type User struct {
	Account Account `json:"kakao_account"`
	Profile Profile `json:"properties"`
}

// Profile for kakao
type Profile struct {
	NickName     string `json:"nickname,omitempty"`
	ProfileImage string `json:"profile_image,omitempty"`
}

// Account for kakao
type Account struct {
	Email string `json:"email"`
}
