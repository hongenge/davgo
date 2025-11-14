## DavGo

`DavGo` 是一个用 Go 语言实现的轻量级 WebDAV 服务器，支持动态配置多个 WebDAV 服务实例，每个实例可以独立设置根目录、认证信息和读写模式。

### 功能特性

* **动态配置**：通过 YAML 文件配置任意数量的 WebDAV 服务。
* **独立实例**：每个服务有独立的根目录、用户名和密码。
* **读写控制**：支持 `readonly`（只读）和 `readwrite`（读写）模式，通过中间件拦截写操作。
* **基本认证**：内置 HTTP Basic Authentication 支持。

### 使用方法

#### 1. 下载

首先从[发布页面](https://github.com/hongenge/davgo/releases)下载适合您的操作系统和架构的最新程序。

#### 2. 配置 `config.yaml`

创建一个 `config.yaml` 文件，示例内容如下：

```yaml
port: 5344
services:
  - name: webdav
    root_dir: "./davroot"
    username: "abc"
    password: "123"
    mode: "readwrite"

  - name: mydav
    root_dir: "./davroot"
    username: "abc"
    password: "123"
    mode: "readwrite"
```

* `port`：服务器监听端口。
* `services`：WebDAV 服务列表，每个实例用 `name` 指定访问路径（如 `/webdav/` 或 `/mydav/`）。
* `root_dir`：文件系统根目录。
* `username` 和 `password`：基本认证凭据。
* `mode`：`readonly` 或 `readwrite`，控制读写权限。

#### 2.1 访问路径示例

假设你的配置如下：

```yaml
services:
  - name: webdav
  - name: mydav
```

* 每个 `services` 项的 `name` 字段决定访问路径。
* 例如：

  * 第一个服务 `webdav` 可以通过 `http://服务器IP:5344/webdav/` 访问。
  * 第二个服务 `mydav` 可以通过 `http://服务器IP:5344/mydav/` 访问。
* 客户端挂载示例：

  * Windows 或 macOS：输入 `\\服务器IP@5344\webdav` 挂载。
  * RaiDrive 或 ES 文件管理器：URL 填 `http://服务器IP:5344/webdav/`。
* 每个实例都有独立的用户名和密码，需要按照 `username` 和 `password` 进行认证。

#### 3. 运行服务器

```bash
./davgo
```

服务器将在指定端口（默认 `5344`）启动。

#### 4. 可用挂载软件

`Potplayer`，`kmplayer`，`RaiDrive`，`kodi`，`Nplayer`，ES 文件管理器，Nova 魔改版。

### 反向代理

`nginx` 反向代理配置示例：

```nginx
location / {
  proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header X-Forwarded-Proto $scheme;
  proxy_set_header Host $http_host;
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header Range $http_range;
  proxy_set_header If-Range $http_if_range;
  proxy_redirect off;
  proxy_pass http://127.0.0.1:5344;
  # 上传文件最大尺寸
  client_max_body_size 20000m;
}
```
