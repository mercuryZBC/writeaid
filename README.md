# WriteAid-知识记录与分享平台
## 项目简介
前后端分离的个人知识库平台，支持用户创建和发布文档、搜索功能、多级评论、点赞功能。支持用户查看最近浏览、最近编辑、最近评论的文档。
## 使用技术
前端：使用React 框架与 Ant Design UI组件库进行开发
后端：
开发语言：Go

web开发框架：Gin

orm框架：Gorm

数据库：Mysql

缓存：Redis

搜索引擎：ElasticSearch
## 项目部署
本项目支持在Linux平台或MacOS平台部署，在部署前用户需要安装docker，docker-compose，nodejs，npm
```bash
#执行一键部署脚本
./deploy.sh
```
大家可以随时提issue，共同开发和完善该项目。有部署上的问题也可以联系我，本人QQ：2950206203
## 项目架构
本项目分为接入层，用户鉴权层，微服务层

接入层主要实现支持api接入，网关服务，以及负载均衡。

用户鉴权层采用jwt实现权限校验。

微服务层将用户服务、搜索服务、知识库服务、文档服务、消息推送服务、评论服务进行划分。

![image](https://github.com/user-attachments/assets/2b063fb0-f30c-4a1e-a560-0922b194fcd4)
## 项目图片
Dashboard页面
![image]("https://github.com/user-attachments/assets/0f10be39-4a76-408a-8f35-6ffd7f43789e")
文档预览页面
![image]("https://github.com/user-attachments/assets/2d56657d-00bd-4f78-a5fc-bfb09715c887")
文档编辑页面
![image]("https://github.com/user-attachments/assets/c4c412a3-7c08-4921-aff6-23894f1ff080")
通知窗口
![image](https://github.com/user-attachments/assets/918f4d00-bc05-46cb-bb4f-75a811e3e80f)
知识库和文档搜索
![image](https://github.com/user-attachments/assets/a758a695-6a24-4925-9c3c-2a63fcfd4c15)




