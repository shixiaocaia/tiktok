package dao

import (
	"github.com/shixiaocaia/tiktok/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserNameIsExist(username string) (bool, error) {
	db := GetDB()
	user := model.User{}
	err := db.Where("user_name = ?", username).First(&user).Error
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func InsertUser(username, password string) (*model.User, error) {
	db := GetDB()
	// 加密密文，明文存储密码不安全
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// mysql创建用户
	user := model.User{
		Name:            username,
		Password:        string(hashPassword),
		Follow:          0,
		Follower:        0,
		TotalFav:        0,
		FavCount:        0,
		Avatar:          "https://tse1-mm.cn.bing.net/th/id/R-C.d83ded12079fa9e407e9928b8f300802?rik=Gzu6EnSylX9f1Q&riu=http%3a%2f%2fwww.webcarpenter.com%2fpictures%2fGo-gopher-programming-language.jpg&ehk=giVQvdvQiENrabreHFM8x%2fyOU70l%2fy6FOa6RS3viJ24%3d&risl=&pid=ImgRaw&r=0",
		BackgroundImage: "https://tse2-mm.cn.bing.net/th/id/OIP-C.sDoybxmH4DIpvO33-wQEPgHaEq?pid=ImgDet&rs=1",
		Signature:       "test sign",
	}
	result := db.Model(&model.User{}).Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	zap.L().Info("create user", zap.Any("user", user))

	// todo redis缓存

	return &user, nil
}