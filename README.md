# 文件说明
## 1. 文件类型
本项目可生成 .list 文件，适用于 `behavior: classical` 且 `format: text` 的 [Clash](https://github.com/Dreamacro/clash) 使用场景，其中包含 `DOMAIN`、`DOMAIN-SUFFIX`、`DOMAIN-KEYWORD` 和 `PROCESS-NAME` 规则类型
## 2. 数据源
① 每天凌晨 2 点（北京时间）自动构建  
② **`ads.list`** 源采用 [privacy-protection-tools/anti-AD](https://github.com/privacy-protection-tools/anti-AD)  
③ **`private.txt`** 源采用 [v2fly/domain-list-community/private](https://github.com/v2fly/domain-list-community/blob/master/data/private)、[blackmatrix7/ios_rule_script/Lan](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Lan)（仅域名）和 [TrackersList](https://github.com/XIU2/TrackersListCollection/blob/master/all.txt)（仅域名）组合，并添加[常用 fake-ip 地址过滤列表](https://github.com/juewuy/ShellCrash/blob/dev/public/fake_ip_filter.list)里的所有 [AdGuardHome](https://github.com/AdguardTeam/AdGuardHome) 相关域名和主流 [Clash dashboard 在线面板](https://github.com/DustinWin/clash_singbox-tools/tree/main/Clash-dashboard)域名（`clash.razord.top`、`clash.metacubex.one`、`yacd.haishan.me`、`yacd.metacubex.one` 和 `d.metacubex.one`）  
④ **`microsoft-cn.list`** 源采用 [v2fly/domain-list-community/microsoft@cn](https://github.com/v2fly/domain-list-community/blob/master/data/microsoft)  
⑤ **`apple-cn.list`** 源采用 [v2fly/domain-list-community/apple@cn](https://github.com/v2fly/domain-list-community/blob/master/data/apple)  
⑥ **`google-cn.list`** 源采用 [v2fly/domain-list-community/google@cn](https://github.com/v2fly/domain-list-community/blob/master/data/google)  
⑦ **`games-cn.list`** 源采用 [v2fly/domain-list-community/category-games@cn](https://github.com/v2fly/domain-list-community/blob/master/data/category-games)、[blackmatrix7/ios_rule_script/SteamCN](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/SteamCN) 和 [blackmatrix7/ios_rule_script/GameDownloadCN](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Game/GameDownloadCN) 组合  
⑧ **`netflix.list`** 源采用 [blackmatrix7/ios_rule_script/Netflix](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Netflix)（仅域名）  
⑨ **`disney.list`** 源采用 [blackmatrix7/ios_rule_script/Disney](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Disney)  
⑩ **`max.list`** 源采用 [blackmatrix7/ios_rule_script/HBO](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/HBO)  
⑪ **`primevideo.list`** 源采用 [blackmatrix7/ios_rule_script/PrimeVideo](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/PrimeVideo)  
⑫ **`appletv.list`** 源采用 [blackmatrix7/ios_rule_script/AppleTV](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/AppleTV)  
⑬ **`youtube.list`** 源采用 [blackmatrix7/ios_rule_script/YouTube](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/YouTube)  
⑭ **`tiktok.list`** 源采用 [blackmatrix7/ios_rule_script/TikTok](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/TikTok)  
⑮ **`bilibili.list`** 源采用 [blackmatrix7/ios_rule_script/BiliBili](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/BiliBili)  
⑯ **`ai.list`** 源采用 [blackmatrix7/ios_rule_script/OpenAI](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/OpenAI)、[blackmatrix7/ios_rule_script/Copilot](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Copilot)、[blackmatrix7/ios_rule_script/Gemini](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Gemini) 和 [blackmatrix7/ios_rule_script/Claude](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Claude) 组合  
⑰ **`networktest.list`** 源采用 [v2fly/domain-list-community/speedtest](https://github.com/v2fly/domain-list-community/blob/master/data/speedtest)、[blackmatrix7/ios_rule_script/Speedtest](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Speedtest) 和 IPv6 测试域名关键字（`keyword`，包括：`ipv6-test`、`test-ipv6`、`ipv6test` 和 `testipv6`）组合  
⑱ **`proxy.list`** 源采用 [cokebar/gfwlist2dnsmasq](https://github.com/cokebar/gfwlist2dnsmasq) 生成的 [gfwlist](https://github.com/gfwlist/gfwlist)、[v2fly/domain-list-community/geolocation-!cn](https://github.com/v2fly/domain-list-community/blob/master/data/geolocation-!cn)（删除了带有 `@cn` 和 `@ads` 的域名，并新增了 v2fly/domain-list-community/cn 中带有 `@!cn` 的域名）和 [blackmatrix7/ios_rule_script/Global](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/Global) 组合  
⑲ **`cn.list`** 源采用 [v2fly/domain-list-community/cn](https://github.com/v2fly/domain-list-community/blob/master/data/cn)（删除了带有 `@!cn` 和 `@ads` 的域名，并新增了 v2fly/domain-list-community/geolocation-!cn 中带有 `@cn` 的域名）和 [blackmatrix7/ios_rule_script/ChinaMax](https://github.com/blackmatrix7/ios_rule_script/tree/master/rule/Clash/ChinaMax)（仅域名）组合  
