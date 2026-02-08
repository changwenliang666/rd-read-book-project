package vo

// 返回给前端的用户信息
type UserInfoVo struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
