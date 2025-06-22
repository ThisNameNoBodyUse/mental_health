package service

import (
	"errors"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"mental/config"
	"mental/constant"
	"mental/dao"
	"mental/models"
	"mental/serializer"
	"mental/utils"
	"strconv"
	"time"
)

type UserService struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Username string `json:"username"`
	RoleId   string `json:"role_id"`
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
	copier.Copy(userLogin, user)

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
	// 获取令牌中的角色id列表
	_, claims, err := utils.ParseJWT(accessToken, true) // 访问令牌解析

	// 获取用户角色列表
	roles, _ := claims["roles"].([]interface{})
	flag := false // 是否是管理员
	for _, role := range roles {
		if role == "1" {
			flag = true
			break
		}
	}
	if flag == true {
		userLogin.Role = "1" // 管理员
	} else {
		userLogin.Role = "2" // 普通用户
	}
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
	if userService.Username == "" {
		return false, errors.New("用户名不能为空")
	}
	if userService.RoleId == "" {
		return false, errors.New("角色不能为空")
	}

	// 转换并验证 roleId
	roleId, err := strconv.Atoi(userService.RoleId)
	if err != nil || (roleId != 1 && roleId != 2) {
		return false, errors.New("角色ID非法")
	}

	// 分布式锁
	key := constant.RegisterPrefix + userService.Account
	lock, err := utils.TryLock(key, 5*time.Second)
	if err != nil {
		return false, errors.New("该账号正在注册中，请勿重复操作")
	}
	defer utils.Unlock(lock)

	// 事务处理
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		userDao := dao.NewUserDao(tx)
		existingUser, err := userDao.GetUserByAccount(userService.Account)
		if err != nil {
			return err
		}
		if existingUser != nil {
			return errors.New("该账号已被注册")
		}

		// 创建用户
		hashedPwd, err := utils.HashPassword(userService.Password)
		if err != nil {
			return err
		}
		// 雪花算法生成用户唯一id
		snowflake, _ := utils.NewSnowflake()
		id := snowflake.GenerateID()

		newUser := &models.User{
			Account:  userService.Account,
			Password: hashedPwd,
			Username: userService.Username,
		}
		newUser.Id = int(id)
		if err := tx.Create(newUser).Error; err != nil {
			return err
		}

		// 创建角色绑定
		userRole := &models.UserRole{
			UserID: newUser.Id,
			RoleID: roleId,
		}
		if err := tx.Create(userRole).Error; err != nil {
			return err
		}

		return nil // 提交事务
	})

	// 判断事务是否成功
	if err != nil {
		return false, err
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

// UserLogout 退出登录
func (userService *UserService) UserLogout(token string) error {
	_, claims, _ := utils.ParseJWT(token, true)
	accessJti := claims["jti"].(string) // 访问令牌的jti
	// 将访问令牌jti存入 Redis黑名单
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	ttl := expirationTime.Sub(time.Now())
	accessKey := constant.BlackListPrefix + accessJti // 访问令牌在redis中的key
	err := utils.Set(accessKey, "true", ttl)          // 访问令牌存入redis
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

// UpdateAvatar 根据用户id修改头像
func (userService *UserService) UpdateAvatar(userId int64, url string) error {
	dao := dao.NewUserDao(config.DB)
	err := dao.UpdateAvatar(userId, url)
	return err
}

// UpdateUsernameOrEmail 根据用户id修改用户名/邮箱
func (userService *UserService) UpdateUsernameOrEmail(id int64, username string, email string) error {
	dao := dao.NewUserDao(config.DB)
	err := dao.UpdateUsernameOrEmail(id, username, email)
	return err
}

// UpdatePassword 尝试根据用户id、原密码，更新密码
func (s *UserService) UpdatePassword(userId int64, oldPwd string, newPwd string) error {
	// 查询用户
	dao := dao.NewUserDao(config.DB)
	user, err := dao.GetUserById(userId)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 校验旧密码
	if !utils.CheckPassword(user.Password, oldPwd) {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashedPwd, err := utils.HashPassword(newPwd)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	return dao.UpdatePassword(userId, hashedPwd)
}
