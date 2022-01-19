# zdeploy
![](https://img.shields.io/github/license/cedar12/zdeploy)
![](https://img.shields.io/github/stars/cedar12/zdeploy)
![](https://img.shields.io/github/forks/cedar12/zdeploy)
![](https://img.shields.io/github/issues/cedar12/zdeploy)
![](https://img.shields.io/github/v/release/cedar12/zdeploy.svg)
> 部署文件工具
* 传输部署文件
* 提供shell/bat执行功能

## 安全
* 使用md5加密传输密码
* 服务端可配置多个ip白名单，支持通配符*
* 客户端只能执行服务端配置的shell/bat

## 配置
使用ini格式作为配置文件
- 服务端
```shell
# server.ini
[server]
# 监听主机
host=0.0.0.0
# 监听端口
port=39093
# 密码
pass=Zdeploy1!2@3#

# 文件传输字节数，网络差推荐128
buf=2048

# ip白名单
white=192.168.1.*,192.168.*.10

[cmd]
# shell/bat 内置解压传输的文件指令unzip
a=echo 123
b=echo 321
```
- 客户端
```shell
# client.ini
[server]
# 服务器主机ip
host=192.168.1.10
# 服务端端口
port=39093
# 密码
pass=Zdeploy1!2@3#

[deploy]
# 要传输的文件源路径
src=/code/dist.zip
# 传输到服务端的文件名
dist=dist.zip

[cmd]
# 运行服务端配置指令
path=unzip,a,b
```

## 使用
```shell
# 服务端
zdeploy -s server.ini
# 客户端
zdeploy -c client.ini
```

## 计划
- [ ] 传输文件做成服务端内置指令
- [ ] 文件格式限制
- [ ] 文件大小限制
- [ ] 延时执行指令
