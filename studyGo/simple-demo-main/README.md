# simple-demo

## 抖音项目服务端简单示例

具体功能内容参考飞书说明文档

工程无其他依赖，直接编译运行即可

```shell
go build 
./simple-demo
```

### 功能说明

接口功能不完善，仅作为示例

* 用户登录数据保存在内存中，单次运行过程中有效
* 视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可
* http://127.0.0.1:8080/static/bear.mp4
* http://127.0.0.1:8080/douyin/comment/list/
* http://127.0.0.1:8080/douyin/user/?token=zhangleidouyin
* 用户注册的数据只是存在map[string]User这个图里
* 找半天数据定义在哪，一度怀疑不需要定义。原来是也是在controller文件夹里，common.go。

### 测试数据

测试数据写在 demo_data.go 中，用于列表接口的 mock 测试