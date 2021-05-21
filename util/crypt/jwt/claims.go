package jwt

//很常见的一个struct ，可以自己定义
type UserClaims struct {
	UID      uint32 `json:"uid"`       //鉴权用户的UID
	Account  string `json:"account"`   //登录用户名或者鉴权 Identifier
	UserType string `json:"user_type"` //用户类型。admin,user,deliver,merchant .暂时大概这些类型
	Role     string `json:"role"`      //身份角色。"身份1,状态1|身份2,状态2"
	GroupID  uint32 `json:"group_id"`
	Rules    string `json:"rules"`
}
