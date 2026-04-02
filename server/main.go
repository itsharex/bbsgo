package main

import (
	"bbsgo/cache"
	"bbsgo/database"
	"bbsgo/routes"
	"bbsgo/utils"
	"log"
	"net/http"
)

func main() {
	database.InitDB()
	database.AutoMigrate()
	cache.Init()

	seedData()

	r := routes.SetupRoutes()

	log.Printf("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func seedData() {
	var forumCount int64
	database.DB.Model(&struct{}{}).Table("forums").Count(&forumCount)
	if forumCount == 0 {
		forums := []struct {
			Name        string
			Description string
			SortOrder   int
		}{
			{"全部", "默认首页，显示所有板块的帖子", 1},
			{"技术交流", "编程语言、框架、架构等纯技术讨论", 2},
			{"提问求助", "发帖求助、解答问题", 3},
			{"业界资讯", "科技新闻、技术动态、行业趋势", 4},
			{"资源分享", "工具、教程、电子书、代码片段", 5},
			{"求职招聘", "内推、招聘信息、面经", 6},
			{"灌水闲聊", "生活、娱乐、非技术话题", 7},
			{"站务管理", "公告、反馈、版务", 8},
		}

		for _, f := range forums {
			database.DB.Exec("INSERT INTO forums (name, description, sort_order, allow_post, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
				f.Name, f.Description, f.SortOrder, true)
		}
		log.Println("Default forums created")
	}

	var configCount int64
	database.DB.Model(&struct{}{}).Table("site_configs").Count(&configCount)
	if configCount == 0 {
		configs := []struct {
			Key   string
			Value string
		}{
			{"site_name", "彩虹BBS"},
			{"site_logo", ""},
			{"site_icon", ""},
			{"site_description", "一个现代化的社区论坛系统"},
			{"email_enabled", "false"},
			{"email_host", ""},
			{"email_port", "465"},
			{"email_user", ""},
			{"email_password", ""},
			{"email_from", ""},
			{"email_from_name", "彩虹BBS"},
			{"qiniu_access_key", ""},
			{"qiniu_secret_key", ""},
			{"qiniu_bucket", ""},
			{"qiniu_domain", ""},
			{"jwt_secret", "change-this-secret-in-production"},
			{"jwt_expire_days", "7"},
			{"cache_num_counters", "10000"},
			{"cache_max_cost", "10000000"},
		}
		for _, c := range configs {
			database.DB.Exec("INSERT INTO site_configs (key, value, created_at, updated_at) VALUES (?, ?, datetime('now'), datetime('now'))",
				c.Key, c.Value)
		}
		log.Println("Default site configs created")
	}

	var tagCount int64
	database.DB.Model(&struct{}{}).Table("tags").Count(&tagCount)
	if tagCount == 0 {
		tags := []struct {
			Name        string
			Icon        string
			Description string
			SortOrder   int
			IsOfficial  bool
		}{
			{"今日份松弛", "😌", "分享慢生活、拒绝焦虑的瞬间", 1, true},
			{"爱你老己", "💖", "对自己好的方式、自我关怀", 2, true},
			{"活人感日常", "🫠", "真实、不完美的生活碎片", 3, true},
			{"邪修一下", "⚡", "找捷径、高效摆烂、反内卷", 4, true},
			{"外耗模式", "😤", "与其内耗自己，不如外耗别人", 5, true},
			{"今日小确幸", "✨", "微小而确定的幸福瞬间", 6, true},
			{"我的互联网嘴替", "🗣️", "说出了我想说但说不出的话", 7, true},
			{"求建议/避雷", "❓", "生活求助、消费避坑", 8, true},
			{"笑死我了", "😂", "搞笑段子、神评论、趣图", 9, true},
			{"真香现场", "🔥", "打脸时刻、意外真香", 10, true},
			{"破防了", "💔", "感动、扎心、被戳中的瞬间", 11, true},
			{"什么水平？", "🤔", "求评价、求鉴定、秀成果", 12, true},
		}
		for _, t := range tags {
			database.DB.Exec("INSERT INTO tags (name, icon, description, sort_order, usage_count, is_official, is_banned, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))",
				t.Name, t.Icon, t.Description, t.SortOrder, 0, t.IsOfficial, false)
		}
		log.Println("Default tags created")
	}

	var userCount int64
	database.DB.Model(&struct{}{}).Table("users").Count(&userCount)
	if userCount == 0 {
		// set auto increment start
		database.DB.Exec("INSERT INTO sqlite_sequence (name, seq) VALUES ('users', 9999)")
		adminPassword, _ := utils.HashPassword("12345678")
		database.DB.Exec(`INSERT INTO users (username, email, nickname, password_hash, role, credits, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))`,
			"admin", "admin@example.com", "管理员", adminPassword, 2, 10000)
		log.Println("Admin user created (username: admin, password: 12345678)")

		users := []struct {
			Username string
			Email    string
			Nickname string
		}{
			{"testuser1", "test1@example.com", "测试用户1"},
			{"testuser2", "test2@example.com", "测试用户2"},
			{"testuser3", "test3@example.com", "测试用户3"},
			{"testuser4", "test4@example.com", "测试用户4"},
			{"testuser5", "test5@example.com", "测试用户5"},
			{"testuser6", "test6@example.com", "测试用户6"},
			{"testuser7", "test7@example.com", "测试用户7"},
			{"testuser8", "test8@example.com", "测试用户8"},
			{"testuser9", "test9@example.com", "测试用户9"},
			{"testuser10", "test10@example.com", "测试用户10"},
		}

		for _, u := range users {
			hashedPassword, _ := utils.HashPassword("123456")
			database.DB.Exec(`INSERT INTO users (username, email, nickname, password_hash, credits, created_at, updated_at) 
				VALUES (?, ?, ?, ?, ?, datetime('now'), datetime('now'))`,
				u.Username, u.Email, u.Nickname, hashedPassword, 1000)
		}
		log.Println("Test users created")
	}

	var topicCount int64
	database.DB.Model(&struct{}{}).Table("topics").Count(&topicCount)
	if topicCount == 0 {
		// set auto increment start
		database.DB.Exec("INSERT INTO sqlite_sequence (name, seq) VALUES ('topics', 9999)")
		topics := []struct {
			Title      string
			Content    string
			UserID     uint
			ForumID    uint
			LikeCount  int
			ReplyCount int
			ViewCount  int
		}{
			{"bbs-go v3.5.0 发布，升级 go1.18", "文档地址：https://docs.bbs-go.com/\n官网交流：https://mlog.club\n问题反馈：https://mlog.club/topic/node/3\n\n本次更新内容：\n1. 升级 Go 1.18 版本\n2. 优化数据库查询性能\n3. 修复已知 bug\n4. 新增配置管理功能", 1, 4, 12, 8, 352},
			{"Vue3 + TypeScript 项目实践分享", "最近用 Vue3 + TypeScript 做了一个项目，分享一些实践经验：\n\n1. 组合式 API 真的很香，逻辑复用更方便了\n2. TypeScript 的类型推导需要好好配置\n3. Pinia 比 Vuex 更简洁好用\n\n有问题的朋友欢迎留言讨论~", 2, 4, 45, 18, 892},
			{"今天天气不错，适合摸鱼", "周末到了，阳光明媚，正是摸鱼好时节。大家最近都在看什么书？有什么好剧推荐吗？\n\n我最近在看《三体》，真的很精彩！强烈推荐给还没看过的朋友。", 3, 9, 32, 21, 687},
			{"求助：MySQL 慢查询优化", "公司有个 MySQL 表数据量大概 500 万，查询越来越慢了。\n\n表结构大概是：\n- id (主键)\n- user_id (索引)\n- created_at (索引)\n- content (text)\n\n查询语句：SELECT * FROM table WHERE user_id = ? ORDER BY created_at DESC LIMIT 20\n\n请问有什么优化建议吗？", 4, 5, 8, 12, 234},
			{"分享一些 Linux 常用命令", "整理了一些常用的 Linux 命令，希望对大家有帮助：\n\n查看端口占用情况：\nnetstat -tunlp | grep 端口号\n\n通过 ssh 将远程端口映射到本地端口：\nssh -L 13306:127.0.0.1:3306 用户名@远程地址 -N\n\n这样远程服务器就不需要开放需要的端口到公网了，更安全。", 5, 7, 18, 5, 221},
			{"C++ 程序返回 value 3221226356 求教！", "return value 3221226356 求教求教！\n\n#include <iostream>\nusing namespace std;\nint main() {\n    int n; \n    double *p=new double[n]; \n    cin>>n;\n    for(int i=0;i<n;i++) { cin>>p[i]; }\n    for(int i=0;i<n;i++) { cout<<p[i]<<\" \"; }\n    return 0;\n}\n\n程序运行时出现这个错误，请问是什么原因？", 6, 5, 3, 7, 126},
			{"分享一张今天拍的美照", "今天去公园玩了，随手拍了一张照片，分享给大家~\n\n[图片]\n\n摄影器材：Sony A7M3\n参数：f/2.8, 1/500s, ISO100", 7, 10, 156, 43, 2341},
			{"网站有个 BUG 反馈", "在使用网站时发现一个问题：\n\n当我在移动端浏览帖子时，点击回复按钮后键盘会遮挡输入框，需要手动收起键盘才能看到输入内容。\n\n浏览器：Safari\n系统：iOS 16\n设备：iPhone 13 Pro\n\n希望能修复一下，谢谢！", 8, 11, 5, 3, 89},
			{"推荐一个很好用的开源项目", "最近发现一个很棒的开源项目：\n\n项目名称：VSCode\nGitHub 地址：https://github.com/microsoft/vscode\n\n功能强大，插件生态丰富，支持几乎所有编程语言。强烈推荐给各位开发者！\n\n大家还有什么好用的工具欢迎分享~", 9, 7, 67, 29, 1523},
			{"2024 年前端技术趋势预测", "随着 AI 的快速发展，前端领域也在不断变化。以下是我对 2024 年前端技术趋势的一些预测：\n\n1. AI 辅助开发将成为标配\n2. Server Components 会更加流行\n3. TypeScript 使用率继续上升\n4. Rust 在前端工具链中的应用会更广泛\n5. Web Components 可能会迎来第二春\n\n大家怎么看？欢迎讨论！", 10, 6, 89, 32, 1523},
		}

		for _, t := range topics {
			database.DB.Exec(`INSERT INTO topics (title, content, user_id, forum_id, like_count, reply_count, view_count, created_at, updated_at) 
				VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now', '-'||abs(random())%72||' hours'), datetime('now', '-'||abs(random())%72||' hours'))`,
				t.Title, t.Content, t.UserID, t.ForumID, t.LikeCount, t.ReplyCount, t.ViewCount)
		}
		log.Println("Test topics created")

	}
}
