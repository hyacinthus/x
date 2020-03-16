package xim

import (
	"fmt"
	"strings"

	"github.com/levigross/grequests"
)

// 企业微信机器人基础链接
const baseURL = `https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=`

// RobotMsg 群机器人消息
type RobotMsg struct {
	MsgType  string       `json:"msgtype"`
	Text     *MsgText     `json:"text,omitempty"`
	MarkDown *MsgMarkdown `json:"markdown,omitempty"`
}

// MsgText 文本消息部分
type MsgText struct {
	Content string `json:"content"`
}

// MsgMarkdown MD消息
type MsgMarkdown struct {
	Content string `json:"content"`
}

// SendRobotMsg 发送机器人消息
func SendRobotMsg(key, tp, content string) {
	var msg = &RobotMsg{
		MsgType: tp,
	}
	switch tp {
	case "text":
		msg.Text = &MsgText{Content: content}
	case "markdown":
		msg.MarkDown = &MsgMarkdown{Content: content}
	}
	_, err := grequests.Post(baseURL+key, &grequests.RequestOptions{
		JSON: msg,
	})
	if err != nil {
		log.WithError(err).Errorf("向企业微信发送通知出错:%s", content)
		return
	}
	log.Infof("向企业微信发送通知成功：%s", content)
}

// SendRobotMarkdown 向机器人发送 Markdown 通知
func SendRobotMarkdown(key string, lines []string) {
	// 一次最多4096字 控制一下
	var buffer = make([]string, 0)
	var count, times int
	for _, line := range lines {
		if len(line) > 4000 {
			// 单行不能超过4000字
			log.Error("发送微信消息时单行超过4000字，取消发送，打印如下。")
			log.Error(strings.Join(lines, "\n"))
			return
		}
		if count+len(line) > 4000 {
			// 加上这行就超长了，赶紧把之前的打印一下，重新开始积累
			SendRobotMsg(key, "markdown", strings.Join(buffer, "\n"))
			buffer = make([]string, 0)
			count = 0
			times++
			if times >= 5 {
				// 已经发了5次了，太长了，不继续发了
				log.Error("发送微信消息时满容量发送了5次还未发完，取消后续发送，打印如下。")
				log.Error(strings.Join(lines, "\n"))
				return
			}
		}
		buffer = append(buffer, line)
		count += len(line)
	}
	// 最后发送 buffer 中残留的最后一部分
	SendRobotMsg(key, "markdown", strings.Join(buffer, "\n"))
}

// =============== 三个通知机器人 出错，紧急和普通 ===================

// wError 程序出错通知
func wError(args ...interface{}) {
	SendRobotMsg(errorKey, "text", fmt.Sprint(args...))
}

// wWarn 紧急通知
func wWarn(args ...interface{}) {
	SendRobotMsg(warnKey, "text", fmt.Sprint(args...))
}

// wInfo 一般通知
func wInfo(args ...interface{}) {
	SendRobotMsg(infoKey, "text", fmt.Sprint(args...))
}

// wErrorMD 出错通知 Markdown
func wErrorMD(lines []string) {
	SendRobotMarkdown(errorKey, lines)
}

// wWarnMD 紧急通知 Markdown
func wWarnMD(lines []string) {
	SendRobotMarkdown(warnKey, lines)
}

// wInfoMD 一般通知 Markdown
func wInfoMD(lines []string) {
	SendRobotMarkdown(infoKey, lines)
}
