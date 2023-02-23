// Code generated by hertz generator.

package user

import (
	"context"

	"BiteDans.com/tiktok-backend/biz/dal/model"
	"BiteDans.com/tiktok-backend/biz/model/douyin/core/user"
	"BiteDans.com/tiktok-backend/pkg/constants"
	"BiteDans.com/tiktok-backend/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
)

// UserInfo .
// @router /douyin/user [GET]
func UserInfo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.DouyinUserRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(user.DouyinUserResponse)

	var curUserId uint

	if curUserId, err = utils.GetIdFromToken(req.Token); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Invalid token"
		resp.User = nil

		c.JSON(consts.StatusOK, resp)
		return
	}

	_user := new(model.User)

	if err = model.FindUserById(_user, uint(req.UserId)); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "User id does not exist"
		resp.User = nil
		c.JSON(consts.StatusOK, resp)
		return
	}

	respUser := &user.User{}
	var isFollow bool

	isFollow, err = model.GetFollowRelation(curUserId, uint(req.UserId))
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Failed to retrieve user info"
		resp.User = nil
		c.JSON(consts.StatusInternalServerError, resp)

		hlog.Errorf("Failed to retrieve user info: %v", err)
		return
	}

	userLikeReceivedCount, err := model.GetUserReceivedLikeCount(int64(_user.ID))
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Cannot get user like count"
		c.JSON(consts.StatusInternalServerError, resp)
		hlog.Errorf("Cannot get user like count for: %s", err.Error())
		return
	}

	userLikeCount, err := model.GetUserLikeCount(int64(_user.ID))
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Cannot get user received like count"
		c.JSON(consts.StatusInternalServerError, resp)
		hlog.Errorf("Cannot get user received like count for: %s", err.Error())
		return
	}

	userWorkCount, err := model.GetUserVideoCount(int64(_user.ID))
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Cannot get user work count"
		c.JSON(consts.StatusInternalServerError, resp)
		hlog.Errorf("Cannot get user work count for: %s", err.Error())
		return
	}

	userAvatar, err := model.FindUserAvatar(int64(_user.ID))
	if err != nil {
		userAvatar = constants.PROFILE_PIC_ADDR
	}

	userBackgroundImage, err := model.FindUserBackgroundImage(int64(_user.ID))
	if err != nil {
		userBackgroundImage = constants.BACKGROUND_PIC_ADDR
	}

	respUser.ID = int64(_user.ID)
	respUser.Name = _user.Username
	respUser.FollowCount = model.GetFollowCount(_user)
	respUser.FollowerCount = model.GetFollowerCount(_user)
	respUser.IsFollow = isFollow
	respUser.TotalFavorited = userLikeReceivedCount
	respUser.WorkCount = userWorkCount
	respUser.FavoriteCount = userLikeCount
	respUser.Signature = constants.SIGNATURE
	respUser.BackgroundImage = userBackgroundImage
	respUser.Avatar = userAvatar

	resp.StatusCode = 0
	resp.StatusMsg = "User info retrieved successfully"
	resp.User = respUser

	c.JSON(consts.StatusOK, resp)
}

// UserLogin .
// @router /douyin/user/login/ [POST]
func UserLogin(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.DouyinUserLoginRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(user.DouyinUserLoginResponse)
	user := new(model.User)

	user.Username = req.Username
	user.Password = req.Password

	var inputpassword = user.Password

	if err = model.FindUserByUsername(user, req.Username); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Username not found"
		c.JSON(consts.StatusOK, resp)
		return
	}

	if inputpassword != user.Password {
		resp.StatusCode = -1
		resp.StatusMsg = "Incorrect password"
		c.JSON(consts.StatusOK, resp)
		return
	}

	var token string

	if token, err = utils.GenerateJWT(user.ID); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Token generation failed"
		c.JSON(consts.StatusInternalServerError, resp)

		hlog.Errorf("Failed to generate token: %v", err)
		return
	}

	resp.StatusCode = 0
	resp.StatusMsg = "User logged in successfully"
	resp.UserId = int64(user.ID)
	resp.Token = token

	c.JSON(consts.StatusOK, resp)
}

// UserRegister .
// @router /douyin/user/register/ [POST]
func UserRegister(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.DouyinUserRegisterRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(user.DouyinUserRegisterResponse)
	user := new(model.User)

	if err = model.FindUserByUsername(user, req.Username); err == nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Username has been used"
		c.JSON(consts.StatusOK, resp)
		return
	}

	user.Username = req.Username
	user.Password = req.Password
	user.Avatar = "https://api.dicebear.com/5.x/bottts-neutral/png?seed=" + uuid.NewString()
	user.BackgroundImage = "https://api.dicebear.com/5.x/thumbs/png?seed=" + uuid.NewString()

	if err = model.CreateUser(user); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Failed to register user"
		c.JSON(consts.StatusInternalServerError, resp)

		hlog.Errorf("Failed to create user into database: %v", err)
		return
	}

	var token string

	if token, err = utils.GenerateJWT(user.ID); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Token generation failed"
		c.JSON(consts.StatusInternalServerError, resp)

		hlog.Errorf("Failed to generate token: %v", err)
		return
	}

	resp.StatusCode = 0
	resp.StatusMsg = "User registered successfully"
	resp.UserId = int64(user.ID)
	resp.Token = token

	c.JSON(consts.StatusOK, resp)
}
