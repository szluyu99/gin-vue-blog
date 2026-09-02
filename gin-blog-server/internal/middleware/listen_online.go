package middleware

import (
	"context"
	g "gin-blog/internal/global"
	"gin-blog/internal/handle"
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// 在线状态的有效期, 每次请求都会重新续期
const onlineExpire = 10 * time.Minute

// 监听在线状态中间件
// 登录时: 移除用户的强制下线标记
// 退出登录时: 添加用户的在线标记
func ListenOnline() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		rdb := c.MustGet(g.CTX_RDB).(*redis.Client)

		auth, err := handle.CurrentUserAuth(c)
		if err != nil {
			handle.ReturnError(c, g.ErrUserAuth, err)
			return
		}

		onlineKey := g.ONLINE_USER + strconv.Itoa(auth.ID)
		offlineKey := g.OFFLINE_USER + strconv.Itoa(auth.ID)

		// 判断当前用户是否被强制下线
		if rdb.Exists(ctx, offlineKey).Val() == 1 {
			slog.Info("用户被强制下线", "user_id", auth.ID, "username", auth.Username)
			// ReturnError 内部已经是 AbortWithStatusJSON, 不需要再 Abort 一次
			handle.ReturnError(c, g.ErrForceOffline, nil)
			return
		}

		// 每次请求都要把在线状态续期 10 分钟。
		// 已经在线的只续期, 不用重新序列化整个用户对象写一遍
		if rdb.Expire(ctx, onlineKey, onlineExpire).Val() {
			c.Next()
			return
		}
		if err := rdb.Set(ctx, onlineKey, auth, onlineExpire).Err(); err != nil {
			slog.Warn("记录在线状态失败", "user_id", auth.ID, "err", err)
		}
		c.Next()
	}
}
