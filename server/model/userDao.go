package model

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// UserDao 实例，全局唯一
var CurrentUserDao *UserDao

type UserDao struct {
	pool *redis.Pool
}

// 初始化一个 UserDao 结构体示例，
func InitUserDao(pool *redis.Pool) (currentUserDao *UserDao) {
	currentUserDao = &UserDao{pool: pool}
	return
}

func idIncr(conn redis.Conn) (id int, err error) {
	res, err := conn.Do("incr", "users_id")
	id = int(res.(int64))
	if err != nil {
		fmt.Printf("id自增错误: %v\n", err)
		return
	}
	return
}

// 根据用户 username 获取用户信息
// 获取成功返回 user 信息，err nil
// 获取失败返回 err，user 为 nil
func (ud *UserDao) GetUserByUsername(username string) (user User, err error) {
	conn := ud.pool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("hget", "users", username))
	if err != nil {
		err = errors.New("用户不存在！")
		return
	}
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Printf("GetUserByUsername方法反序列化错误: %v\n", err)
		return
	}
	return
}

// 注册用户
// 用户名不能重复
func (ud *UserDao) Register(username, password, passwordConfirm string) (user User, err error) {
	// 判断密码是否正确
	if password != passwordConfirm {
		err = errors.New("两次密码不一致！")
		return
	}

	// 保证用户名不重复
	user, err = ud.GetUserByUsername(username)
	if err == nil {
		fmt.Printf("用户已存在!\n")
		err = errors.New("用户已存在！")
		return
	}

	conn := ud.pool.Get()
	defer conn.Close()
	// id 自增 1，作为下个用户 id
	id, err := idIncr(conn)
	if err != nil {
		return
	}

	user = User{ID: id, Name: username, Password: password}
	info, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Register序列化错误: %v", err)
		return
	}
	_, err = conn.Do("hset", "users", username, info)
	if err != nil {
		fmt.Printf("存储用户信息到redis错误: %v", err)
		return
	}
	return
}

func (ud *UserDao) Login(username, password string) (user User, err error) {
	user, err = ud.GetUserByUsername(username)
	if err != nil {
		fmt.Printf("GetUserByUsername错误！: %v\n", err)
		return
	}

	if user.Password != password {
		err = errors.New("密码错误！")
		return
	}

	return
}
