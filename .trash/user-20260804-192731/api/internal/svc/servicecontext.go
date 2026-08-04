package svc

import (
	"context"
	"fmt"

	"tjxt/pkg/xerr"
	"tjxt/apps/user/api/internal/config"
	"tjxt/apps/user/api/internal/model"
	"tjxt/apps/user/api/internal/service"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/syncx"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	Config config.Config

	Conn  sqlx.SqlConn

	UserModel   *model.UserModel
	DetailModel *model.UserDetailModel
	AuthModel   *model.AuthModel

	Redis *redis.Redis

	AuthService service.AuthService
	UserService service.UserService
}

// NewServiceContext 创建服务上下文，初始化数据库、Redis、模型
func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DB.DataSource)

	// 初始化 Redis（可选：留空则不使用缓存）
	var (
		r      *redis.Redis
		cCache cache.Cache
	)
	if c.Cache.Host != "" && c.Cache.Port > 0 {
		addr := fmt.Sprintf("%s:%d", c.Cache.Host, c.Cache.Port)
		conf := redis.RedisConf{Host: addr}
		rdb, err := redis.NewRedis(conf)
		if err == nil {
			r = rdb
			cCache = cache.NewNode(rdb, syncx.NewSingleFlight(), cache.NewStat("user.cache"), sqlx.ErrNotFound)
		}
	}

	userModel := model.NewUserModel(conn, cCache)
	detailModel := model.NewUserDetailModel(conn, cCache)
	authModel := model.NewAuthModel(conn, "tj_auth.login_record", c.RoleTable)

	return &ServiceContext{
		Config:      c,
		Conn:        conn,
		UserModel:   userModel,
		DetailModel: detailModel,
		AuthModel:   authModel,
		Redis:       r,
		AuthService: service.NewAuthService(userModel, authModel, c.Auth.AccessSecret, c.Auth.AccessExpire),
		UserService: service.NewUserService(userModel, detailModel, authModel, c.DefaultPassword),
	}
}

// CheckCtx 检查上下文，返回一个错误
func CheckCtx(ctx context.Context, checkFn func() error) error {
	if err := checkFn(); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "检查失败")
	}
	return nil
}