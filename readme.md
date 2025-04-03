## 介绍
一个简单的[MaaAssistantArknights](https://github.com/MaaAssistantArknights/MaaAssistantArknights)的docker版本，实现任务调度和持久化，使用[maa-cli](https://github.com/MaaAssistantArknights/maa-cli)对MAA调度。本项目仅实现docker化和前端。

## 命令
`docker run -p 8000:8080 --rm -itd -e CLIENT_TYPE=Offical -v /etc/localtime:/etc/localtime:ro  -v ./config:/app/config crestfallmax/maadocker:latest `
docker run环境变量
| Key           | 默认值         | 可选值                     | 描述                                   |
|---------------|----------------|----------------------------|----------------------------------------|
| `PROXY`     |  空   | `任意`  | 代理，git代理和api代理，对于一些连不上github的网络环境这个几乎是必填，但是也可以手动去cli.toml里修改url                   |
| `CLIENT_TYPE`     | `Bilibili`    | `Bilibili`、`Offical`         | 客户端类型，目前只支持b服和官服                      |

## 使用
1. 打开安卓模拟器或者安卓容器，下载安装游戏。
2. 访问容器端口，右上角设置里修改`device`为目标ip端口，如`device = "192.168.1.2:5555"`
3. 配置任务簇，已经有两个簇可供使用和参考
4. 修改任务簇状态为启用

## 注意
1. 对于当前正在执行的簇，所有修改，包括删除，都只会在下一次运行时生效，选择右上角`强行结束任务簇`使修改快速生效，但是会强行结束当前任务簇。
2. 因设计原因，infrast文件是全局生效的，所以只需要上传一次就行。[排班表生成器](https://ark.yituliu.cn/tools/schedule)生成infrast文件。

## 预览
![one](https://raw.githubusercontent.com/CoronaAustralis/maa_docker/master/doc/assets/one.png)
![two](https://raw.githubusercontent.com/CoronaAustralis/maa_docker/master/doc/assets/two.png)