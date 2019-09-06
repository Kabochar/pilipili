# 霹雳霹雳 pilipili

一个简单纯粹的在线点播案例



## 项目相关(简单说明)

-   main 函数执行，首先读取相关配置文件，初始化项目，包括，读取项目运行变量，设置日志级别，连接数据库
-   其次，装载路由，设置 middleware 中间件，设置和获取 session ，cores，currentUser 相关内容
-   相关路由内容，设置版本类型，链接相关功能接口，分发到 api  文件夹中，实现服务分发，service 文件夹中代码即是功能相关处理逻辑。

#### 文件夹存在意义

-   serializar 序列化操作，统一化实现，定义一个通用规则，避免重复造轮子。
-   model ，具体模型的对应内容，对应字段，相关函数，
-   util 工具，辅助操作相关内容。
-   conf 配置文件内容，带国际化操作。
-   cache 内容，Redis 操作，辅助数据的读取和存储。



## 历史版本操作

### V0.1 基础模型搭建

基础内容创建

#### 模块划分层面

-   模型 model
-   控制器 server
-   服务 service
-   视图
    -   序列化器 serializer



### V0.2 视频系列模块

#### 相关接口

-   投稿视频接口：POST api/v1/videos
    -   测试：JSON 带上 title 和 info
-   视频详情接口：GET api/v1/video/:id
-   视频列表接口：GET api/v1/videos
-   更新视频接口：PUT api/v1/video/:id
    -   测试：JSON 带上 title 和 info
-   删除视频接口：DELETE api/v1/video/:id
    -   注意：删除记录，只是软删除，记录删除而已。当然是对外不可见

#### 相关知识

后端中的 MVC，view 现在是直接操作 JSON 数据了。

关于状态码的返回，一般只操作 200，无论操作是否成功，返回 HTTP 200，JSON 记录并返回相关错误信息。项目中统一化了错误信息的返回，更加单纯唯一。

#### 改动文件

modified:   README.md

new file:   api/video.go

modified:   conf/i18n.go

modified:   model/migration.go

new file:   model/video.go

new file:   serializer/video.go

modified:   server/router.go

new file:   service/create_video_service.go

new file:   service/delete_video_service.go

new file:   service/list_video_service.go

new file:   service/show_video_service.go

new file:   service/update_video_service.go



### V0.3 热度排行榜

#### 相关接口

热度查询接口：GET /api/v1/rank/daily

#### 相关知识

##### Redis

获取某一视频打开量：GET "view:video:[ id ]"

播放量+1：INCR "view:video:[ id ]"，给 view:video:[ id ] 给 这个内容+1

删除某个键，删除数据：DEL "rank:daily"，删除 rank:daily 数据

清空 redis 所有数据：FLUSHALL，清空 redis 数据

排行榜加分：ZINCRBY “rank:daily” 1 “10” ，给 rank:daily 10 号视频 +1

获取排行：ZREVRANGE “rank:daily” 0 9 ，获取排行前 10 的内容

#### 改动文件

new file:   api/rank.go

new file:   cache/keys.go

new file:   cache/main.go

new file:   service/daily_rank_service.go

new file:   tasks/cron.go

new file:   tasks/rank.go

modified:   conf/conf.go

modified:   model/video.go

modified:   serializer/video.go

modified:   server/router.go

modified:   service/daily_rank_service.go

modified:   service/show_video_service.go

modified:   tasks/cron.go

modified:   tasks/rank.go





