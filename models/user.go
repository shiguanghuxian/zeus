package models

import (
	"encoding/json"
	"errors"

	"53it.net/zues/internal"
)

type User struct {
	Id       int    `json:"id" xorm:"int pk autoincr 'id'"`
	Name     string `json:"name" xorm:"VARCHAR(30)"`
	Username string `json:"username" xorm:"varchar(60) notnull unique 'username'"`
	Password string `json:"password" xorm:"varchar(40) notnull 'password'"`
	Phone    string `json:"phone" xorm:"varchar(25) notnull 'phone'"`
	Email    string `json:"email" xorm:"varchar(60) notnull 'email'"`
	Icon     string `json:"icon" xorm:"varchar(255) notnull 'icon'"`
	Uptime   int64  `json:"uptime" xorm:"int 'uptime'"`
	State    string `json:"state" xorm:"default '1' index ENUM('0','1')"`
	NameEn   string `json:"name_en" xorm:"VARCHAR(60)"`
	Sex      string `json:"sex" xorm:"default '0' CHAR(1)"`
	Addtime  int    `json:"addtime" xorm:"INT(11)"`
	GroupId  int    `json:"group_id" xorm:"default 0 index INT(11)"`
	Info     string `json:"info" xorm:"TEXT"`
	Level    int    `json:"level" xorm:"index TINYINT(4)"`
	IsDelete string `json:"is_delete" xorm:"default '0' index ENUM('0','1')"`
	Language string `json:"language" xorm:"default 'zh-CN' varchar(60)"`
}

// 实际表名
func (this *User) TableName() string {
	return "zn_user"
}

func (this *User) String() string {
	s, _ := json.Marshal(this)
	return string(s)
}

// 验证登录
func VerifyLogin(username, password string) (*User, error) {
	user := new(User)
	user.Username = username
	user.Password = internal.UserPwdEncrypt(password, "")
	isexist, err := engine.Get(user)
	if isexist == true && err == nil {
		user.Password = ""
		return user, nil
	} else {
		internal.LogFile.E("用户名或密码错误:" + err.Error())
		return new(User), errors.New("用户名或密码错误")
	}
}

// 根据用户id查询用户信息
func GetIdUserInfo(id int) (*User, error) {
	user := new(User)
	user.Id = id
	isexist, err := engine.Get(user)
	if isexist == true && err == nil {
		return user, nil
	} else {
		internal.LogFile.E("未查询到信息:" + err.Error())
		return new(User), errors.New("未查询到信息")
	}
}

// 根据用户id修改用户信息
func UpdateIdUserInfo(id int, u *User) error {
	if u.Password != "" {
		u.Password = internal.UserPwdEncrypt(u.Password, "")
	}
	affected, err := engine.Id(id).Update(u)
	if affected > 0 && err == nil {
		return nil
	} else {
		internal.LogFile.E("更新用户信息失败:" + err.Error())
		return errors.New("更新用户信息失败")
	}
}

// 根据条件查询总数
func GetUsersWhereCount(where string) (c int64, err error) {
	user := new(User)
	if where != "" {
		c, err = engine.Where(where).Count(user)
	} else {
		c, err = engine.Count(user)
	}
	if err != nil {
		internal.LogFile.E("查询用户数量失败：" + err.Error())
	}
	return c, err
}

// 根据条件查询用户列表
func GetUsersWhereList(page, pageCount int, where string) ([]User, error) {
	pageStart := (page - 1) * pageCount
	users := make([]User, 0)
	err := engine.Where(where).Asc("id").Limit(pageCount, pageStart).Find(&users)
	if err != nil {
		internal.LogFile.W("查询用户列表数据失败：" + err.Error())
	}
	return users, err
}
