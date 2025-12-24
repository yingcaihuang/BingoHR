# BingoHR

## 简介
一个包含Go API的示例项目，`api/Dockerfile` 会构建并运行服务（默认暴露端口 8000）。

## CI/CD：使用 GitHub Actions（非Docker）部署到 Azure VM
本仓库已添加工作流以实现：
- 在CI中为 Linux/amd64 构建 Go 二进制
- 打包二进制、`conf/app.ini` 与字体文件
- 通过SSH上传到Azure VM，创建systemd服务并启动

### 前置条件
- Azure VM 可通过 SSH 访问，SSH 用户具备 `sudo` 权限
- VM 无需安装 Docker；只需可写入 `/opt/bingohr` 并允许 `systemd`
- 数据库与Redis可用（默认 `conf/app.ini` 配置为本机：MySQL `127.0.0.1:3306`、Redis `127.0.0.1:6379`）

### 必需的仓库 Secrets
在 GitHub 仓库 Settings → Secrets → Actions 配置：
- `AZURE_VM_HOST`：VM 公网 IP 或 DNS
- `AZURE_VM_USER`：SSH 用户名
- `AZURE_VM_SSH_KEY`：上述用户的私钥内容（PEM 格式）
（不再需要GHCR相关的Secrets）

### 工作流说明
工作流文件：`.github/workflows/deploy.yml`
- 触发条件：推送到 `main` 或手动触发（Workflow Dispatch）
- 构建：`api` 目录下构建Linux二进制（`bingohr-api`）
- 打包：二进制+配置+字体为 `bingohr-api-linux-amd64.tar.gz`
- 部署：解压到 `/opt/bingohr/api`，生成并重启 `bingohr` systemd 服务，默认监听端口 `8000`

### 本地验证
```bash
cd api
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bingohr-api .
./bingohr-api # 本地运行（需MySQL/Redis可用）
```

