## DavGo

`DavGo` 是一个用 Go 语言实现的轻量级 WebDAV 服务器，支持动态配置多个 WebDAV 服务实例，每个实例可以独立设置根目录、认证信息和读写模式。

### 功能特性
- **动态配置**：通过 YAML 文件配置任意数量的 WebDAV 服务。
- **独立实例**：每个服务有独立的根目录、用户名和密码。
- **读写控制**：支持 `readonly`（只读）和 `readwrite`（读写）模式，通过中间件拦截写操作。
- **基本认证**：内置 HTTP Basic Authentication 支持。


### 使用方法

####  1. 下载

首先从[发布页面](https://github.com/hongenge/davgo/releases)下载适合您的操作系统和架构的最新程序。

#### 2. 配置 `config.yaml`

创建一个 `config.yaml` 文件，示例内容如下：
```yaml
port: 5344
services:
  dav1:
    root_dir: "./davroot1"
    username: "user"
    password: "pwd"
    mode: "readonly"
  dav2:
    root_dir: "./davroot2"
    username: "user"
    password: "pwd"
    mode: "readwrite"
```
- `port`：服务器监听端口。
- `services`：WebDAV 服务列表，键（如 `dav1`）决定访问路径（`/dav1/`）。
- `root_dir`：文件系统根目录。
- `username` 和 `password`：基本认证凭据。
- `mode`：`readonly` 或 `readwrite`，控制读写权限。

#### 3. 运行服务器
```bash
./davgo
```
服务器将在指定端口（默认 `5344`）启动。

#### 4. 可以用来挂载`WebDav`的软件

`Potplayer`，`kmplayer`，`RaiDrive`，`kodi`，`Nplayer`，ES文件管理器，nova魔改

