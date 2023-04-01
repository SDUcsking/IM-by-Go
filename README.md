# IM-by-Go,代码在分支master中
使用go语言的即时通讯项目,完成了基础的数据库和后端
# Gin+WebSocket项目IM(Instant Messaging 即时通信)

## 核心功能

发送和接受消息(文字,表情,图片,音频),访客

## 技术栈

前端,后端(webSocket,channel/goroutine,gin,template,gorm,sql)

## 系统架构

四层:前端,接入层,逻辑层,持久层

## 消息发送流程

A>登录>鉴权(游客)>消息类型>B

## 具体内容:

### 将请求和数据关联

1. 在main方法里面,初始化配置文件及数据库

utils.InitConfig()

utils.InitMySql()

2. 在config中app.yaml

```yaml
mysql:
	dns: root:root@tcp(127.0.0.1:3306)ginchat?charset=utf8mb4&parseTime=True&loc=Local
```

3. 新建utils包system_init.go
4. 在models包中user_basic.go里加上

```go
func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

```

5. 在service包中新建userservice.go

```go
package service

import (
	"IM/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"message": data,
	})
}

```

6. 到route包的app.go

```go
r.GET("/user/getUserList", service.GetUserList)
```

### 与我之前使用的django对比

+ gin框架明显轻量化,什么都是由自己创建,而django就非常完备,我们只需要在给定的位置填入合适的内容

+ #### 两者作为web框架是共通的

+ gin使用router放置路由,相当于django中的urls.py

+ 都是用models放置数据库模型

+ 在config中放置配置文件

+ 在utils中完成初始化

+ 在service中放置具体服务

  

  

  

## swagger

swag init创建docs目录

改造route包的app.go

再到service中加入注释

### 日志引入

在数据库初始化的时候加入自己的logger

### 完成用户模块基本功能

##### 增删改查

基本模式:

1. 在app.go中设置到路由,函数在service中
2. 在service中编写对应函数,利用models进行数据库操作
3. 在models中写出对数据库操作的函数,与service对应

#### 校验

利用govalidator进行校验

在struct字段后边加入校验规则

###### 避免重复注册:

在models中设置查询函数,在service中使用,若已经注册,则return

### 加密

利用md5和随机数加盐进行加密



### 即时通信

利用websocket进行即时通信,从redis中读取数据,注意要打开redis(先在终端cd进入redis目录"C:\Users\86152\Downloads\Redis-x64-5.0.14.1",再执行命令"redis-server redis.windows.conf")



### 初步的数据库

使用gorm作为ORM,注意为蛇形命名,即小写字母+下划线

contact记录好友关系

group_basic记录群组关系

message记录历史消息



### 发送消息

需要:发送者ID,接收者ID,消息类型,发送的内容

利用websocket发送消息,涉及channels

