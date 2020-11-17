### sonar-qyhook

[![Build Status](https://travis-ci.org/CharLemAznable/sonar-qyhook.svg?branch=master)](https://travis-ci.org/CharLemAznable/sonar-qyhook)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/CharLemAznable/sonar-qyhook)
[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.svg?v=103)](https://opensource.org/licenses/mit-license.php)
![GitHub code size](https://img.shields.io/github/languages/code-size/CharLemAznable/sonar-qyhook)

SonarQube项目分析后回调Webhook服务, 将分析结果转发为企业微信应用推送消息.

#### 配置文件

```config.toml``` [示例](https://github.com/CharLemAznable/sonar-qyhook/blob/master/config.toml)

```toml
Port = 17258
ContextPath = ""
LogLevel = "info"

VarysBaseUrl = ""       # varys服务地址
QyWxAgentId = ""        # 企业微信应用ID, 即varys配置的企业应用codeName

ProjectKeyPattern = ""  # 按正则匹配需要发送消息的Sonar项目名称, 默认为: "^.*$"

ShieldsBadgeUrl = ""    # 徽章反向代理地址, 默认为: "https://img.shields.io/static/v1"
```

#### 部署执行

1. 下载最新的可执行文件压缩包并解压

    下载地址: [sonar-qyhook release](https://github.com/CharLemAznable/sonar-qyhook/releases)

```bash
$ tar -xvJf sonar-qyhook-[version].[arch].[os].tar.xz
```

2. 新建/编辑配置文件, 启动运行

```bash
$ nohup ./sonar-qyhook-[version].[arch].[os].bin &
```
