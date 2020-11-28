package main

import "errors"

var (
	CAS_UNAUTHORISE_SERVICE    = errors.New("未认证授权的服务")
	CAS_LOGIN_TICKET_NOT_FOUND = errors.New("获取 Login Ticket 失败")
	CAS_BAD_REQUEST            = errors.New("抱歉！您的请求出现了异常，请稍后再试。")
	CAS_ACCOUNT_ERROR          = errors.New("用户名密码错误")
)
