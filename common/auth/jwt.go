package auth

import (
	"strconv"
	"time"

	"common/xerr"

	"github.com/golang-jwt/jwt/v4"
)

// Claims JWT 载荷。字段与 go-zero rest.WithJwt 的密钥名保持一致：
// userId 为数值型 claim，go-zero 鉴权中间件会将其注入 ctx。
type Claims struct {
	Role   string `json:"role"`
	UserId int64  `json:"userId"`
	jwt.RegisteredClaims
}

// Sign 签发访问令牌。
func Sign(userID int64, role, secret string, expireSeconds int64) (string, error) {
	now := time.Now()
	claims := Claims{
		Role:   role,
		UserId: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expireSeconds) * time.Second)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", xerr.Wrap(err, xerr.CodeInternal, "签发令牌失败")
	}
	return signed, nil
}

// Parse 校验并解析访问令牌，返回 userId 与 role。
func Parse(tokenString, secret string) (int64, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, xerr.Unauthorized("")
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return 0, "", xerr.Wrap(err, xerr.CodeUnauthorized, "")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, "", xerr.Unauthorized("")
	}
	userID := claims.UserId
	if userID <= 0 {
		parsed, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil || parsed <= 0 {
			return 0, "", xerr.Unauthorized("")
		}
		userID = parsed
	}
	role := claims.Role
	if role == "" {
		role = RoleUser
	}
	return userID, role, nil
}
