# cfg 配置中心
- 读取配置顺序 内存缓存 -> Redis -> Mysql -> 本地Yaml文件 -> 环境变量
- 以Mysql数据为基准数据, 每次设置会更新Mysql , 同时清空缓存
- 提供grpc 接口 ,提供grpc-gateway 形式作为http接口,进行数据读取

### todo : 
1. etcd 作为存储介质,然后增加同步内存缓存策略
2. websocket接口




