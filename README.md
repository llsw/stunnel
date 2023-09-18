## stunnel
go 语言实现的ssh隧道转发，非国产自研。抄袭自[go-ssh-tunnel](https://github.com/dtapps/go-ssh-tunnel)，改成了秘钥登录  
使用场景：假设有数据库C, 只允许服务器A通过私有的vpc网络访问，不对外公开，需要ssh登录A，然后再访问C上的数据库。 
大概是这样  
```bash
ssh -> A:22 -> mysql -> C:3306
```
现在可以通过本程序建一个A服务器的ssh隧道，监听在本地127.0.0.1:4306，
```bash
访问127.0.0.1:4306 -> 等价于访问C:3306
```
当然一般的数据库管理工具都自带ssh隧道功能，或者你可以使用ssh本身自带的隧道命令。此工具适用于不方便执行ssh命令的地方。
## 使用
1. 把config-template.yaml复制一份，改名变成config.yaml，  
修改config.yaml中的相应的服务器配置
2. 直接在release下载二进制程序运行
```bash
# config。yaml需要和二进制程序放在同一个文件夹
# windows
stunel_x86_64.exe
# linux
stunel_linux_amd64
# macosx 暂不支持arm架构
stunel_darwin_amd64
```
3. 编译运行(需要安装make和go)
```bash
# config。yaml需要和Makefile文件放在同一个文件夹
# 直接go运行
make run
# 或者编译后运行
# linux平台
make linux
# macosx平台
make macosx
```