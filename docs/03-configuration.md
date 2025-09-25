# 03. 配置指南

hwcctl 支持多种配置方式，满足不同使用场景的需求。

## 配置优先级

配置参数按以下优先级应用（高优先级覆盖低优先级）：

1. **命令行参数** (最高优先级)
2. **环境变量**
3. **配置文件** (最低优先级)

## 配置文件

### 配置文件位置

```
~/.hwcctl/config
```

### 配置文件格式

```yaml
default:
  # 华为云认证信息
  access_key_id: "your-access-key-id"
  secret_access_key: "your-secret-access-key"
  region: "cn-north-4"
  domain_id: "your-domain-id" # CDN 服务必需

  # 输出设置
  output: "table" # table, json, yaml

  # 重试设置
  enable_retry: false # 是否启用重试
  max_retries: 3 # 最大重试次数
```

### 创建配置文件

#### 方式一：交互式配置

```bash
hwcctl configure
```

#### 方式二：手动创建

```bash
# 创建配置目录
mkdir -p ~/.hwcctl

# 创建配置文件
cat > ~/.hwcctl/config << 'EOF'
default:
    access_key_id: "your-access-key-id"
    secret_access_key: "your-secret-access-key"
    region: "cn-north-4"
    domain_id: "your-domain-id"
    output: "table"
    enable_retry: false
    max_retries: 0
EOF
```

## 环境变量

```bash
# 华为云认证
export HUAWEICLOUD_ACCESS_KEY="your-access-key"
export HUAWEICLOUD_SECRET_KEY="your-secret-key"
export HUAWEICLOUD_REGION="cn-north-4"
export HUAWEICLOUD_DOMAIN_ID="your-domain-id"
```

## 命令行参数

```bash
hwcctl --access-key-id "key" --secret-access-key "secret" --region "cn-north-4" --domain-id "domain" cdn refresh --urls "https://example.com/file.jpg"
```

## 获取认证信息

### 1. Access Key 和 Secret Key

1. 登录华为云控制台
2. 点击右上角用户名 → "我的凭证"
3. 在"访问密钥"页面创建或查看 Access Key

### 2. Domain ID（CDN 必需）

1. 登录华为云控制台
2. 点击右上角用户名 → "我的凭证"
3. 在"API 凭证"页面查看"账号 ID"（即 Domain ID）

### 3. 区域 (Region)

常用华为云区域：

- `cn-north-1` - 华北-北京一
- `cn-north-4` - 华北-北京四
- `cn-east-2` - 华东-上海二
- `cn-east-3` - 华东-上海一
- `cn-south-1` - 华南-广州

## 配置验证

```bash
# 验证配置是否正确
hwcctl version

# 使用调试模式查看配置加载过程
hwcctl --debug cdn --help
```

## 安全建议

1. **文件权限**：确保配置文件权限为 `600`

   ```bash
   chmod 600 ~/.hwcctl/config
   ```

2. **环境变量**：在脚本中使用环境变量而非硬编码
3. **密钥轮换**：定期轮换访问密钥
4. **最小权限**：为应用创建专用的 IAM 用户，仅授予必要权限

## 多环境配置

目前 hwcctl 支持单一配置 profile (`default`)。如需多环境支持，可以：

1. **使用不同配置文件**

   ```bash
   # 开发环境
   hwcctl --profile dev cdn refresh --urls "..."

   # 生产环境
   hwcctl --profile prod cdn refresh --urls "..."
   ```

2. **使用环境变量切换**

   ```bash
   # 开发环境
   source ~/.hwcctl/dev.env
   hwcctl cdn refresh --urls "..."

   # 生产环境
   source ~/.hwcctl/prod.env
   hwcctl cdn refresh --urls "..."
   ```

## 故障排查

如果配置有问题，请查看 [故障排查指南](./08-troubleshooting.md)。
