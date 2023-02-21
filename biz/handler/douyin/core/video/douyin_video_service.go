// Code generated by hertz generator.

package video

import (
	"context"
	"path/filepath"
	"strconv"
	"time"

	"os"

	"BiteDans.com/tiktok-backend/biz/dal/model"
	"BiteDans.com/tiktok-backend/biz/model/douyin/core/user"
	"BiteDans.com/tiktok-backend/biz/model/douyin/core/video"
	"BiteDans.com/tiktok-backend/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// VideoFeed .
// @router /douyin/feed [GET]
func VideoFeed(ctx context.Context, c *app.RequestContext) {
	// var err error
	var req video.DouyinVideoFeedRequest
	_ = c.Bind(&req)

	resp := new(video.DouyinVideoFeedResponse)

	now := time.Now()
	latestTime := req.LatestTime

	if latestTime == 0 {
		latestTime = now.UnixMilli()
	}
	unixTime := time.UnixMilli(latestTime)

	videos, err := model.FindLatestVideos(unixTime)
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Failed to retrieve videos"
		resp.VideoList = nil
		hlog.Errorf("Failed to find videos from the database with error: %s", err.Error())
		c.JSON(consts.StatusInternalServerError, resp)
		return
	}

	var curUserId uint
	curUserId, _ = utils.GetIdFromToken(req.Token)

	resp.StatusCode = 0
	resp.StatusMsg = "Publishing list info retrieved successfully"
	resp.VideoList = []*video.Video{}
	resp.NextTime = now.UnixMilli()

	for _, _video := range videos {
		author := new(model.User)
		author.ID = uint(_video.AuthorId)

		var isFollow bool

		if curUserId == 0 {
			isFollow = false
		} else {
			isFollow, _ = model.GetFollowRelation(curUserId, uint(_video.AuthorId))
		}

		userLikeReceivedCount, err := model.GetUserReceivedLikeCount(_video.AuthorId)
		if err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "Cannot get user like count"
			resp.VideoList = nil
			c.JSON(consts.StatusInternalServerError, resp)
			hlog.Errorf("Cannot get user like count for: %s", err.Error())
			return
		}

		userLikeCount, err := model.GetUserLikeCount(_video.AuthorId)
		if err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "Cannot get user received like count"
			resp.VideoList = nil
			c.JSON(consts.StatusInternalServerError, resp)
			hlog.Errorf("Cannot get user received like count for: %s", err.Error())
			return
		}

		userWorkCount, err := model.GetUserVideoCount(_video.AuthorId)
		if err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "Cannot get user work count"
			resp.VideoList = nil
			c.JSON(consts.StatusInternalServerError, resp)
			hlog.Errorf("Cannot get user work count for: %s", err.Error())
			return
		}

		the_user := &user.User{
			ID:            int64(_video.AuthorId),
			Name:          _video.AuthorUsername,
			FollowCount:   model.GetFollowCount(author),
			FollowerCount: model.GetFollowerCount(author),
			IsFollow:      isFollow,
			TotalFavorited: userLikeReceivedCount,
			WorkCount:	userWorkCount,
			FavoriteCount: userLikeCount,
		}

		like := new(model.Like)
		like.UserId = int64(curUserId)
		like.VideoId = int64(_video.ID)
		isFavorite := true
		err = model.IsVideoLiked(like)
		if err != nil {
			isFavorite = false
		}

		likeCount, err := model.GetVideoLikeCount(int64(_video.ID))
		if err != nil {
			hlog.Errorf("Failed to get video like count from database with error: %s", err.Error())
			c.String(consts.StatusInternalServerError, err.Error())
			return
		}

		commentCount, err := model.GetCommentCount(int64(_video.ID))
		if err != nil {
			hlog.Errorf("Failed to get video comment count from database with error: %s", err.Error())
			c.String(consts.StatusInternalServerError, err.Error())
			return
		}

		theVideo := &video.Video{
			ID:            int64(_video.ID),
			Author:        (*video.User)(the_user),
			PlayUrl:       _video.PlayUrl,
			CoverUrl:      _video.CoverUrl,
			FavoriteCount: likeCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
			Title:         _video.Title,
		}
		resp.VideoList = append(resp.VideoList, theVideo)
		if resp.NextTime > _video.CreatedAt.UnixMilli() {
			resp.NextTime = _video.CreatedAt.UnixMilli()
		}
	}

	c.JSON(consts.StatusOK, resp)
}

