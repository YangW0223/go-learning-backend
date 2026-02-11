package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword 对明文密码做 bcrypt hash。
// 默认成本由 bcrypt.DefaultCost 决定，兼顾安全与性能。
func HashPassword(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// ComparePassword 验证明文密码和 hash 是否匹配。
// 当密码不匹配时会返回非 nil 错误。
func ComparePassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
