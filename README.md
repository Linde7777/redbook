# redbook
模仿小红书

# 目录结构
借鉴DDD(Domain-Driven Design)思想
## internal
domain存放的是业务对象；repository负责操作业务对象，对外屏蔽实现细节（数据库、缓存等）；
service是业务逻辑，调用repository实现功能，可以组合其他的业务逻辑；web是对外提供的接口，调用service实现功能。

## ioc
ioc是inversion of control的缩写，即控制反转。用于Wire实现自动依赖注入
