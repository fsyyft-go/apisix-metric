# 日志组件使用示例

本目录包含了日志组件的各种使用示例，帮助你快速了解和使用日志组件。每个示例都经过精心设计，展示了不同场景下的最佳实践。

## 示例目录

### 1. 基础示例（basic）

展示了日志组件的基本使用方法，包括：
- 初始化日志系统，配置输出目标
- 使用不同级别记录日志
- 通过结构化字段增强日志信息

```go
// 初始化日志系统，使用标准输出
if err := log.InitLogger(log.LogTypeStd, ""); err != nil {
    panic(err)
}

// 记录一条信息级别的日志
log.Info("这是一条信息日志。")
```

### 2. Logrus 示例（logrus）

展示了如何使用 Logrus 作为日志后端的高级特性，包括：
- 配置日志输出到文件系统
- 使用结构化字段记录上下文信息
- 集成错误处理和日志记录

```go
// 初始化 Logrus 日志系统
if err := log.InitLogger(log.LogTypeLogrus, "logs/app.log"); err != nil {
    panic(err)
}

// 使用结构化字段记录用户操作
log.WithFields(map[string]interface{}{
    "user": "admin",
    "action": "login",
}).Info("用户登录。")
```

## 运行示例

每个示例都可以独立运行，使用以下命令：

```bash
# 运行基础示例
cd basic && go run main.go

# 运行 Logrus 示例
cd logrus && go run main.go
```

## 注意事项

1. 在使用日志组件之前，必须先调用 InitLogger 进行初始化。
2. 根据应用场景选择合适的日志类型：
   - 开发环境建议使用标准输出（LogTypeStd）
   - 生产环境推荐使用 Logrus（LogTypeLogrus）
3. 合理使用日志级别：
   - Debug：仅在开发环境使用
   - Info：记录正常操作信息
   - Warn：记录潜在问题
   - Error：记录错误信息
   - Fatal：记录致命错误（会导致程序退出）
4. 使用结构化字段时的最佳实践：
   - 字段名使用小写字母
   - 避免使用特殊字符
   - 保持字段名的一致性
   - 合理组织字段信息 