package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	//JwtKey for encrypt
	defaultKey    = "zomeisagoodguy"
	defaultIssuer = "zome"
)

//Jwt 统一的鉴权jwt数据
type Jwt struct {
	Key    string
	Token  string //加密后的串。解密的时候使用
	Claims interface{}
	jwt.StandardClaims
}

//Option Option
type Option func(o *Jwt)

func Claims(claims interface{}) Option {
	return func(o *Jwt) {
		o.Claims = claims
	}
}
func Key(key string) Option {
	return func(o *Jwt) {
		o.Key = key
	}
}
func Token(token string) Option {
	return func(o *Jwt) {
		o.Token = token
	}
}

// Expire expire
func Expire(expire int64) Option {
	return func(o *Jwt) {
		o.ExpiresAt = expire
	}
}

//Issuer Issuer
func Issuer(issuer string) Option {
	return func(o *Jwt) {
		o.Issuer = issuer
	}
}

//DefaultOptions default options
func DefaultOptions() *Jwt {
	o := &Jwt{}
	o.Key = defaultKey
	o.Issuer = defaultIssuer
	o.ExpiresAt = time.Now().Add(24 * 30 * time.Hour).Unix()
	return o
}

//New 产生token
func New(options ...Option) *Jwt {
	opts := DefaultOptions()
	for _, o := range options {
		o(opts)
	}
	return opts

}
func (o *Jwt) Encode() (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, o)
	token, err := claims.SignedString([]byte(o.Key))
	return token, err
}
func (o *Jwt) SetToken(token string) *Jwt {
	o.Token = token
	return o
}

//Decode 解析
func (o *Jwt) Decode() (map[string]interface{}, error) {
	tokenClaims, err := jwt.ParseWithClaims(o.Token, &Jwt{}, func(token *jwt.Token) (interface{}, error) {
		//返回秘钥
		return []byte(o.Key), nil
	})
	if tokenClaims != nil {
		//检查token的有效性tokenClaims.Valid,类型断言
		if claims, ok := tokenClaims.Claims.(*Jwt); ok && tokenClaims.Valid {
			return claims.Claims.(map[string]interface{}), nil
		}
	}

	return nil, err
}
