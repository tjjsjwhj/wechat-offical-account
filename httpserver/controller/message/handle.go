/*
 * @Description: 消息处理
 * @Date: 2021-03-27 17:11:48
 * @LastEditors: wanghaijie01
 * @LastEditTime: 2021-04-18 19:42:35
 */

package message

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"github.com/wechat-official-account/model/qabot"
)

// EventHandler 事件消息处理
func EventHandler(c *gin.Context, event string) (interface{}, error) {
	return nil, nil
}

// MsgHandler 被动回复用户消息
func MsgHandler(c *gin.Context, msgType string) (interface{}, error) {
	switch msgType {
	case "text":
		return TextHandler(c)
	}
	return nil, nil
}

// TextHandler 文本消息处理
func TextHandler(c *gin.Context) (interface{}, error) {
	clientMsg := &ClientTextMsg{}
	if err := c.ShouldBindBodyWith(clientMsg, binding.XML); err != nil {
		log.Warn("bind body err: ", err)
		return nil, err
	}

	// 先看配置关键词回复
	if resp, ok := keyReplay(clientMsg.Content, clientMsg.FromUserName, clientMsg.ToUserName); ok {
		return resp, nil
	}

	// 暂时只回复文本消息
	serverMsg := &ServerTextMsg{
		ToUserName:   clientMsg.FromUserName,
		FromUserName: clientMsg.ToUserName,
		MsgType:      "text",
		CreateTime:   time.Now().Unix(),
	}

	// 先获取关键词回复
	content, err := qabot.QA(c, clientMsg.Content)
	if err != nil {
		content = "success"
	}
	serverMsg.Content = content.(string)
	return serverMsg, nil
}
