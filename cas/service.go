package cas

// Service 为接入 CAS 的服务，实现 Callback() 方法后 CAS 将在登录后执行。
type Service interface {
	GetServiceURL() string
	LoginCallback(string) error // 有些服务 CAS 登录访问后还需要二次操作，因此需要有获取登录后的回调函数，交给服务单独处理
	SetCasSession(*Session)
}
