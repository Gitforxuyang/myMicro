package server

//定义server的基本方法

type Server interface {
	//获取配置对象
	Options() Options
	//init
	Init(...Option)
}

type Option func(*Options)