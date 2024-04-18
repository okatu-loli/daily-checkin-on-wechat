# Daily-Checkin

Daily-Checkin-On-WeChat 是一个Go编写的自动化微信提醒工具，用于每天通知工作室成员值班情况。它利用了 `wechatbot-webhook` 接口来发送消息。

**依赖说明：** 本项目依赖于 [wechatbot-webhook](https://github.com/danni-cool/wechatbot-webhook) 服务。使用前，请确保已正确部署 `wechatbot-webhook`。详细的部署和配置指南，请访问 [项目页面](https://github.com/danni-cool/wechatbot-webhook)。

## 功能

- **自动提醒**：根据设定的时间自动向指定的微信群发送值班提醒。
- **灵活的配置文件**：可以通过配置文件自定义 webhook 地址、群名、值班时间表等。

## 配置文件示例

在使用之前，需要将 `config.yaml.example` 文件复制为 `config.yaml` 并根据实际情况进行修改：

```yaml
# Webhook 配置部分
webhook:
  # Webhook URL: 替换为实际的 webhook 接收地址
  url: "http://example.com/webhook"
  # 发送消息的微信群名或者房间名
  groupName: "ExampleGroup"
  # 是否为微信群聊，true为群聊，false为私聊
  isRoom: false

# 值班安排
dutySchedule:
  # 每个星期一的值班人员
  Monday: "Alice"
  # 每个星期二的值班人员
  Tuesday: "Bob"
  # 每个星期三的值班人员
  Wednesday: "Charlie"
  # 每个星期四的值班人员
  Thursday: "Dana"
  # 每个星期五的值班人员
  Friday: "Eve"
  # 每个星期六的值班人员
  Saturday: "Frank"
  # 每个星期日的值班人员
  Sunday: "Grace"

# 消息内容
messages:
  # 值班提醒消息内容
  duty: "You are on duty today!"

# 提醒时间表
schedule:
  # 提醒的小时（24小时制）
  hour: 9
  # 提醒的分钟
  minute: 0
  # 提醒的秒
  second: 0
```

## 部署

您可以选择以下两种方式来部署程序：

1. **直接运行**：
   ```bash
   go run .
   ```

2. **编译后运行**：
   ```bash
   go build -o daily_checkin
   ./daily_checkin
   ```

## 开始使用

1. 克隆仓库到本地：
   ```bash
   git clone https://github.com/your-username/daily-checkin.git
   ```
2. 配置 `config.yaml` 文件，根据你的需要修改示例配置。

3. 选择一种部署方法启动程序。

## 贡献

欢迎通过 Issue 提出问题或直接提交 Pull Request。

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 致谢

感谢 [danni-cool](https://github.com/danni-cool) 提供的 [wechatbot-webhook](https://github.com/danni-cool/wechatbot-webhook) 接口，使得消息推送实现变得简单。
