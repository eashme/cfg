package coord_cfg

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
	"github.com/xormplus/xorm/names"
	"log"
	"time"
)

type ConfigStatus int8

type Config struct {
	Code    string `xorm:"varchar(64) pk" remark:"配置编码"`
	Name    string `xorm:"varchar(32)" remark:"配置名称"`
	Value   string `xorm:"text" remark:"配置值"`
	Remark  string `xorm:"text" remark:"配置备注"`
	Version int64  `xorm:"version index"`
	Created int64  `xorm:"created" remark:"创建时间"`
	Updated int64  `xorm:"updated" remark:"更新时间"`
	Deleted int64  `xorm:"deleted default(0)" remark:"删除时间"`
}

var (
	eg xorm.EngineInterface
)

func connDB(user, pwd, host string, port uint16, db string) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, pwd, host, port, db)
	eg, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		return err
	}
	// 设置 下划线链接命名规则
	eg.SetMapper(&names.SnakeMapper{})
	// 设置最大数据库空闲连接数
	eg.SetMaxIdleConns(200)
	// 设置最大数据库连接数
	eg.SetMaxOpenConns(200)
	// 设置时区为中国东八区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	eg.SetTZLocation(loc)
	eg.SetTZDatabase(loc)
	err = eg.Ping()
	if err != nil {
		log.Printf("failed ping databases %s  err:%v", dsn, err)
		return err
	}
	return nil
}

func getFromDB(ctx context.Context, k string) string {
	ses := eg.NewSession().Context(ctx)
	defer ses.Close()

	conf := &Config{
		Code: k,
	}
	_, err := ses.MustCols("code", "status").Get(conf)
	if err != nil {
		log.Printf("failed get cfg [%s] err: %v ", k, err)
		return ""
	}
	return conf.Value
}

func set2DB(ctx context.Context, k string, v string, remark ...string) error {
	ses := eg.NewSession().Context(ctx)
	defer ses.Close()

	bean := &Config{
		Code: k,
	}
	ok, err := ses.Get(bean)
	if err != nil {
		return err
	}
	var (
		desc string
		n    int64
	)
	if len(remark) > 0 {
		desc = remark[0]
	}

	if ok {
		n, err = ses.Update(&Config{
			Code:    k,
			Value:   v,
			Remark:  desc,
			Version: bean.Version,
		})
	} else {
		n, err = ses.InsertOne(&Config{
			Code:   k,
			Value:  v,
			Remark: desc,
		})
	}
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("no affected ")
	}
	return nil
}
