package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

/*
样例内容数据: 分类、标签、文章、评论、留言、友链

和 seed_resource.go / generate-data 里的系统数据分开:
系统数据(菜单/资源/角色/配置)是跑起来就必须有的, 生产环境也要灌;
这里的内容纯粹是本地测试用的 —— 空库能验证"服务起得来", 但验证不了列表分页、
归档按月分组、评论回复、标签云这些必须有数据才看得出问题的地方。

只在库里一篇文章都没有时才灌, 避免重复执行灌出一堆重复内容。
*/

// 文章数量刻意超过前台每页条数(9), 这样一进列表页就能看到分页
const demoArticleCount = 15

func SeedDemoContent(db *gorm.DB) error {
	var count int64
	if err := db.Model(&Article{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	var admin UserAuth
	if err := db.Where("username", "admin").First(&admin).Error; err != nil {
		return fmt.Errorf("找不到 admin 用户, 先灌系统数据: %w", err)
	}
	var guest UserAuth
	if err := db.Where("username", "guest").First(&guest).Error; err != nil {
		return fmt.Errorf("找不到 guest 用户, 先灌系统数据: %w", err)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		categories, err := seedDemoCategories(tx)
		if err != nil {
			return err
		}
		tags, err := seedDemoTags(tx)
		if err != nil {
			return err
		}
		articles, err := seedDemoArticles(tx, admin.ID, categories, tags)
		if err != nil {
			return err
		}
		if err := seedDemoComments(tx, articles, admin.ID, guest.ID); err != nil {
			return err
		}
		if err := seedDemoMessages(tx); err != nil {
			return err
		}
		return seedDemoLinks(tx)
	})
}

// 样例内容里的图片同样走本仓库的 images/ 目录
const demoImgBase = "https://raw.githubusercontent.com/szluyu99/gin-vue-blog/main/images"

var demoCovers = []string{
	demoImgBase + "/page/home.jpg",
	demoImgBase + "/page/archive.png",
	demoImgBase + "/page/category.png",
	demoImgBase + "/page/tag.png",
	demoImgBase + "/page/link.jpg",
	demoImgBase + "/page/about.jpg",
	demoImgBase + "/page/message.jpeg",
	demoImgBase + "/page/article_list.jpg",
}

func seedDemoCategories(tx *gorm.DB) ([]Category, error) {
	categories := []Category{
		{Name: "后端开发"},
		{Name: "前端开发"},
		{Name: "随笔"},
	}
	for i := range categories {
		if err := tx.Create(&categories[i]).Error; err != nil {
			return nil, err
		}
	}
	return categories, nil
}

func seedDemoTags(tx *gorm.DB) ([]Tag, error) {
	tags := []Tag{
		{Name: "Go"},
		{Name: "Gin"},
		{Name: "GORM"},
		{Name: "Vue"},
		{Name: "Vite"},
		{Name: "生活"},
	}
	for i := range tags {
		if err := tx.Create(&tags[i]).Error; err != nil {
			return nil, err
		}
	}
	return tags, nil
}

// 一篇样例文章的差异部分, 公共部分(正文骨架/封面/时间)由 seedDemoArticles 补齐
type demoArticle struct {
	title       string
	desc        string
	category    int   // seedDemoCategories 返回值的下标
	tags        []int // seedDemoTags 返回值的下标
	status      int   // 0 表示 STATUS_PUBLIC
	typ         int   // 0 表示 TYPE_ORIGINAL
	isTop       bool
	originalUrl string
}

var demoArticleSpecs = []demoArticle{
	{title: "Gin 项目的分层结构该怎么划分", desc: "handler 只做参数校验和响应封装, 业务逻辑下沉到 service, 数据访问收在 model 里。", category: 0, tags: []int{0, 1}, isTop: true},
	{title: "GORM 的 Updates 为什么会漏掉零值", desc: "结构体更新会跳过零值字段, 想把 false / 0 写进去只有 Select 或 map 两条路。", category: 0, tags: []int{0, 2}},
	{title: "用中间件统一处理 JWT 鉴权", desc: "把 token 解析、用户加载、权限判断拆开, 顺序比实现更容易出错。", category: 0, tags: []int{0, 1}},
	{title: "接口分页返回值的几种写法", desc: "总数和列表一起返回, 前端才能算出总页数, 只返回列表的接口迟早要改。", category: 0, tags: []int{0, 1}},
	{title: "多对多关联的更新陷阱", desc: "全量替换关联表意味着漏传的 id 会被删掉, 局部更新要么补齐要么单独写 SQL。", category: 0, tags: []int{0, 2}},
	{title: "slog 结构化日志入门", desc: "标准库自带的结构化日志, 不引第三方也够用, 关键是把 handler 配好。", category: 0, tags: []int{0}, typ: TYPE_TRANSLATE},
	{title: "Vue 3 里解构 pinia store 会丢响应式", desc: "store 重新赋值整个对象时, 解构出来的变量还指向旧的引用。", category: 1, tags: []int{3}},
	{title: "v-for 里的模板 ref 没有顺序保证", desc: "ref 数组的下标和渲染顺序不一定对得上, 用函数式 ref 按 key 存更稳。", category: 1, tags: []int{3}},
	{title: "Vite 构建产物的体积优化", desc: "先看清 chunk 是怎么切的, 再决定要不要手动分包。", category: 1, tags: []int{3, 4}},
	{title: "UnoCSS 默认不扫描 js 文件", desc: "写在 js 里的类名不会生成 CSS, 动态拼类名前先确认 content.pipeline 的范围。", category: 1, tags: []int{3, 4}, typ: TYPE_REPRINT, originalUrl: "https://unocss.dev/guide/extracting"},
	{title: "组件测试从哪里开始写", desc: "先覆盖加载态和空数据, 这两处的 bug 比交互逻辑更多。", category: 1, tags: []int{3}},
	{title: "写博客这件事坚持了三年", desc: "记录的价值不在阅读量, 而在于半年后自己还能看懂当时的思路。", category: 2, tags: []int{5}},
	{title: "一次线上事故的复盘", desc: "配置项的默认值太宽松, 于是没人发现它其实没生效。", category: 2, tags: []int{5}},
	{title: "还没写完的草稿", desc: "草稿状态的文章不会出现在前台, 用来验证状态过滤。", category: 2, tags: []int{5}, status: STATUS_DRAFT},
	{title: "只有自己能看的私密文章", desc: "私密文章同样不出现在前台列表和归档里。", category: 2, tags: []int{5}, status: STATUS_SECRET},
}

// 正文骨架: 带多级标题、列表、代码块和引用, 用来验证 markdown 渲染、目录和代码复制按钮
func demoContent(a demoArticle) string {
	return fmt.Sprintf(`## 背景

%s

这是一篇样例文章, 内容仅用于本地测试前台的渲染效果。

## 要点

- 结论先写在前面, 过程放在后面
- 代码片段尽量能直接跑
- 需要注意的坑单独列出来

## 代码

`+"```go"+`
func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.Run(":8765")
}
`+"```"+`

## 小结

> %s

更多内容见项目仓库。
`, a.desc, a.title)
}

func seedDemoArticles(tx *gorm.DB, userId int, categories []Category, tags []Tag) ([]Article, error) {
	if len(demoArticleSpecs) != demoArticleCount {
		return nil, fmt.Errorf("样例文章数量应为 %d, 实际 %d", demoArticleCount, len(demoArticleSpecs))
	}

	now := time.Now()
	articles := make([]Article, 0, len(demoArticleSpecs))

	// 时间从新到旧每篇差 12 天, 跨越半年多, 归档页才有多个月份分组;
	// 靠后的几篇超过 90 天, 顺带能看到"老文章"提示
	for i, spec := range demoArticleSpecs {
		createdAt := now.AddDate(0, 0, -12*i)

		status := spec.status
		if status == 0 {
			status = STATUS_PUBLIC
		}
		typ := spec.typ
		if typ == 0 {
			typ = TYPE_ORIGINAL
		}

		articleTags := make([]*Tag, 0, len(spec.tags))
		for _, idx := range spec.tags {
			articleTags = append(articleTags, &tags[idx])
		}

		article := Article{
			Model:       Model{CreatedAt: createdAt, UpdatedAt: createdAt},
			Title:       spec.title,
			Desc:        spec.desc,
			Content:     demoContent(spec),
			Img:         demoCovers[i%len(demoCovers)],
			Type:        typ,
			Status:      status,
			IsTop:       spec.isTop,
			OriginalUrl: spec.originalUrl,
			CategoryId:  categories[spec.category].ID,
			UserId:      userId,
			Tags:        articleTags,
		}
		if err := tx.Create(&article).Error; err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// 评论: 前两篇文章各挂一条评论 + 一条回复, 再留一条未审核的用来验证后台审核流程
func seedDemoComments(tx *gorm.DB, articles []Article, adminId, guestId int) error {
	if len(articles) < 2 {
		return fmt.Errorf("样例文章不足, 无法生成评论")
	}

	now := time.Now()
	type spec struct {
		userId      int
		replyUserId int
		topicId     int
		parentIdx   int // -1 表示顶级评论, 否则指向本切片中已创建的评论下标
		typ         int
		content     string
		isReview    bool
	}

	specs := []spec{
		{userId: guestId, topicId: articles[0].ID, parentIdx: -1, typ: TYPE_ARTICLE, content: "分层这块讲得很清楚, 之前我一直把 SQL 写在 handler 里。", isReview: true},
		{userId: adminId, replyUserId: guestId, parentIdx: 0, typ: TYPE_ARTICLE, content: "刚开始都这样, 拆出来之后写测试会轻松很多。", isReview: true},
		{userId: guestId, topicId: articles[1].ID, parentIdx: -1, typ: TYPE_ARTICLE, content: "零值这个坑我踩过, 后来统一改成 map 更新了。", isReview: true},
		{userId: adminId, replyUserId: guestId, parentIdx: 2, typ: TYPE_ARTICLE, content: "map 确实更直观, 缺点是字段名没有编译期检查。", isReview: true},
		{userId: guestId, topicId: articles[1].ID, parentIdx: -1, typ: TYPE_ARTICLE, content: "这条评论还没审核, 前台看不到。", isReview: false},
		{userId: guestId, parentIdx: -1, typ: TYPE_LINK, content: "申请友链: 我的博客也是用这个项目搭的。", isReview: true},
	}

	created := make([]Comment, 0, len(specs))
	for i, s := range specs {
		createdAt := now.AddDate(0, 0, i-len(specs))
		comment := Comment{
			Model:       Model{CreatedAt: createdAt, UpdatedAt: createdAt},
			UserId:      s.userId,
			ReplyUserId: s.replyUserId,
			TopicId:     s.topicId,
			Content:     s.content,
			Type:        s.typ,
			IsReview:    s.isReview,
		}
		if s.parentIdx >= 0 {
			parent := created[s.parentIdx]
			comment.ParentId = parent.ID
			comment.TopicId = parent.TopicId
			comment.Type = parent.Type
		}
		if err := tx.Create(&comment).Error; err != nil {
			return err
		}
		created = append(created, comment)
	}

	return nil
}

// 留言(弹幕): speed 影响滚动速度, 留一条未审核的
func seedDemoMessages(tx *gorm.DB) error {
	now := time.Now()
	messages := []Message{
		{Nickname: "路过的猫", Content: "博客界面挺好看的!", Speed: 8, IsReview: true},
		{Nickname: "老王", Content: "什么时候更新说说功能?", Speed: 6, IsReview: true},
		{Nickname: "匿名", Content: "第一次留言, 打个卡。", Speed: 10, IsReview: true},
		{Nickname: "小明", Content: "这条留言还没审核。", Speed: 5, IsReview: false},
	}

	for i := range messages {
		messages[i].Avatar = DefaultAvatar
		messages[i].IpAddress = "127.0.0.1"
		messages[i].IpSource = "内网IP"
		messages[i].CreatedAt = now.AddDate(0, 0, -i)
		messages[i].UpdatedAt = messages[i].CreatedAt
		if err := tx.Create(&messages[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedDemoLinks(tx *gorm.DB) error {
	links := []FriendLink{
		{Name: "Go 官方博客", Avatar: demoImgBase + "/common/header.jpeg", Address: "https://go.dev/blog/", Intro: "Go 团队的官方博客"},
		{Name: "Vue 官方文档", Avatar: demoImgBase + "/common/header.jpeg", Address: "https://cn.vuejs.org/", Intro: "Vue 3 中文文档"},
		{Name: "gin-vue-blog", Avatar: demoImgBase + "/common/header.jpeg", Address: "https://github.com/szluyu99/gin-vue-blog", Intro: "本项目的源码仓库"},
	}

	for i := range links {
		if err := tx.Create(&links[i]).Error; err != nil {
			return err
		}
	}

	return nil
}