// VideoPublish .
// @router /douyin/publish/action/ [POST]
func VideoPublish(ctx context.Context, c *app.RequestContext) {
	var err error
	var req video.DouyinVideoPublishRequest
	resp := new(video.DouyinVideoPublishResponse)

	err = c.BindAndValidate(&req)
	if err != nil {
		resp.StatusMsg = err.Error()
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	userId, err := utils.GetIdFromToken(req.Token)
	if err != nil {
		resp.StatusMsg = "Invalid token"
		c.JSON(consts.StatusOK, resp)
		return
	}

	totalRows, err := model.GetVideoCount()
	if err != nil {
		hlog.Errorf("Failed to get the row count of videos with error: %s", err.Error())
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	form, _ := c.MultipartForm()
	file := form.File["data"][0]
	ext := filepath.Ext(file.Filename)

	// unique video name: number of total videos + 1.ext
	filename := strconv.Itoa(int(totalRows) + 1)
	fullFilename := filename + ext
	fullImagename := filename + ".png"

	_ = os.Mkdir("./files", 0755)
	c.SaveUploadedFile(file, "./files/"+fullFilename)

	_, err = utils.GetSnapshot("./files/"+fullFilename, "./files/"+filename, 1)

	if err != nil {
		hlog.Errorf("Failed to create snapshot image: %s", err.Error())
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	videoOutput, err := utils.UploadFile(fullFilename)
	if err != nil {
		hlog.Errorf("Unable to upload video to AWS with error: %s", err.Error())
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	coverOutput, err := utils.UploadFile(fullImagename)
	if err != nil {
		hlog.Errorf("Unable to upload image to AWS with error: %s", err.Error())
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	_ = os.Remove("./files/" + fullFilename)
	_ = os.Remove("./files/" + fullImagename)

	author := new(model.User)
	err = model.FindUserById(author, userId)
	if err != nil {
		hlog.Errorf("Failed to find user from database with error: %s", err.Error())
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	_video := new(model.Video)
	_video.AuthorId = int64(userId)
	_video.AuthorUsername = author.Username
	_video.PlayUrl = videoOutput
	_video.CoverUrl = coverOutput
	_video.Title = req.Title
	err = model.CreateVideo(_video)
	if err != nil {
		hlog.Errorf("Failed to save new video to the database with error: %s", err.Error())
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}

	resp.StatusCode = 0
	resp.StatusMsg = "Video published successfully"
	c.JSON(consts.StatusOK, resp)
}

// VideoPublishList .
// @router /douyin/publish/list/ [GET]
func VideoPublishList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req video.DouyinVideoPublishListRequest

	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(video.DouyinVideoPublishListResponse)
	userId, err := utils.GetIdFromToken(req.Token)
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Invalid token"
		resp.VideoList = nil
		c.JSON(consts.StatusOK, resp)
		return
	}

	_user := new(model.User)

	if err = model.FindUserById(_user, uint(req.UserId)); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "User id does not exist"
		resp.VideoList = nil
		c.JSON(consts.StatusOK, resp)
		return
	}

	videos, err := model.FindVideosByUserId(int64(_user.ID))
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Failed to retrieve videos"
		resp.VideoList = nil
		hlog.Errorf("Failed to find videos from the database with error: %s", err.Error())
		c.JSON(consts.StatusInternalServerError, resp)
		return
	}

	resp.StatusCode = 0
	resp.StatusMsg = "Publishing list info retrieved successfully"
	resp.VideoList = []*video.Video{}

	userLikeReceivedCount, err := model.GetUserReceivedLikeCount(int64(_user.ID))
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Cannot get user like count"
		resp.VideoList = nil
		c.JSON(consts.StatusInternalServerError, resp)
		hlog.Errorf("Cannot get user like count for: %s", err.Error())
		return
	}

	userLikeCount, err := model.GetUserLikeCount(int64(_user.ID))
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Cannot get user received like count"
		resp.VideoList = nil
		c.JSON(consts.StatusInternalServerError, resp)
		hlog.Errorf("Cannot get user received like count for: %s", err.Error())
		return
	}

	userWorkCount, err := model.GetUserVideoCount(int64(_user.ID))
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Cannot get user work count"
		resp.VideoList = nil
		c.JSON(consts.StatusInternalServerError, resp)
		hlog.Errorf("Cannot get user work count for: %s", err.Error())
		return
	}

	author := &user.User{
		ID:            int64(_user.ID),
		Name:          _user.Username,
		FollowCount:   model.GetFollowCount(_user),
		FollowerCount: model.GetFollowerCount(_user),
		IsFollow:      true,
		TotalFavorited: userLikeReceivedCount,
		WorkCount:	userWorkCount,
		FavoriteCount: userLikeCount,
	}

	isFollowingAuthor, err := model.GetFollowRelation(userId, _user.ID)
	if err != nil {
		hlog.Errorf("Failed to get user relation from database with error: %s", err.Error())
		c.String(consts.StatusInternalServerError, err.Error())
		return
	}
	author.IsFollow = isFollowingAuthor

	for _, _video := range videos {

		like := new(model.Like)
		like.UserId = int64(userId)
		like.VideoId = int64(_video.ID)
		isFavorite := true
		err = model.IsVideoLiked(like)
		if err != nil {
			isFavorite = false
		}

		likeCount, err := model.GetVideoLikeCount(int64(_video.ID))
		if err != nil {
			hlog.Errorf("Failed to get video like count from database with error: %s", err.Error())
			c.String(consts.StatusInternalServerError, err.Error())
			return
		}

		commentCount, err := model.GetCommentCount(int64(_video.ID))
		if err != nil {
			hlog.Errorf("Failed to get video comment count from database with error: %s", err.Error())
			c.String(consts.StatusInternalServerError, err.Error())
			return
		}

		theVideo := &video.Video{
			ID:            int64(_video.ID),
			Author:        (*video.User)(author),
			PlayUrl:       _video.PlayUrl,
			CoverUrl:      _video.CoverUrl,
			FavoriteCount: likeCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
			Title:         _video.Title,
		}
		resp.VideoList = append(resp.VideoList, theVideo)
	}

	c.JSON(consts.StatusOK, resp)
}
