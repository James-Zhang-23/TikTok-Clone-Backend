// Code generated by hertz generator.

package interaction

import (
	"context"

	"BiteDans.com/tiktok-backend/biz/dal/model"
	"BiteDans.com/tiktok-backend/biz/model/douyin/core/user"
	interaction "BiteDans.com/tiktok-backend/biz/model/douyin/extra/interaction"
	"BiteDans.com/tiktok-backend/pkg/constants"
	"BiteDans.com/tiktok-backend/pkg/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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

	var user_id uint
	if user_id, err = utils.GetIdFromToken(req.Token); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Invalid token"
		c.JSON(consts.StatusUnauthorized, resp)
		return
	}

	_user := new(model.User)
	if err = model.FindUserById(_user, user_id); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "User id does not exist"
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	_video := new(model.Video)
	if err = model.FindVideoById(_video, uint(req.VideoId)); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Video id does not exist"
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	if req.ActionType == constants.LIKE_VIDEO {
		like := new(model.Like)
		like.UserId = int64(user_id)
		like.VideoId = req.VideoId

		if err = model.IfUserLikedVideo(like); err != nil {
			if err = model.LikeVideo(like); err != nil {
				resp.StatusCode = -1
				resp.StatusMsg = "Failed to like the video"
				c.JSON(consts.StatusInternalServerError, resp)

				hlog.Error("Failed to create like record into database")
				return
			}

			if err = model.VideoLikeCountIncrease(req.VideoId); err != nil {
				resp.StatusCode = -1
				resp.StatusMsg = "Failed to increase liked video count"
				c.JSON(consts.StatusInternalServerError, resp)

				hlog.Error("Failed to increase liked video count in video database")
				return
			}

			resp.StatusCode = 0
			resp.StatusMsg = "Liked video successfully!"
			c.JSON(consts.StatusOK, resp)
			return
		} else {
			resp.StatusCode = 0
			resp.StatusMsg = "Already liked the video"

			c.JSON(consts.StatusOK, resp)
			return
		}

	} else if req.ActionType == constants.UNLIKE_VIDEO {
		like := new(model.Like)
		like.UserId = int64(user_id)
		like.VideoId = req.VideoId
		if err = model.IfUserLikedVideo(like); err != nil {
			resp.StatusCode = 0
			resp.StatusMsg = "Already unliked the video"

			c.JSON(consts.StatusOK, resp)
			return
		} else {
			if _video.FavoriteCount > 0 {
				if err = model.UnlikeVideo(like); err != nil {
					resp.StatusCode = -1
					resp.StatusMsg = "Failed to unlike the video"
					c.JSON(consts.StatusInternalServerError, resp)

					hlog.Error("Failed to delete like record into database")
					return
				}
				if err = model.VideoLikeCountDecrease(req.VideoId); err != nil {
					resp.StatusCode = -1
					resp.StatusMsg = "Failed to decrease liked video count"
					c.JSON(consts.StatusInternalServerError, resp)

					hlog.Error("Failed to decrease liked video count in video database")
					return
				}
				resp.StatusCode = 0
				resp.StatusMsg = "Unliked video successfully!"
				c.JSON(consts.StatusOK, resp)
				return
			} else {
				resp.StatusCode = -1
				resp.StatusMsg = "Like count is not positive number, not deductible"
				c.JSON(consts.StatusInternalServerError, resp)

				hlog.Error("Like count is not positive number, not deductible")
				return
			}
		}

	}
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

	if _, err = utils.GetIdFromToken(req.Token); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Invalid token"
		c.JSON(consts.StatusUnauthorized, resp)
		return
	}

	_user := new(model.User)
	if err = model.FindUserById(_user, uint(req.UserId)); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "User id does not exist"
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	var _liked_videos_id []*model.Like
	_liked_videos_id, err = model.FindLikedVideosIdByUserId(req.UserId)
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Failed to retrieve liked videos ID"
		resp.VideoList = nil
		c.JSON(consts.StatusInternalServerError, resp)
		hlog.Error("Failed to query liked videos ID by user id in database")
		return
	}

	resp.StatusCode = 0
	resp.StatusMsg = "Liked videos retrieved successfully"
	resp.VideoList = []*interaction.Video{}

	for _, liked_video_id := range _liked_videos_id {
		video := new(model.Video)
		if err = model.FindVideoById(video, uint(liked_video_id.VideoId)); err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "Liked video id does not exist"
			c.JSON(consts.StatusBadRequest, resp)
			return
		}

		the_user := new(model.User)
		if err = model.FindUserById(the_user, uint(video.AuthorId)); err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "Author of the video does not exist"
			resp.VideoList = nil
			c.JSON(consts.StatusInternalServerError, resp)
			hlog.Error("The user with author id of the video is not exist")
			return
		}

		format_user := &user.User{
			ID:            int64(the_user.ID),
			Name:          the_user.Username,
			FollowCount:   int64(len(the_user.Followings)),
			FollowerCount: int64(len(the_user.Followers)),
			IsFollow:      false,
		}

		the_like := new(model.Like)
		the_like.UserId = req.UserId
		the_like.VideoId = int64(video.ID)

		the_video := &interaction.Video{
			ID:            int64(video.ID),
			Author:        (*interaction.User)(format_user),
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    true,
			Title:         video.Title,
		}
		resp.VideoList = append(resp.VideoList, the_video)
	}

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

	var user_id uint
	if user_id, err = utils.GetIdFromToken(req.Token); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Invalid token"
		resp.Comment = nil

		c.JSON(consts.StatusUnauthorized, resp)
		return
	}

	_user := new(model.User)
	if err = model.FindUserById(_user, user_id); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "User id does not exist"
		resp.Comment = nil
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	_video := new(model.Video)
	if err = model.FindVideoById(_video, uint(req.VideoId)); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Video id does not exist"
		resp.Comment = nil
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	if req.ActionType != constants.POST_COMMENT && req.ActionType != constants.DELETE_COMMENT {
		resp.StatusCode = -1
		resp.StatusMsg = "Fail to get action type"
		resp.Comment = nil
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	if req.ActionType == constants.POST_COMMENT {
		comment := new(model.Comment)
		comment.UserId = int64(user_id)
		comment.VideoId = req.VideoId
		comment.Content = req.CommentText
		if err = model.CreateComment(comment); err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "Failed to create comment"
			resp.Comment = nil
			c.JSON(consts.StatusInternalServerError, resp)

			hlog.Error("Failed to create comment into database")
			return
		}
		resp.StatusCode = 0
		resp.StatusMsg = "comment on video successfully!"
		resp.Comment = &interaction.Comment{
			ID:         0,
			User:       nil,
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.Format("01-02"),
		}
		c.JSON(consts.StatusOK, resp)
		return
	}

	//Delete comment
	comment := new(model.Comment)
	if comment, err = model.FindCommentById(req.CommentId); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Comment id does not exist"
		resp.Comment = nil
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	if comment.UserId != (int64(user_id)) {
		resp.StatusCode = -1
		resp.StatusMsg = "You can not delete comment that does not belong to you"
		resp.Comment = nil
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	if err = model.DeleteComment(comment); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Fail to delete comment"
		resp.Comment = nil
		c.JSON(consts.StatusInternalServerError, resp)
		hlog.Error("Failed to delete comment in database")
		return
	}
	resp.StatusCode = 0
	resp.StatusMsg = "delete comment successfully"
	resp.Comment = nil

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

	var user_id uint
	if user_id, err = utils.GetIdFromToken(req.Token); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Invalid token"
		resp.CommentList = nil

		c.JSON(consts.StatusUnauthorized, resp)
		return
	}

	_user := new(model.User)
	if err = model.FindUserById(_user, user_id); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "User id does not exist"
		resp.CommentList = nil
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	_video := new(model.Video)
	if err = model.FindVideoById(_video, uint(req.VideoId)); err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Video id does not exist"
		resp.CommentList = nil
		c.JSON(consts.StatusBadRequest, resp)
		return
	}

	var _comments []*model.Comment
	_comments, err = model.FindCommentsByVideoId(req.VideoId)
	if err != nil {
		resp.StatusCode = -1
		resp.StatusMsg = "Failed to retrieve comments of the video"
		resp.CommentList = nil
		c.JSON(consts.StatusInternalServerError, resp)
		hlog.Error("Failed to query comments id with video id in database")
		return
	}

	resp.StatusCode = 0
	resp.StatusMsg = "Video comments retrieved successfully"
	resp.CommentList = []*interaction.Comment{}

	for _, comment := range _comments {
		the_user := new(model.User)
		if err = model.FindUserById(the_user, uint(comment.UserId)); err != nil {
			resp.StatusCode = -1
			resp.StatusMsg = "User id does not exist in comment"
			resp.CommentList = nil
			c.JSON(consts.StatusInternalServerError, resp)
			hlog.Error("The user with user id in comment is not exist")
			return
		}

		format_user := &user.User{
			ID:            int64(the_user.ID),
			Name:          the_user.Username,
			FollowCount:   int64(len(the_user.Followings)),
			FollowerCount: int64(len(the_user.Followers)),
			IsFollow:      false,
		}

		the_comment := &interaction.Comment{
			ID:         int64(comment.ID),
			User:       (*interaction.User)(format_user),
			Content:    comment.Content,
			CreateDate: comment.CreatedAt.Format("01-02"),
		}
		resp.CommentList = append(resp.CommentList, the_comment)
	}

	c.JSON(consts.StatusOK, resp)
}
