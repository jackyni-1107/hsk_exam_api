package bo

type TokenPayload struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
}
