// Code generated by hertz generator.

package follow

import (
	"context"
	"fmt"

	"BiteDans.com/tiktok-backend/biz/dal/model"
	"BiteDans.com/tiktok-backend/pkg/utils"

	follow "BiteDans.com/tiktok-backend/biz/model/douyin/extra/follow"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// FollowAction .
// @router /douyin/relation/action [POST]
func FollowAction(ctx context.Context, c *app.RequestContext) {
	var err error
	var req follow.DouyinRelationActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(follow.DouyinRelationActionResponse)
	curUserReq := req.Token
	curUserId, err1 := utils.GetIdFromToken(curUserReq)
	if err1 != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	curUser := new(model.User)
	if err := model.FindUserById(curUser, curUserId); err != nil {
		return
	}
	toUserId := req.ToUserId
	toUser := new(model.User)
	if err := model.FindUserById(toUser, uint(toUserId)); err != nil {
		return
	}
	actionType := req.ActionType
	if err := model.UserFollowAction(curUser, toUser, uint(actionType)); err != nil {
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// FollowList .
// @router /douyin/relation/follow/list/ [GET]
func FollowList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req follow.DouyinRelationFollowListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(follow.DouyinRelationFollowListResponse)
	/* get curUserId */
	var curUserId uint
	if curUserId, err = utils.GetIdFromToken(req.Token); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Invalid token"
		resp.UserList = nil

		c.JSON(consts.StatusUnauthorized, resp)
		return
	}

	/* get target User */
	_user := new(model.User)
	targetUserId := uint(req.UserId)
	if err = model.FindUserById(_user, targetUserId); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "User id does not exist"
		resp.UserList = nil
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	var uList []*model.User

	if err = model.GetFollowListByUser(&uList, _user); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "DB Find error"
		resp.UserList = nil
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	var respList []*follow.User
	for _, u := range uList {
		userResp := new(follow.User)
		if err := model.GetFollowInfoByIDs(curUserId, u.ID, userResp); err != nil {
			resp.StatusMsg = "wrong"
			return
		}
		respList = append(respList, userResp)
	}
	resp.StatusCode = 0
	resp.StatusMsg = "FollowList retrieved successfully"
	resp.UserList = respList

	c.JSON(consts.StatusOK, resp)
}

// FollowerList .
// @router /douyin/relation/follower/list/ [GET]
func FollowerList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req follow.DouyinRelationFollowerListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := new(follow.DouyinRelationFollowerListResponse)
	curUserReq := req.Token
	curUserId, err1 := utils.GetIdFromToken(curUserReq)
	if err1 != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	toUserId := req.UserId

	var idList []uint
	curUser := new(model.User)
	model.FindUserById(curUser, uint(toUserId))
	idList, _ = model.GetFollowersId(curUser)

	fmt.Println(idList)
	var followerList []*follow.User
	for i := 0; i < len(idList); i++ {
		user := new(follow.User)
		if err := model.GetFollowInfoByIDs(curUserId, idList[i], user); err != nil {
			resp.StatusMsg = "wrong"
			return
		}
		followerList = append(followerList, user)
	}

	fmt.Println(followerList)

	resp.StatusCode = 0
	resp.StatusMsg = "Get follower list successfully"
	resp.UserList = followerList

	c.JSON(consts.StatusOK, resp)

}
