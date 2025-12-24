# BingoHR

## 简介
一个包含Go API的示例项目，`api/Dockerfile` 会构建并运行服务（默认暴露端口 8000）。

## CI/CD：使用 GitHub Actions（非Docker）部署到 Azure VM
本仓库已添加工作流以实现：
- 在CI中为 Linux/amd64 构建 Go 二进制
- 打包二进制、`conf/app.ini` 与字体文件
- 通过SSH上传到Azure VM，创建systemd服务并启动
- 支持多环境：推送 `main` 分支 → dev 环境，推送 `tag` (v*) → pro 环境

### 前置条件
- Azure VM（开发和生产各一台或相同）可通过 SSH 访问，SSH 用户具备 `sudo` 权限
- VM 无需安装 Docker；只需可写入 `/opt/bingohr` 并允许 `systemd`
- 数据库与Redis可用（默认 `conf/app.ini` 配置为本机：MySQL `127.0.0.1:3306`、Redis `127.0.0.1:6379`）

### 环境配置
在 GitHub 仓库 Settings → Environments 中创建两个环境：`dev` 和 `pro`，各环境配置其对应的 Secrets：

**dev 环境 Secrets：**
- `AZURE_VM_HOST`：开发VM的公网IP或DNS
- `AZURE_VM_USER`：SSH用户名
- `AZURE_VM_SSH_KEY`：SSH私钥内容（PEM）
- `AZURE_VM_SSH_PORT`（可选）：SSH端口，默认22

**pro 环境 Secrets：**
- `AZURE_VM_HOST`：生产VM的公网IP或DNS
- `AZURE_VM_USER`：SSH用户名
- `AZURE_VM_SSH_KEY`：SSH私钥内容（PEM）
- `AZURE_VM_SSH_PORT`（可选）：SSH端口，默认22

**仓库级 Secrets（两个环境共用）：**
在 GitHub 仓库 Settings → Secrets and variables → Actions → Repository secrets 中配置：
- `TEAMS_WEBHOOK_URL`：Microsoft Teams Incoming Webhook URL，用于部署成功或失败的通知推送

### 工作流说明
工作流文件：`.github/workflows/deploy.yml`

**触发条件与环境选择：**
- 推送到 `main` 分支 → 自动使用 `dev` 环境部署
- 推送 tag 标签（格式 `v*`） → 自动使用 `pro` 环境部署
- 手动触发（Workflow Dispatch） → 从下拉菜单选择 `dev` 或 `pro` 环境

**工作流步骤：**
1. 构建：Linux/amd64 二进制编译
2. 打包：二进制+配置+字体为 tar.gz
3. 部署：SSH连接到对应环境的VM，解压到 `/opt/bingohr/api`，启动 systemd 服务，监听端口 `8000`
4. 通知：部署后独立的通知Job（成功或失败都会向该环境的Teams Webhook推送卡片），包含状态/版本/提交信息/运行耗时

### 快速开始
**部署到dev环境：**
```bash
git push origin main
```

**部署到pro环境：**
```bash
git tag v1.0.0
git push origin v1.0.0
```

### 本地验证
```bash
cd api
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bingohr-api .
./bingohr-api # 本地运行（需MySQL/Redis可用）
```

