package vo

// 返回给前端的用户信息
type LoginUserInfoVo struct {
	Token string `json:"token"`
}

type UserInfoVo struct {
	Id       int    `json:"id"`
	UserName string `json:"username"`
}
