# Aoi (葵)

由 ChatGPT 驱动的 Ghost in the Shell

Aoi 是由 ChatGPT 驱动的基于人工智能的对话代理程序。使用 Aoi，可以在终端中与 AI 进行自然语言对话，它能够理解您的查询并执行适当的命令。

## 功能
 Aoi 可以用作 ChatGPT 的终端版本，此外，Aoi 还带有几个内置功能提高使用效率：

- `/code` - 生成代码片段并复制到剪贴板，例如 `/code go generate random numbers`
- `/db` - 生成 SQL 并在数据库上执行，例如 `/db {url} list tables`
- `/shell` - 生成 shell 命令并执行，例如 `/shell view listening ports`
- `/ssh` - 生成远程 shell 命令并执行，例如 `/ssh {host} view listening tcp ports`
- `/tldr` - 获取命令的 tl;dr 格式的解释
- `/trans` - 将文本翻译为指定语言

## 入门指南
可以从 GitHub 的[发布页面](https://github.com/shellfly/aoi/releases)下载 Aoi。或者，可以使用 Go 在系统上安装 Aoi：

```bash
go install github.com/shellfly/aoi@latest
```
将 OpenAI API 密钥设置为环境变量，然后运行 aoi 命令。

```bash

export OPENAI_API_KEY=<your_api_key>
aoi
```

## 演示
[![asciicast](https://asciinema.org/a/XjCGaMNf8Qp2nQ1UDlehjm5AN.svg)](https://asciinema.org/a/XjCGaMNf8Qp2nQ1UDlehjm5AN)
## 贡献
如果在使用 Aoi 时发现任何问题或有新功能的建议，请在 GitHub 存储库上创建问题或提交拉取请求。
