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
func SendRobotMsg(key string, msg *RobotMsg) error {
	_, err := grequests.Post(baseURL+key, &grequests.RequestOptions{
		JSON: msg,
	})
	if err != nil {
		return err
	}
	return nil
}

// =============== 三个通知机器人 出错，紧急和普通 ===================

// wError 程序出错通知
func wError(args ...interface{}) {
	text := fmt.Sprint(args...)
	err := SendRobotMsg(errorKey, &RobotMsg{
		MsgType: "text",
		Text: &MsgText{
			Content: text,
		},
	})
	if err != nil {
		log.WithError(err).Error("向企业微信发送出错通知出错")
	} else {
		log.Errorf("发送出错通知成功：%s", text)
	}
}

// wWarn 紧急通知
func wWarn(args ...interface{}) {
	text := fmt.Sprint(args...)
	err := SendRobotMsg(warnKey, &RobotMsg{
		MsgType: "text",
		Text: &MsgText{
			Content: text,
		},
	})
	if err != nil {
		log.WithError(err).Error("向企业微信发送紧急通知出错")
	} else {
		log.Warnf("发送紧急通知成功：%s", text)
	}
}

// wWarnMD 紧急通知 Markdown
func wWarnMD(lines []string) {
	content := strings.Join(lines, "\n")
	err := SendRobotMsg(warnKey, &RobotMsg{
		MsgType: "markdown",
		MarkDown: &MsgMarkdown{
			Content: content,
		},
	})
	if err != nil {
		log.WithError(err).Error("向企业微信发送紧急通知出错")
	} else {
		log.Warnf("发送紧急通知成功：%s", content)
	}
}

// wInfo 一般通知
func wInfo(args ...interface{}) {
	text := fmt.Sprint(args...)
	err := SendRobotMsg(infoKey, &RobotMsg{
		MsgType: "text",
		Text: &MsgText{
			Content: text,
		},
	})
	if err != nil {
		log.WithError(err).Error("向企业微信发送一般通知出错")
	} else {
		log.Infof("发送一般通知成功：%s", text)
	}
}
