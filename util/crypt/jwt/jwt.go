package jwt

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	//JwtKey for encrypt
	JwtKey = "zomeisagoodguy"
)

const ()

const (
	Ot_Owner  = "owner"  //所有人，主人，上级
	Ot_Member = "member" //普通成员
	//其他名字在扩充

)

type JwtNumber uint32

// Options 统一的鉴权jwt数据
type Options struct {
	UID     uint32 `json:"uid"`       //鉴权用户的UID
	Account string `json:"account"`   //登录用户名或者鉴权 Identifier
	OwnType string `json:"user_type"` //用户类型。//需要区分等级的时候才需要用到
	Role    string `json:"role"`      //身份角色。"身份1,状态1|身份2,状态2"
	Pid     uint32 `json:"pid"`       //如果是子账号。这个表示关联主体的ID， principal id
	Rules   string `json:"rules"`     //权限节点。"节点1|节点2"
	Etype   string `json:"etype"`     //财务中心账户类型
	Org     string `json:"org"`       //组织。实际上这个在某些时候和pid可以更替

	//Test JwtNumber `json:"test"` //测试用
	//jwtClaims jwt.StandardClaims
	jwt.RegisteredClaims
}

// Option Option
type Option func(o *Options)

// UID set uid
func UID(uid uint32) Option {
	return func(o *Options) {
		o.UID = uid
	}
}
func (o JwtNumber) String() string {
	return string(strconv.Itoa(int(o)))
}
func (o JwtNumber) Uint32() uint32 {
	return uint32(o)
}

// Account account
func Account(account string) Option {
	return func(o *Options) {
		o.Account = account
	}
}

// OwnType usertype
func OwnType(usertype string) Option {
	return func(o *Options) {
		o.OwnType = usertype
	}
}

// Pid group
func Pid(pid uint32) Option {
	return func(o *Options) {
		o.Pid = pid
	}
}

// Pid group
func Org(org string) Option {
	return func(o *Options) {
		o.Org = org
	}
}

// Rules rules
func Rules(rules string) Option {
	return func(o *Options) {
		o.Rules = rules
	}
}

// Role roles
func Role(role string) Option {
	return func(o *Options) {
		o.Role = role
	}
}

// Expire expire
func Expire(expire int64) Option {
	return func(o *Options) {
		o.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expire)))
		//o.ExpiresAt = expire
	}
}

// Issuer Issuer
func Issuer(issuer string) Option {
	return func(o *Options) {
		o.Issuer = issuer
	}
}

// func App(app string) Option {
// 	return func(o *Options) {
// 		o.App = app
// 	}
// }

// DefaultOptions default options
func DefaultOptions() Options {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * 30 * time.Hour)

	o := Options{
		Pid:   0,  //GroupId
		Rules: "", //rules
		Role:  "",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "zome",
		},
	}
	return o
}

// DefaultOptions default options
func TestOptions() Options {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * 30 * time.Hour)

	o := Options{
		UID:   1,
		Pid:   1,  //GroupId
		Rules: "", //rules
		Role:  "",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    "zome",
		},
	}
	return o
}

// New 产生token
func New(opts ...Option) (string, error) {
	options := DefaultOptions()
	for _, o := range opts {
		o(&options)
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, options)
	token, err := claims.SignedString([]byte(JwtKey))
	return token, err

}

// Parse 分析token
func Parse(token string) (*Options, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Options{}, func(token *jwt.Token) (interface{}, error) {
		//返回秘钥
		return []byte(JwtKey), nil
	})

	if tokenClaims != nil {
		//检查token的有效性tokenClaims.Valid,类型断言
		if claims, ok := tokenClaims.Claims.(*Options); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
