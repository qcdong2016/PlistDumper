
PlistDumper
--------------
PlistDumper 是一个拆图工具。读取plist或者json文件，导出图片。

# 中文 (Chinese)

从 `Zwoptex` 或者 `TexturePacker` 的 `plist|json` 文件中导出图片，还原图片真实大小。支持各种版本的plist文件，和部分json文件。类似于 `TextureUnpacker`。

采用 `golang` 开发，支持各种操作系统。

## 使用说明

```
$ PlistDumper [plistfile|jsonfile|targetdir]
```
* 第一个参数传plist文件/json文件路径，或者目录
* 不传参数等于传当前目录

# English

Export the image from the `TexturePacker` or `Zwoptex`'s  `plist` file, support all version of plist, similar to `TextureUnpacker`.

Developed with `golang`, and support most popular operating systems.

## Usage

```
    PlistDumper [plistfile|jsonfile|targetdir]
```

# preview / 预览

![preview](./preview.jpg)
