package main

import (
	"context"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"os"
)

const (
	msgFmt = "==== %s ====\n"
)

func printKV(cfg *ini.File) {
	secs := cfg.Sections()
	for _, s := range secs {
		// 排除名为DEFAULT的section
		if s.Name() == "DEFAULT" {
			continue
		}
		fmt.Println("打印配置文件:")
		fmt.Printf("===== %s =====\n", s.Name())
		keys := s.KeyStrings()
		for _, key := range keys {
			fmt.Printf("%s:%s\n", key, s.Key(key).String())
		}
		fmt.Println()
	}
}

func main() {
	ctx := context.Background()
	// 1. 加载INI配置文件
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("无法加载配置文件: %v", err)
		return
	}
	printKV(cfg)

	// 2. 读取配置项
	// 指定section名称
	section := cfg.Section("milvus_server")
	if section == nil {
		fmt.Println("找不到指定的section")
		return
	}

	milvusAddr := section.Key("milvusAddr").String()
	username := section.Key("username").String()
	password := section.Key("password").String()

	fmt.Printf(msgFmt, "start connecting to Milvus")
	c, err := client.NewClient(ctx, client.Config{
		Address:  milvusAddr,
		Username: username,
		Password: password,
	})
	if err != nil {
		fmt.Printf("failed to connect to milvus, err: %s\n", err.Error())
		os.Exit(1)
	}
	defer c.Close()

	dbs, _ := c.ListDatabases(ctx)
	for _, db := range dbs {
		fmt.Printf(msgFmt, db)
		c.UsingDatabase(ctx, db.Name)
		colls, _ := c.ListCollections(ctx)
		var cns = make([]string, len(colls))
		// 设置隔离级别
		func1 := func(option *client.SearchQueryOption) {
			option.ConsistencyLevel = entity.ClEventually
		}
		for i := 0; i < len(colls); i++ {
			collName := colls[i].Name
			// 获取collection隔离级别
			ct, _ := c.DescribeCollection(ctx, collName)
			// 获取collection近似数量
			nums, _ := c.GetCollectionStatistics(ctx, collName)
			// 获取collection精确数量
			fieldstr := "count(*)"
			outFields := []string{fieldstr}
			rs, err := c.Query(ctx, collName, nil, "", outFields, func1)
			if err != nil {
				fmt.Printf("%s:%s\n", collName, err.Error())
				cns[i] = fmt.Sprintf("%s,ConsistencyLevel:%s,approxCount:%s,exactCount:???", collName, ct.ConsistencyLevel.CommonConsistencyLevel().String(), nums["row_count"])
				continue
			}
			column := rs.GetColumn(fieldstr)
			count, _ := column.GetAsInt64(0)
			cns[i] = fmt.Sprintf("%s,ConsistencyLevel:%s,approxCount:%s,exactCount:%d", collName, ct.ConsistencyLevel.CommonConsistencyLevel().String(), nums["row_count"], count)
		}

		for i := 0; i < len(cns); i++ {
			fmt.Printf("%d: %s\n", (i + 1), cns[i])
		}
		fmt.Println()
	}

}
