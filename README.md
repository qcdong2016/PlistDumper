
PlistDumper
--------------

Original author: [qcdong2016](https://github.com/qcdong2016/PlistDumper) | Modified by: [shines77](https://github.com/shines77/PlistDumper)

Special thanks to: https://github.com/DHowett/go-plist

# 中文 (Chinese)

从 `TexturePacker` 的 `plist` 文件中导出 `sprite` 图片，类似于 `TextureUnpacker`。

采用 `golang` 开发，支持各种操作系统。

## 使用说明

```
$ PlistDumper [format] [plistfile]
```

* `format`: `plist` 文件的格式，可选格式有：`cocos2dx`, `std`。
* `plistfile`：`plist` 文件名，例如：`abc.plist`。

### 范例

```
$ PlistDumper abc.plist
$ PlistDumper cocos2dx abc.plist
$ PlistDumper std abc.plist
```

### 依赖库

安装 `go-plist` 库：

```
$ go get github.com/DHowett/go-plist
```
或者
```
$ go get howett.net/plist
```

# English

Export the `sprite` image from the `TexturePacker`'s  `plist` file, similar to `TextureUnpacker`.

Developed with `golang`, and support most popular operating systems.

## Usage

```
    PlistDumper [format] [plistfile]
```

* `format`: `plist` file format, options include：`cocos2dx`, `std`.
* `plistfile`：`plist` file name, example: `abc.plist`.

### Examples

```
    PlistDumper abc.plist
    PlistDumper cocos2dx abc.plist
    PlistDumper std abc.plist
```

### Dependent library

Install `go-plist` library：

```
$ go get github.com/DHowett/go-plist
```
Or
```
$ go get howett.net/plist
```

# preview / 预览

![preview](./preview.jpg)
