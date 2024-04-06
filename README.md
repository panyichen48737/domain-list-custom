# 文件说明
## 1. 文件类型
本项目生成的文件为 [v2fly/domain-list-community/data](https://github.com/v2fly/domain-list-community) 项目的 .txt 格式，包含 `domain`、`full` 和 `regexp` 类型
## 2. 数据源
① 每天凌晨 2 点（北京时间）自动构建  
② private.txt 源采用 [v2fly/domain-list-community/private](https://github.com/v2fly/domain-list-community/blob/master/data/private)  
③ microsoft.txt 源采用 [v2fly/domain-list-community/microsoft](https://github.com/v2fly/domain-list-community/blob/master/data/microsoft)  
④ category-games.txt 源采用 [v2fly/domain-list-community/category-games](https://github.com/v2fly/domain-list-community/blob/master/data/category-games)  
⑤ speedtest.txt 源采用 [v2fly/domain-list-community/speedtest](https://github.com/v2fly/domain-list-community/blob/master/data/speedtest)  
⑥ geolocation-!cn.txt 源采用 [v2fly/domain-list-community/geolocation-!cn](https://github.com/v2fly/domain-list-community/blob/master/data/geolocation-!cn)（删除了带有 `@cn` 和 `@ads` 的域名，并新增了 v2fly/domain-list-community/cn 中带有 `@!cn` 的域名）  
⑦ cn.txt 源采用 [v2fly/domain-list-community/cn](https://github.com/v2fly/domain-list-community/blob/master/data/cn)（删除了带有 `@!cn` 和 `@ads` 的域名，并新增了 v2fly/domain-list-community/geolocation-!cn 中带有 `@cn` 的域名）
