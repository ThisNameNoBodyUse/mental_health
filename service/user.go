package service

import (
	"errors"
	"github.com/jinzhu/copier"
	"mental/config"
	"mental/constant"
	"mental/dao"
	"mental/models"
	"mental/serializer"
	"mental/utils"
	"time"
)

type UserService struct {
	Account  string `json: "account"`
	Password string `json: "password"`
}

// UserLogin 登录，返回登录信息 + err
func (userService *UserService) UserLogin() (*serializer.UserLogin, error) {
	if userService.Account == "" {
		return nil, errors.New("账号不能为空")
	}
	if userService.Password == "" {
		return nil, errors.New("密码不能为空")
	}
	// 调用dao层查询用户是否存在
	userDao := dao.NewUserDao(config.DB)
	// 根据账号查询用户
	user, err := userDao.GetUserByAccount(userService.Account)
	if user == nil { // 用户不存在
		return nil, errors.New("账号或密码错误")
	}
	// 将加密后的密码和数据库中的比对
	valid := utils.CheckPassword(user.Password, userService.Password) // 加密密码，原文密码
	if !valid {                                                       // 密码错误
		return nil, errors.New("账号或密码错误")
	}

	// 比对成功，进行登录，申请访问令牌和刷新令牌，封装到UserLogin中返回
	userLogin := new(serializer.UserLogin)
	copier.Copy(userLogin, user) // 后面拷贝到前面

	// 生成双令牌
	accessToken, err := utils.GenerateJWT(user, true)
	if err != nil {
		return nil, errors.New("生成访问令牌失败")
	}
	refreshToken, err := utils.GenerateJWT(user, false)
	if err != nil {
		return nil, errors.New("生成刷新令牌失败")
	}
	userLogin.Authentication = accessToken
	userLogin.RefreshToken = refreshToken
	return userLogin, nil
}

// UserRegister 用户注册，判断是否能成功注册
func (userService *UserService) UserRegister() (bool, error) {
	if userService.Account == "" {
		return false, errors.New("账号不能为空")
	}
	if userService.Password == "" {
		return false, errors.New("密码不能为空")
	}
	// 引入分布式锁
	key := constant.RegisterPrefix + userService.Account // 注册操作分布式锁的key
	lock, err := utils.TryLock(key, 5*time.Second)       // 尝试上锁5s,获取锁对象
	if err != nil {
		return false, errors.New("该账号正在注册中，请勿重复操作")
	}
	defer utils.Unlock(lock) // 任务完成后释放锁 defer 后的语句在该函数结束后执行

	userDao := dao.NewUserDao(config.DB)
	user, err := userDao.GetUserByAccount(userService.Account)
	// 如果用户不为空
	if user != nil {
		// 说明该账号已被注册
		return false, errors.New("该账号已被注册")
	}
	// 账号没被注册过，可以进行注册
	user = new(models.User)
	user.Account = userService.Account
	password, err := utils.HashPassword(userService.Password)
	if err != nil {
		return false, err
	}
	user.Password = password
	// 插入数据库
	save := userDao.Save(user)
	if save.Error != nil {
		return false, save.Error
	}
	return true, nil
}

// GetUserInfoById 根据id查询用户基本信息
func (userService *UserService) GetUserInfoById(id int64) (*serializer.UserInfo, error) {
	userDao := dao.NewUserDao(config.DB)
	user, err := userDao.GetUserById(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	userInfo := new(serializer.UserInfo)
	copier.Copy(userInfo, user)
	return userInfo, err
}

// AdminLogout 退出登录
func (userService *UserService) AdminLogout(token string) error {
	_, claims, _ := utils.ParseJWT(token, true)
	access_jti := claims["jti"].(string) // 访问令牌的jti
	// 将访问令牌的 jti 存入 Redis 黑名单，并设置过期时间
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	ttl := expirationTime.Sub(time.Now())
	access_key := constant.BlackListPrefix + access_jti // 访问令牌在redis中的key
	err := utils.Set(access_key, "true", ttl)           // 访问令牌存入redis
	if err != nil {
		return errors.New("访问令牌拉黑异常")
	}
	return nil
}

// RefreshToken 根据刷新令牌，生成新的访问令牌并返回
func (userService *UserService) RefreshToken(refreshToken string) (string, error) {
	_, claims, err := utils.ParseJWT(refreshToken, false) // 尝试解析刷新令牌
	if err != nil {
		return "", errors.New("invalid refresh token")
	}
	var user = new(models.User)
	// 获取刷新令牌中用户基本信息
	id := claims["id"].(float64)
	account := claims["account"].(string)
	username := claims["username"].(string)
	// 基本信息设置到user结构体中
	user.Id = int(id)
	user.Account = account
	user.Username = username
	// 生成新的访问令牌
	newAccessToken, err := utils.GenerateJWT(user, true) // 生成访问令牌
	if err != nil {
		return "", errors.New("访问令牌生成失败")
	}
	return newAccessToken, nil
}
