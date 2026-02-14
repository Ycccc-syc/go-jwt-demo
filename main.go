package main

import (
	"go-jwt/database"
	"go-jwt/models"
	"go-jwt/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	// 初始化数据库
	database.InitDB()

	r := gin.Default()
	r.Use(GlobalRecovery())
	// r.GET("/test-panic", func(c *gin.Context) {
	// 	// 故意触发 panic
	// 	a := 0
	// 	b := 1 / a
	// 	utils.Success(c, gin.H{"result": b})
	// })
	// 登录接口（从数据库验证）
	r.POST("/login", func(c *gin.Context) {
		var req LoginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, "用户名或密码不能为空")
			return
		}

		// 从数据库查用户
		var user models.User
		if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
			utils.Error(c, http.StatusUnauthorized, "账号或密码错误")
			return
		}

		// 校验密码
		if !user.CheckPassword(req.Password) {
			utils.Error(c, http.StatusUnauthorized, "账号或密码错误")
			return
		}

		// 生成 token
		token, _ := utils.GenerateToken(user.ID)
		utils.Success(c, gin.H{
			"msg":   "登录成功",
			"token": token,
		})
	})

	// 需要登录才能访问的接口
	auth := r.Group("/")
	auth.Use(JWTMiddleware())
	{
		auth.GET("/user/info", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			utils.Success(c, gin.H{
				"msg":    "已登陆",
				"userID": userID,
			})
		})
	}

	r.Run(":8080")
}

// JWT 登录校验中间件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			utils.Error(c, http.StatusUnauthorized, "未提供token")
			c.Abort()
			return
		}

		// 去掉 Bearer 前缀
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		}

		// 解析 token
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			utils.Error(c, http.StatusUnauthorized, "token已无效或已过期")
			c.Abort()
			return
		}

		// 把用户ID存到上下文
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// GlobalRecovery 全局异常处理中间件
func GlobalRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		// defer 确保即使 panic 也能执行到这里
		defer func() {
			if err := recover(); err != nil {
				// 打印错误信息到日志，方便排查
				log.Printf("服务发生 panic: %v", err)
				// 返回统一的错误响应
				utils.Error(c, http.StatusInternalServerError, "服务器内部错误")
				// 终止后续处理
				c.Abort()
			}
		}()
		// 继续执行后续的中间件和 handler
		c.Next()
	}
}
