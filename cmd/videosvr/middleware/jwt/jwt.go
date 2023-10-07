package jwtoken

import "github.com/golang-jwt/jwt/v4"

// JWT 签名密钥，定义JWT中存放usrname和userID

var mySigningKey = []byte("mini_tiktok")

type JWTClaims struct {
	Username string `json:"user_name"`
	UserID   int64  `json:"user_id"`
	jwt.RegisteredClaims
}

func GenToken(userid int64, username string) (string, error) {
	claims := &JWTClaims{
		Username: username,
		UserID:   userid,
		RegisteredClaims: jwt.RegisteredClaims{
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			// 方便测试，不设置过期时间
			Issuer: "server",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
