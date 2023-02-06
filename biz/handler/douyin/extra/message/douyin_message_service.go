// Code generated by hertz generator.

package message

import (
	"BiteDans.com/tiktok-backend/biz/dal/model"
	"BiteDans.com/tiktok-backend/pkg/utils"
	"context"

	message "BiteDans.com/tiktok-backend/biz/model/douyin/extra/message"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// MessageSend .
// @router /douyin/message/action/ [POST]
func MessageSend(ctx context.Context, c *app.RequestContext) {
	var err error
	var req message.DouyinRelationActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(message.DouyinRelationActionResponse)

	senderuser := new(model.User)
	receiveruser := new(model.User)
	var userID uint
	if userID, err = utils.GetIdFromToken(req.Token); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Invalid token"

		c.JSON(consts.StatusUnauthorized, resp)
		return
	}

	if err = model.FindUserById(senderuser, userID); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Sender user id does not exist"

		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	if err = model.FindUserById(receiveruser, uint(req.ToUserId)); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Receiver user id does not exist"

		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	if req.ActionType == 1 {
		_message := new(model.Message)
		_message.ToUserId = req.ToUserId
		_message.FromUserId = int64(userID)
		_message.Content = req.Content
		if err = model.CreateMessage(_message); err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "Message creation failed"

			c.JSON(consts.StatusInternalServerError, resp)
			return
		}
	} else {
		resp.StatusCode = -1
		resp.StatusMsg = "Action type not defined"

		c.JSON(consts.StatusBadRequest, resp)
		return

	}
	resp.StatusCode = 0
	resp.StatusMsg = "Message sent successfully"
	c.JSON(consts.StatusOK, resp)
}

// MessageHistory .
// @router /douyin/message/chat/ [GET]
func MessageHistory(ctx context.Context, c *app.RequestContext) {
	var err error
	var req message.DouyinMessageChatRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(message.DouyinMessageChatResponse)

	var user_id uint

	if user_id, err = utils.GetIdFromToken(req.Token); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Invalid token"
		resp.MessageList = nil

		c.JSON(consts.StatusUnauthorized, resp)
		return
	}

	_user := new(model.User)

	if err = model.FindUserById(_user, user_id); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "This User id does not exist"
		resp.MessageList = nil
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	receiver := new(model.User)

	if err = model.FindUserById(receiver, uint(req.ToUserId)); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Receiver id does not exist"
		resp.MessageList = nil
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	var _messages []*model.Message

	if _messages, err = model.FindMessageBySenderandReceiverId(_messages, user_id, uint(req.ToUserId)); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Fail to retrieve messages from this user to specified user"
		resp.MessageList = nil
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	resp.StatusCode = 0
	resp.StatusMsg = "Messages retrieved successfully"
	resp.MessageList = []*message.Message{}

	for _, _message := range _messages {
		create_time := _message.CreatedAt.String()
		the_message := &message.Message{
			ID:         int64(_message.ID),
			ToUserId:   _message.ToUserId,
			FromUserId: _message.FromUserId,
			Content:    _message.Content,
			CreateTime: create_time,
		}
		resp.MessageList = append(resp.MessageList, the_message)
	}

	c.JSON(consts.StatusOK, resp)
}