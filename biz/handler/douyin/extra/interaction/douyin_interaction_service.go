// Code generated by hertz generator.

package interaction

import (
	"context"

	interaction "BiteDans.com/tiktok-backend/biz/model/douyin/extra/interaction"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// FavoriteInteraction .
// @router /douyin/favorite/action/ [POST]
func FavoriteInteraction(ctx context.Context, c *app.RequestContext) {
	var err error
	var req interaction.DouyinFavoriteActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(interaction.DouyinFavoriteActionResponse)

	c.JSON(consts.StatusOK, resp)
}

// FavoriteList .
// @router /douyin/favorite/list/ [GET]
func FavoriteList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req interaction.DouyinFavoriteListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(interaction.DouyinFavoriteListResponse)

	c.JSON(consts.StatusOK, resp)
}

// CommentInteraction .
// @router /douyin/comment/action/ [POST]
func CommentInteraction(ctx context.Context, c *app.RequestContext) {
	var err error
	var req interaction.DouyinCommentActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(interaction.DouyinCommentActionResponse)

	c.JSON(consts.StatusOK, resp)
}

// CommentList .
// @router /douyin/comment/list/ [GET]
func CommentList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req interaction.DouyinCommentListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(interaction.DouyinCommentListResponse)

	c.JSON(consts.StatusOK, resp)
}
