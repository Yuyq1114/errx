# errx — Go 轻量级错误堆栈扩展库

errx 用于增强 Go 的错误处理能力，使其在保持原生 error 类型的同时，具备 堆栈信息、链式包装、可格式化输出 等特性，让错误定位更接近 panic 时的边界可见性。

在多级函数调用场景中，原生错误仅携带字符串信息，不包含调用链位置。errx 的目标是让你在 不依赖手动打印日志 的情况下，也能获取完整的错误调用链，极大提升调试效率。

✨ 特性

支持自定义错误并自动记录错误发生位置

支持错误包装（链式扩展），记录完整调用链

支持类似 panic 的堆栈输出

支持格式化为字符串，可直接输出到标准日志组件中

API 极简，保持与 Go 原生错误处理风格一致

## 1. 安装

go get github.com/yuyq1114/errx

## 2. 使用方法

errx 提供四个核心方法：

### （1）New —— 创建新错误（带堆栈）

用于手动构造一个错误，同时记录创建位置。例如：数据库查询返回 0 行时，你希望将其视为错误。

return errx.New("no rows found")


等价于创建一个 error + 记录 stack 信息。

### （2）Warp —— 包装错误（追加堆栈）

用于在多级函数中继续包装已有错误，自动记录当前函数位置，无需重复写日志。

原逻辑：

if err != nil {
log.Error(err)
return err
}


使用 errx 后：

if err != nil {
return errx.Warp(err, "query user failed")
}


顶层统一处理即可。

### （3）Print —— 直接打印错误堆栈

适用于 CLI 或 panic 捕获后的打印。

if err != nil {
errx.Print(err)
}


输出类似：

[ERRX] query user failed
at service/user.go:42
at dao/user.go:30
at database/sql.go:91

### （4）Format —— 格式化为字符串（用于日志系统）

如果你仍使用旧的日志打印方式，可以将错误堆栈格式化成字符串：

log.Println(errx.Format(err))


获得可完整输出的堆栈字符串。


## 3. 推荐的错误处理模式

✔ 在底层创建错误
✔ 在中间层 Warp 错误
✔ 在最顶层统一打印 / 格式化错误
✔ 不随处打印日志（避免重复打印）

## 4. 适用场景

多层调用的业务服务（Web/微服务）

数据库操作、RPC 调用等需要明确调用链

复杂调试场景，需要比 error 更多信息

想接近 panic 可见性的 error 处理

## 5. 许可证

MIT License。