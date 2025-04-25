package serializer

// UserLogin 登录后返回给前端的用户数据 + 双令牌
type UserLogin struct {
	Account        string `json:"account"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Avatar         string `json:"avatar"`
	Authentication string `json:"authentication"` // 访问令牌
	RefreshToken   string `json:"refresh_token"`  // 刷新令牌
}

// UserInfo 用户基本信息数据（管理员）
type UserInfo struct {
	Account  string `json:"account"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}
