# 文件说明
## 1. 文件类型
① 重构上游项目 [v2fly/domain-list-community](https://github.com/v2fly/domain-list-community)，生成供下游项目 [panyichen48737/ruleset_geodata](https://github.com/panyichen48737/ruleset_geodata) 使用的 DOMAIN 数据源文件  
② 数据源文件为 [mihomo 内核](https://github.com/MetaCubeX/mihomo) rule-set 规则集文件（.list 格式），包含 `DOMAIN`、`DOMAIN-SUFFIX`、`DOMAIN-KEYWORD`、`DOMAIN-REGEX` 和 `PROCESS-NAME` 规则类型，配置 `behavior: classical` 和 `format: text` 后可直接使用
## 2. 数据去重
① 去除 `DOMAIN-SUFFIX` 类型中冗余的子域名，如：
```
DOMAIN-SUFFIX,test.cn
DOMAIN-SUFFIX,cn
DOMAIN-SUFFIX,123.example.com
DOMAIN-SUFFIX,abc.example.com
DOMAIN-SUFFIX,example.com
```
去除冗余子域名后，仅保留：
```
DOMAIN-SUFFIX,cn
DOMAIN-SUFFIX,example.com
```
② 去除 `DOMAIN-SUFFIX` 类型中已存在，`DOMAIN` 类型中冗余的域名，如：
```
DOMAIN,test.cn
DOMAIN-SUFFIX,cn
DOMAIN,123.example.com
DOMAIN,test.abc.example.com
DOMAIN-SUFFIX,example.com
```
去除 `DOMAIN` 类型中冗余的域名后，仅保留：
```
DOMAIN-SUFFIX,cn
DOMAIN-SUFFIX,example.com
```
## 3. 上游 `include` 语法
① `include:文件名`：纳入该文件的全部域名  
② `include:文件名 @属性`：仅纳入带有该属性的域名，如 `include:google @cn`  
③ `include:文件名 @-属性`：纳入**不**带有该属性的域名（未标属性的域名也计入），即上游 v2fly 的 BanAttrs 语义，如 `include:tencent @-!cn`  
④ 同一行可写多个属性，条件同时生效：`@属性` 是「必须全部命中」（与上游 v2fly 的 `MustAttrs` 一致），`@-属性` 是「一个都不能命中」。如 `include:child @cn @ads @-!cn` 表示「带 `@cn` 和 `@ads`，且不带 `@!cn`」
## 4. 数据源
① 每天凌晨 2 点（北京时间 UTC+8）自动构建  
② **`fakeip-filter.list`** 源采用 [DustinWin/fake_ip_filter.list](https://github.com/DustinWin/ShellCrash/blob/dev/public/fake_ip_filter.list)  
③ **`fakeip-filter-lite.list`** 源采用 [DustinWin/fake_ip_filter.list](https://github.com/DustinWin/ShellCrash/blob/dev/public/fake_ip_filter.list)，仅保留主要域名（推荐搭配 [AdGuard Home](https://github.com/AdguardTeam/AdGuardHome) 且 DNS 配置 mix 混合模式时使用）  
④ **`private.list`** 源采用 [v2fly/domain-list-community/private](https://github.com/v2fly/domain-list-community/blob/master/data/private) 和 [blackmatrix7/ios_rule_script/Lan](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Lan)（仅域名）组合，并添加 `*` 和主流 [Dashboard 在线面板](https://github.com/DustinWin/proxy-tools/releases/tag/Dashboard)域名（`yacd.haishan.me`、`yacd.metacubex.one`、`metacubex.github.io`、`metacubexd.pages.dev` 和 `board.zash.run.place`）  
⑤ **`trackerslist.list`** 源采用 [XIU2/TrackersListCollection](https://github.com/XIU2/TrackersListCollection/blob/master/all.txt)（仅域名）和 [ngosang/trackerslist](https://github.com/ngosang/trackerslist/blob/master/trackers_all.txt) 组合  
⑥ **`microsoft-cn.list`** 源采用 [v2fly/domain-list-community/microsoft@cn](https://github.com/v2fly/domain-list-community/blob/master/data/microsoft)，并删除 `onedrive.list`、`bing.list` 和 `ai.list` 中已存在的域名  
⑦ **`onedrive.list`** 源采用 [blackmatrix7/ios_rule_script/OneDrive](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/OneDrive)  
⑧ **`bing.list`** 源采用 [blackmatrix7/ios_rule_script/Bing](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Bing)  
⑨ **`apple-proxy.list`** 源采用 [blackmatrix7/ios_rule_script/AppleProxy](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/AppleProxy)  
⑩ **`apple-cn.list`** 源采用 [v2fly/domain-list-community/apple@cn](https://github.com/v2fly/domain-list-community/blob/master/data/apple) 和 [felixonmars/dnsmasq-china-list/apple.china.conf](https://github.com/felixonmars/dnsmasq-china-list/blob/master/apple.china.conf) 组合，并补充 [blackmatrix7/ios_rule_script/Apple](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Apple) 的进程、关键字和 `DOMAIN` 规则（不含其 `DOMAIN-SUFFIX` 规则），最后删除 `apple-proxy.list` 和 `appletv.list` 中已存在的域名（含被其 `DOMAIN-SUFFIX` 父域罩住的域名）  
⑪ **`google-cn.list`** 源采用 [v2fly/domain-list-community/google@cn](https://github.com/v2fly/domain-list-community/blob/master/data/google) 和 [felixonmars/dnsmasq-china-list/google.china.conf](https://github.com/felixonmars/dnsmasq-china-list/blob/master/google.china.conf) 组合  
⑫ **`games-cn.list`** 源采用 [v2fly/domain-list-community/category-game-accelerator-cn](https://github.com/v2fly/domain-list-community/blob/master/data/category-game-accelerator-cn)、[v2fly/domain-list-community/category-game-platforms-download@cn](https://github.com/v2fly/domain-list-community/blob/master/data/category-game-platforms-download)、[v2fly/domain-list-community/category-games-cn](https://github.com/v2fly/domain-list-community/blob/master/data/category-games-cn)、[blackmatrix7/ios_rule_script/SteamCN](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/SteamCN)、[blackmatrix7/ios_rule_script/GameDownloadCN](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Game/GameDownloadCN) 和 [blackmatrix7/ios_rule_script/Game/GameDownload](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Game/GameDownload) 组合  
⑬ **`netflix.list`** 源采用 [v2fly/domain-list-community/netflix](https://github.com/v2fly/domain-list-community/blob/master/data/netflix) 和 [blackmatrix7/ios_rule_script/Netflix](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Netflix)（仅域名）组合  
⑭ **`disney.list`** 源采用 [v2fly/domain-list-community/disney](https://github.com/v2fly/domain-list-community/blob/master/data/disney) 和 [blackmatrix7/ios_rule_script/Disney](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Disney) 组合  
⑮ **`hbo.list`** 源采用 [v2fly/domain-list-community/hbo](https://github.com/v2fly/domain-list-community/blob/master/data/hbo) 和 [blackmatrix7/ios_rule_script/HBO](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/HBO) 组合  
⑯ **`primevideo.list`** 源采用 [v2fly/domain-list-community/primevideo](https://github.com/v2fly/domain-list-community/blob/master/data/primevideo) 和 [blackmatrix7/ios_rule_script/PrimeVideo](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/PrimeVideo) 组合  
⑰ **`appletv.list`** 源采用 [v2fly/domain-list-community/apple-tvplus](https://github.com/v2fly/domain-list-community/blob/master/data/apple-tvplus) 和 [blackmatrix7/ios_rule_script/AppleTV](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/AppleTV) 组合  
⑱ **`youtube.list`** 源采用 [v2fly/domain-list-community/youtube](https://github.com/v2fly/domain-list-community/blob/master/data/youtube) 和 [blackmatrix7/ios_rule_script/YouTube](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/YouTube) 组合  
⑲ **`tiktok.list`** 源采用 [v2fly/domain-list-community/tiktok](https://github.com/v2fly/domain-list-community/blob/master/data/tiktok) 和 [blackmatrix7/ios_rule_script/TikTok](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/TikTok) 组合  
⑳ **`bilibili.list`** 源采用 [v2fly/domain-list-community/bilibili](https://github.com/v2fly/domain-list-community/blob/master/data/bilibili) 和 [blackmatrix7/ios_rule_script/BiliBili](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/BiliBili) 组合  
㉑ **`spotify.list`** 源采用 [v2fly/domain-list-community/spotify](https://github.com/v2fly/domain-list-community/blob/master/data/spotify) 和 [blackmatrix7/ios_rule_script/Spotify](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Spotify) 组合  
㉒ **`games.list`** 源采用 [v2fly/domain-list-community/category-game-platforms-download](https://github.com/v2fly/domain-list-community/blob/master/data/category-game-platforms-download)、[v2fly/domain-list-community/category-games-!cn](https://github.com/v2fly/domain-list-community/blob/master/data/category-games-!cn)、[blackmatrix7/ios_rule_script/Game/Nikke](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Game/Nikke)、[blackmatrix7/ios_rule_script/Game/Roblox](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Game/Roblox)、[blackmatrix7/ios_rule_script/Game/TencentLoLMobile](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Game/TencentLoLMobile)、[blackmatrix7/ios_rule_script/Game](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Game) 和 [blackmatrix7/ios_rule_script/Game/GameDownload](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Game/GameDownload) 组合，最后删除 `games-cn.list` 中已存在的域名（含被其 `DOMAIN-SUFFIX` 父域罩住的域名）  
㉓ **`ai.list`** 源采用 [v2fly/domain-list-community/category-ai-!cn](https://github.com/v2fly/domain-list-community/blob/master/data/category-ai-!cn)、[ACL4SSR/ACL4SSR/AI.list](https://github.com/ACL4SSR/ACL4SSR/blob/master/Clash/Ruleset/AI.list)、[blackmatrix7/ios_rule_script/OpenAI](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/OpenAI)、[blackmatrix7/ios_rule_script/Copilot](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Copilot)、[blackmatrix7/ios_rule_script/Gemini](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Gemini) 和 [blackmatrix7/ios_rule_script/Claude](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Claude) 组合  
㉔ **`networktest.list`** 源采用 [v2fly/domain-list-community/category-speedtest](https://github.com/v2fly/domain-list-community/blob/master/data/category-speedtest)、[v2fly/domain-list-community/category-ip-geo-detect](https://github.com/v2fly/domain-list-community/blob/master/data/category-ip-geo-detect)、[blackmatrix7/ios_rule_script/Speedtest](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Speedtest) 和 IPv6 测试域名关键字（`keyword`，包括：`ipv6-test`、`test-ipv6`、`ipv6test` 和 `testipv6`）组合  
㉕ **`tld-proxy.list`** 源采用 [v2fly/domain-list-community/tld-!cn](https://github.com/v2fly/domain-list-community/blob/master/data/tld-!cn)  
㉖ **`tld-cn.list`** 源采用 [v2fly/domain-list-community/tld-cn](https://github.com/v2fly/domain-list-community/blob/master/data/tld-cn)  
㉗ **`proxy.list`** 源采用 [v2fly/domain-list-community/geolocation-!cn](https://github.com/v2fly/domain-list-community/blob/master/data/geolocation-!cn)（删除了带有 `@cn` 和 `@ads` 的域名，并新增了 [v2fly/domain-list-community/cn](https://github.com/v2fly/domain-list-community/blob/master/data/cn) 中带有 `@!cn` 的域名）、[gfwlist](https://github.com/gfwlist/gfwlist) 和 [blackmatrix7/ios_rule_script/Global](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Global) 组合  
㉘ **`cn.list`** 源采用 [v2fly/domain-list-community/cn](https://github.com/v2fly/domain-list-community/blob/master/data/cn)（删除了带有 `@!cn` 和 `@ads` 的域名，并新增了 [v2fly/domain-list-community/geolocation-!cn](https://github.com/v2fly/domain-list-community/blob/master/data/geolocation-!cn) 中带有 `@cn` 的域名）、[blackmatrix7/ios_rule_script/China](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/China)、[blackmatrix7/ios_rule_script/ChinaMax](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/ChinaMax) 和 [felixonmars/dnsmasq-china-list/accelerated-domains.china.conf](https://github.com/felixonmars/dnsmasq-china-list/blob/master/accelerated-domains.china.conf) 组合，最后删除 `proxy.list` 中已存在的域名  
㉙ **`cn-lite.list`** 源采用 [v2fly/domain-list-community/cn](https://github.com/v2fly/domain-list-community/blob/master/data/cn)（删除了带有 `@!cn` 和 `@ads` 的域名，并新增了 [v2fly/domain-list-community/geolocation-!cn](https://github.com/v2fly/domain-list-community/blob/master/data/geolocation-!cn) 中带有 `@cn` 的域名）和 [blackmatrix7/ios_rule_script/China](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/China) 组合，最后删除 `proxy.list` 中已存在的域名