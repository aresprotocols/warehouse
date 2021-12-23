package service

import (
	conf "price_api/price_server/config"
	"price_api/price_server/internal/pkg/jwt"
	"price_api/price_server/internal/util"
	"price_api/price_server/internal/vo"
)

type AuthService struct {
}

func newAuth() *AuthService {
	return &AuthService{}
}

func (s *AuthService) ValidateUserAndPassword(user vo.AdminUser) bool {
	md5Password := util.Md5Str(conf.GCfg.Password)

	if user.User != conf.GCfg.User || user.Password != md5Password {
		return false
	} else {
		return true
	}
}

func (s *AuthService) GenerateToken(user vo.AdminUser) (string, error) {

	return jwt.GenToken(user.User, []byte(conf.GCfg.Password))

}
