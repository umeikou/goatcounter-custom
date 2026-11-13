[English](README.md) | [简体中文](README.zh-CN.md)

GoatCounter 是一个开源网站分析平台，可以作为（免费的）托管服务或自行托管的应用使用。它旨在提供易于使用、有实际意义且保护隐私的网站分析，作为 Google Analytics 或 Matomo 的替代方案。

有两种运行方式：使用 [goatcounter.com][www] 上的托管服务，或在自己的服务器上运行。源代码完全是开源/自由软件，可以不受限制地自行托管。

在线演示位于 [https://stats.arp242.net](https://stats.arp242.net)。

请考虑[提供资金支持][sponsor]。

[sponsor]: http://www.goatcounter.com/contribute
[www]: https://www.goatcounter.com


功能
----
- **注重隐私**；不使用唯一标识符跟踪用户，也不需要 GDPR 通知。可以精细地**控制收集哪些数据**。另请参阅[隐私政策][privacy]和 [GDPR 同意通知][gdpr]。

- **轻量且快速**；只会给网站增加约 3.5K 的额外数据。还提供无需 JavaScript 的“跟踪像素”选项，也可以从应用的中间件使用，或**从日志文件导入**。

- 使用不可识别身份的哈希，在不使用 Cookie 的情况下识别**独立访问**（[技术细节][sessions]）。

- 保留有用的统计信息，例如**浏览器**信息、**位置**和**屏幕尺寸**。跟踪**引荐网站**和**推广活动**。

- **简单易用**；如果 Google Analytics 和 Matomo 中大量您并不需要的选项与灵活性曾让您困惑，那么 GoatCounter 会让您耳目一新。

- **无障碍访问**是高优先级功能，界面能够很好地配合屏幕阅读器等辅助技术使用。

- 100% 致力于**开源**；您可以准确查看代码的行为并加以改进，也可以出于任何目的<strong>自行托管</strong>。

- **拥有自己的数据**；您始终可以导出所有数据，并且可以**随时取消**。

- 只需一个**脚本标签**即可集成到网站中：

      <script data-goatcounter="https://yoursite.goatcounter.com/count"
              async src="//gc.zgo.at/count.js"></script>

- JavaScript 集成适合大多数情况，但也可以使用**基于图片且无需 JavaScript 的跟踪器**、集成到**后端中间件**，或**解析日志文件**。

[privacy]: https://www.goatcounter.com/privacy
[gdpr]: https://www.goatcounter.com/gdpr
[sessions]: http://www.goatcounter.com/help/sessions


将数据导入 GoatCounter
----------------------
有三种方式：

1. 将 JavaScript 代码添加到网站；这是最简单、最常见的方法。详细文档位于 https://www.goatcounter.com/code

2. 使用 HTTP/REST API，例如从后端服务器中间件调用。详细文档位于 https://www.goatcounter.com/api#backend-integration

3. 解析 nginx、Apache、Caddy、CloudFront 或任何其他 HTTP 服务器或代理的日志文件。详细文档请参阅 `goatcounter help import`（自行托管版本和 goatcounter.com 均支持此方式）。

自行托管 GoatCounter
--------------------
[发布页面][releases]提供了适用于多个平台的二进制文件。这些文件采用静态编译，包含所需的一切，基本可以在任何环境中运行。唯一的依赖是一个用于存储 SQLite 数据库文件的位置或 PostgreSQL 连接。也可以使用 Docker，具体方法见下文。

[releases]: https://github.com/arp242/goatcounter/releases

### 运行
可以使用以下命令启动服务器：

    % goatcounter serve

这会在 `*:8080` 上启动服务器。默认使用位于 `./goatcounter-data/db.sqlite3` 的 SQLite 数据库；如果该数据库不存在，系统会创建它。

GoatCounter 同时支持 SQLite 和 PostgreSQL。SQLite 对大多数小型网站来说应该表现良好，而 PostgreSQL 的性能更好，尤其适用于较大的网站。主要瓶颈与其说是页面浏览量，不如说是这些浏览量分布在多少不同页面上（即使使用 SQLite，分布在 5 个页面上的 1,000 万次页面浏览也很快；但如果分布在 100 万个不同页面上，速度就会慢得多）。

要创建第一个站点，请使用 http://localhost:8080 上的向导，或使用 CLI：

    % goatcounter db create site -vhost=stats.example.com -user.email=me@example.com

该命令会要求输入密码；也可以通过 `-password` 在命令行中添加密码。如果使用的不是默认数据库，还必须在此处传递 `-db` 标志。

GoatCounter 内置 TLS 和自动 ACME 证书生成；在生产环境中运行时，可能需要使用类似以下的命令：

    % goatcounter serve -listen=:443 -tls=tls,rdr,acme

详情请参阅 `goatcounter help serve`。

### PostgreSQL
要使用 PostgreSQL，请通过自定义的 `-db` 标志运行 GoatCounter。例如：

    % goatcounter serve -db 'postgresql+dbname=goatcounter'
    % goatcounter serve -db 'postgresql+host=/run/postgresql dbname=goatcounter sslmode=disable'

其格式与 `psql` CLI 相同；也可以使用 `PG*` [环境变量](https://www.postgresql.org/docs/current/libpq-envars.html)：

    % PGDATABASE=goatcounter PGHOST=/run/postgresql goatcounter serve -db 'postgresql'

如果条件允许，数据库会自动创建；如果要为特定用户创建数据库，可以使用：

    % createuser --interactive --pwprompt goatcounter
    % createdb --owner goatcounter goatcounter

可以使用以下命令手动导入数据库结构：

    % goatcounter db schema-pgsql | psql --user=goatcounter --dbname=goatcounter

更多详情请参阅 `goatcounter help db` 和 [pq 文档][pq]。

[pq]: https://pkg.go.dev/github.com/lib/pq#NewConfig

### 使用 Docker 运行
GoatCounter 在 DockerHub 上以 [arp242/goatcounter] 提供。

运行新容器的示例：

    % docker run \
        -p 8080:8080 \
        -v goatcounter-data:/home/goatcounter/goatcounter-data \
        arp242/goatcounter

此命令使用命名卷。建议使用这种方式，因为该卷存储 SQLite 数据库和 ACME 证书（使用 ACME 时），而匿名卷很容易被意外删除。

要创建第一个站点，请使用 http://localhost:8080 上的向导，或使用 CLI：

    % docker exec -it [..] goatcounter db create site -vhost=stats.example.com -user.email=me@example.com

可以使用 `GOATCOUNTER_..` 环境变量设置选项。例如，要启用 TLS 和自动证书生成：

    % docker run \
        -p 80:80 \
        -p 443:443 \
        -v goatcounter-data:/home/goatcounter/goatcounter-data \
        -e GOATCOUNTER_LISTEN=:443 \
        -e GOATCOUNTER_TLS=tls,rdr,acme \
        arp242/goatcounter

设置 `GOATCOUNTER_DB` 以使用 PostgreSQL。例如：

    % docker run \
        -p 8080:8080 \
        -v goatcounter-data:/home/goatcounter/goatcounter-data \
        -e GOATCOUNTER_DB='postgresql+postgresql://goatcounter:goatcounter@postgres:5432/goatcounter?sslmode=disable' \
        arp242/goatcounter

所有选项请参阅 `goatcounter help serve`（或运行：`docker run --rm arp242/goatcounter help serve`）。

以上所有方式也应适用于 Podman。

还可以通过 `docker compose` 使用 compose.yaml 运行 GoatCounter。对于基本的 SQLite 配置：

    % docker compose up -d goatcounter-sqlite

或使用 PostgreSQL（也会从 compose.yaml 启动 PostgreSQL）：

    % docker compose up -d goatcounter-postgres

[arp242/goatcounter]: https://hub.docker.com/r/arp242/goatcounter

### 管理
服务器管理页面位于*设置 → 服务器管理*，可用于调试和测试某些功能。

状态 URL 为 `/status`，可以用于健康监测。

### 更新
更新时可能需要运行数据库迁移。使用 `goatcounter serve -automigrate` 可在启动时始终运行所有待处理的迁移。

使用 `goatcounter db migrate <file>` 或 `goatcounter db migrate all` 手动运行迁移。

使用 `goatcounter db migrate pending` 获取待处理迁移的列表，或使用 `goatcounter db migrate list` 显示所有迁移。

### 从源代码构建
需要 Go 1.21 或更新版本以及 C 编译器。如果使用 `CGO_ENABLED=0` 编译，则不需要 C 编译器，但只能使用 PostgreSQL。

可以使用以下命令从源代码构建：

    % git clone --branch=release-2.7 https://github.com/arp242/goatcounter
    % cd goatcounter
    % go build ./cmd/goatcounter

这会在当前目录中生成 `goatcounter` 二进制文件。

要使用最新的开发版本，请切换到 `main` 分支。

要构建完全静态链接的二进制文件：

    % go build -trimpath -ldflags='-s -w -extldflags=-static' \
        -tags='osusergo,netgo,sqlite_omit_load_extension' \
        ./cmd/goatcounter

建议像上述命令一样使用最新发布版本。main 分支应该比较稳定，但不作任何保证；有时在正式发布之前不会编写详细的发布/升级说明，因此可能遇到意外情况。

如果计划使用 PostgreSQL 且不使用 SQLite，可以在不使用 cgo 的情况下编译 goatcounter：

    % CGO_ENABLED=0 go build ./cmd/goatcounter

默认情况下，这会创建一个静态链接的二进制文件，无需额外标志。功能上没有太大差异，但构建会更轻松、更快速，因为不需要 C 编译器。

### 开发/测试
可以使用以下命令启动测试/开发服务器：

    % goatcounter serve -dev

`-dev` 标志会让一些开发工作更方便：应用会在重新编译时自动重启，模板和静态文件会直接从文件系统读取，并有其他一些细微调整。

有关如何运行开发服务器、编写补丁等内容的更多详情，请参阅 [.github/CONTRIBUTING.md](/.github/CONTRIBUTING.md)。
