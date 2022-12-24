package middleware

import (
	"gin-blog/model/dto"
	"gin-blog/utils"
	"gin-blog/utils/r"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	KEY_USER   = "user:"
	KEY_DELETE = "delete:"
)

// 监听在线状态中间件
func ListenOnline() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := utils.GetFromContext[string](c, "uuid")
		userKey := KEY_USER + uuid

		session := sessions.Default(c)
		var sessionInfo dto.SessionInfo
		// 尝试从 redis 中取用户的 SessionInfo, 取不到则从 session 中获取
		if s, err := utils.Redis.GetResult(userKey); err == nil {
			utils.Json.Unmarshal(s, &sessionInfo)
		} else if s, err := utils.Redis.GetResult(KEY_DELETE + uuid); err == nil { // 被强制离线
			utils.Json.Unmarshal(s, &sessionInfo)
		} else { // * session 中还是取不到, 就返回 401 让前端退回登录界面
			// ! 考虑一下
			val := session.Get(userKey)
			if val == nil {
				r.Send(c, http.StatusUnauthorized, r.ERROR_TOKEN_RUNTIME, nil)
				c.Abort()
				return
			}
			utils.Json.Unmarshal(val.(string), &sessionInfo)
		}

		// 判断当前是否是退出登录请求
		// if strings.Contains(c.FullPath(), "logout") {
		// 	fmt.Println("logout: ", userKey)
		// 	utils.Redis.Del(userKey) // 删除 redis 中缓存
		// 	session.Delete(userKey)  //? FIXME: 删除后 redis 还会有一条记录?
		// 	session.Save()
		// 	c.Abort()
		// 	return
		// }

		// 已经是强制下线状态
		if sessionInfo.IsOffline == 1 { // 被强制下线
			utils.Redis.Del(userKey) // 删除 redis 中缓存
			session.Delete(userKey)
			session.Save()
			// utils.Redis.Del(KEY_DELETE + uuid) // 删除 redis 中缓存
			r.Send(c, http.StatusUnauthorized, r.FORCE_OFFLINE, nil)
			c.Abort()
			return
		}

		// * 每次发送请求会更新 Redis 中的在线状态: 重新计算 10 分钟
		utils.Redis.Set(KEY_USER+uuid, utils.Json.Marshal(sessionInfo), 10*time.Minute)

		c.Next()
	}
}
