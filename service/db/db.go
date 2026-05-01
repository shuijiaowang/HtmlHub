package db

import (
	"context"
	"errors"
	"fmt"
	"htmlhub/config"
	"htmlhub/model"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// CustomLogger 自定义 GORM 日志器，用于屏蔽 ErrRecordNotFound 的日志
type CustomLogger struct {
	logger.Interface // 嵌入默认的日志接口
}

// Trace 重写 Trace 方法，拦截 SQL 执行日志
func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// 如果错误是 RecordNotFound，则忽略该错误，不打印日志
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果你希望连 SQL 语句也不打印，直接 return 即可
		// return

		// 如果你希望打印 SQL 但不打印错误，就把 err 置为 nil 传给下层
		l.Interface.Trace(ctx, begin, fc, nil)
		return
	}
	// 其他情况（正常查询或其他错误），走默认日志逻辑
	l.Interface.Trace(ctx, begin, fc, err)
}

func InitDB() error {
	mysqlCfg := config.AppConfig.MySQL
	// 拼接 MySQL DSN（数据源名称）：user:password@tcp(host:port)/dbname?config
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		mysqlCfg.Username, // 用户名
		mysqlCfg.Password, // 密码
		mysqlCfg.Path,     // 数据库地址（host）
		mysqlCfg.Port,     // 端口
		mysqlCfg.DBName,   // 数据库名
		mysqlCfg.Config,   // 连接参数（字符集、时区等）
	)
	// 1. 配置基础日志策略（保留原有的日志级别，如 Info/Warn）
	// 你可以根据需要调整 logger.Info / logger.Warn / logger.Error
	defaultLogger := logger.Default.LogMode(logger.Error)
	// 2. 使用我们自定义的日志器包装一层
	customLogger := &CustomLogger{
		Interface: defaultLogger,
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		Logger: customLogger, // 3. 注入自定义日志器
	})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %w", err) // 包装错误，便于上层排查
	}

	err = DB.AutoMigrate(

		&model.User{},
		&model.HtmlRecord{},
		&model.HtmlRecordData{},

		// 其他模型...
	)
	if err != nil {
		return fmt.Errorf("模型迁移失败: %w", err)
	}
	//if err := model.InitDefaultConfig(DB); err != nil {
	//	return err
	//}
	//// （可选）优化数据库连接池（推荐配置）
	//sqlDB, err := DB.DB()
	//if err != nil {
	//	return fmt.Errorf("获取数据库连接池失败: %w", err)
	//}
	//sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	//sqlDB.SetMaxOpenConns(50)           // 最大打开连接数
	//sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大存活时间

	log.Println("数据库初始化成功")
	return nil
}
