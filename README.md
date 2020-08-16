
PlistDumper 是一个拆图工具。游戏发布的时候通常会采用和图来提高游戏运行效率，PlistDumper可以根据和图的配置文件拆分出子图片，并且还原图片真实大小。

* 支持TexturePacker各种版本的plist文件导出
* 支持TexturePacker部分json文件导出
* 支持fnt位图字体文件导出
* `golang` 开发，跨平台，可运行在Windows、Mac、Linux。

## 安装
* 首先安装golang环境
* 执行 go get -u -v github.com/qcdong2016/PlistDumper


## 使用说明
```
$ PlistDumper [plist|json|fnt|dir]
```
* 第一个参数传文件或者目录
* 不传参数等于传当前目录

# 预览

![preview](./preview.jpg)
