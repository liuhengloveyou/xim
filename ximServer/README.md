#xim

Im by GO

##接口

###发消息
	POST /send?token="认证串或者放在cookie里"
	Body: {
		type: txt|icon|img|sound
		message: xxx xxx
	}
	
	成功: 200
	失败: ^200 {"message":"错误信息."}
	
###接收消息

	GET /recv?token="认证串或者放在cookie里"

	成功: 200 {"message":"消息."}
	失败: ^200 {"message":"错误信息."}
	
###拉取好友列表
	GET /friends/list?v="客户端当前好友列表版本"&token="认证串或者放在cookie里"
	
	成功: 200 ["好友userid", ...]
	失败: ^200 {"message":"错误信息."}	

###添加用户
	POST /user/add
	Body:
	{
		"cellphone":"18510511015", 
		"email":"liuhengloveyou@gmail.com",
		"nickname":"恒恒 ",
		"password":"123456"
	}

###用户登录

	POST /user/login
	Body:{
		"cellphone":"15236379552", 
		"email":"liuhengloveyou@gmail.com",
		"nickname":"恒恒",
		"password":"123456"
	}		

	

##数据库
	
	CREATE DATABASE `xim` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_bin */;

	CREATE TABLE `xim`.`user` (
	  `id` VARCHAR(45) NOT NULL,
	  `add_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	  `version` int(11) DEFAULT NULL,
	  PRIMARY KEY (`id`),
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
	
