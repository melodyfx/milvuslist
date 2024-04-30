package main

import (
	"context"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
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

	root := &Root{}
	dbs, _ := c.ListDatabases(ctx)

	for _, db := range dbs {
		d := NewDatabase(db.Name)
		root.Add(d)
		c.UsingDatabase(ctx, db.Name)
		colls, _ := c.ListCollections(ctx)
		for i := 0; i < len(colls); i++ {
			collName := colls[i].Name
			co := NewCollection(collName)
			co.id = i + 1
			d.Add(co)

			ct, _ := c.DescribeCollection(ctx, collName)
			// 处理collection
			co.consistency_level = ct.ConsistencyLevel.CommonConsistencyLevel().String()
			stats, _ := c.GetLoadState(ctx, collName, nil)
			co.loadStat = stats

			// 获取collection近似数量
			nums, _ := c.GetCollectionStatistics(ctx, collName)
			co.nums = nums["row_count"]
			co.shardNum = ct.ShardNum

			// 处理schema
			schema := NewSchema()
			co.Add(schema)
			schema.desc = ct.Schema.Description
			schema.dynamicField = ct.Schema.EnableDynamicField
			// 处理fields
			fields := ct.Schema.Fields
			for _, field := range fields {
				f := NewField()
				f.name = field.Name
				f.ftype = field.DataType.Name()
				f.primaryKey = field.PrimaryKey
				f.autoId = field.AutoID
				f.typeParam = field.TypeParams
				idxs, _ := c.DescribeIndex(ctx, collName, field.Name)
				f.indexs = idxs
				schema.Add(f)
			}
		}
	}
	root.show()
}
