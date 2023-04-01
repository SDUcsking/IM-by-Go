package service

import (
	"IM/models"
	"IM/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code":,"message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": data,
	})
}

// FindUserByNameAndPwd
// @Summary 用户登录
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code":,"message"}
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	name := c.Query("name")
	password := c.Query("password")
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户不存在",
		})
		return
	}
	flag := utils.ValidPassword(password, user.Salt, user.PassWord)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "密码不正确",
		})
		return
	}
	pwd := utils.MakePassword(password, user.Salt)

	data = models.FindUserByNameAndPwd(name, pwd)
	c.JSON(200, gin.H{
		"code":    0, //0成功  -1失败
		"message": "登陆成功",
		"data":    data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @param phone query string false "电话号码"
// @param email query string false "电子邮箱"
// @Success 200 {string} json{"code":,"message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")
	user.Phone = c.Query("phone")

	user.Email = c.Query("email")
	salt := fmt.Sprintf("%06d", rand.Int31())
	data := models.FindUserByName(user.Name)
	if data.Name != "" {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "用户名已注册",
		})
		return
	}
	if password != repassword {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "两次密码不一致",
		})
		return
	}

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "参数不匹配",
		})
		return
	}
	user.PassWord = utils.MakePassword(password, salt)
	user.Salt = salt

	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "新增用户成功",
	})

}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code":,"message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.Id = int(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "删除用户成功",
		"data":    user,
	})

}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param phone formData string false "电话"
// @param email formData string false "邮箱"
// @Success 200 {string} json{"code":,"message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.Id = int(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "参数不匹配",
		})
		return
	}
	models.UpdateUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "更新用户成功",
	})

}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {

	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(c, ws)
}

func MsgHandler(c *gin.Context, ws *websocket.Conn) {

	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)

		if err != nil {
			fmt.Println("MsgHandler发送失败", err)
		}

		tm := time.Now().Format("2006-01-01 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
