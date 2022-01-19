# zdeploy
[中文](https://github.com/cedar12/zdeploy/readme_zh.md)
> Deployment file tool
* Transfer deployment files
* Provide shell/bat execution function

## Secure
* Use md5 encryption to transmit password
* The server can be configured with multiple ip whitelists, supporting wildcards *
* The client can only execute the shell/bat configured by the server

## Config
Use ini format as configuration file
- Server
```shell
# server.ini
[server]
# listening host
host=0.0.0.0
# listening port
port=39093
# password
pass=Zdeploy1!2@3#

# File transfer bytes, the recommended network difference is 128
buf=2048

# ip whitelist
white=192.168.1.*,192.168.*.10

[cmd]
# Shell/bat has built-in unzip command to unzip the transferred file
a=echo 123
b=echo 321
```
- Client
```shell
# client.ini
[server]
# server host ip
host=192.168.1.10
# server port
port=39093
# password
pass=Zdeploy1!2@3#

[deploy]
# The source path of the file to transfer
src=/code/dist.zip
# Filename to transfer to the server
dist=dist.zip

[cmd]
# Run server configuration commands
path=unzip,a,b
```

## Usage
```shell
# Server
zdeploy -s server.ini
# Client
zdeploy -c client.ini
```
