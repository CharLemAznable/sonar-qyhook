### sonar-qyhook

[![Build Status](https://travis-ci.org/CharLemAznable/sonar-qyhook.svg?branch=master)](https://travis-ci.org/CharLemAznable/sonar-qyhook)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/CharLemAznable/sonar-qyhook)
[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.svg?v=103)](https://opensource.org/licenses/mit-license.php)
![GitHub code size](https://img.shields.io/github/languages/code-size/CharLemAznable/sonar-qyhook)

SonarQube项目分析后回调Webhook服务, 将分析结果转发为企业微信应用推送消息.

#### 配置文件

1. ```appConfig.toml```

```toml
Port = 17258
ContextPath = ""
VarysBaseUrl = ""       # varys服务地址
QyWxAgentId = ""        # 企业微信应用ID, 即varys配置的企业应用codeName
ProjectKeyPattern = ""  # 按正则匹配需要发送消息的Sonar项目名称
```

2. ```logback.xml```

```xml
<logging>
    <filter enabled="true">
        <tag>file</tag>
        <type>file</type>
        <level>INFO</level>
        <property name="filename">sonar-qyhook.log</property>
        <property name="format">[%D %T] [%L] (%S) %M</property>
        <property name="rotate">false</property>
        <property name="maxsize">0M</property>
        <property name="maxlines">0K</property>
        <property name="daily">false</property>
    </filter>
</logging>
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
