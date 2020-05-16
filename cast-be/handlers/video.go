package handlers

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gitlab.com/daystram/cast/cast-be/config"
	"gitlab.com/daystram/cast/cast-be/constants"
	data "gitlab.com/daystram/cast/cast-be/datatransfers"
	"gitlab.com/daystram/cast/cast-be/util"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *module) CastList(variant string, count int, offset int, userID ...primitive.ObjectID) (videos []data.Video, err error) {
	switch {
	case variant == constants.VideoListTrending:
		if videos, err = m.db.videoOrm.GetTrending(count, offset); err != nil {
			return nil, errors.New(fmt.Sprintf("[CastList] error retrieving trending videos. %+v", err))
		}
	case variant == constants.VideoTypeLive || variant == constants.VideoTypeVOD:
		if videos, err = m.db.videoOrm.GetRecent(variant, count, offset); err != nil {
			return nil, errors.New(fmt.Sprintf("[CastList] error retrieving recent videos. %+v", err))
		}
	case variant == constants.VideoListLiked:
		if len(userID) != 1 {
			return nil, errors.New(fmt.Sprintf("[CastList] userID not provided"))
		}
		if videos, err = m.db.videoOrm.GetLiked(userID[0], count, offset); err != nil {
			return nil, errors.New(fmt.Sprintf("[CastList] error retrieving liked videos. %+v", err))
		}
	case variant == constants.VideoListSubscribed:
		if len(userID) != 1 {
			return nil, errors.New(fmt.Sprintf("[CastList] userID not provided"))
		}
		if videos, err = m.db.videoOrm.GetSubscribed(userID[0], count, offset); err != nil {
			return nil, errors.New(fmt.Sprintf("[CastList] error retrieving subscribed videos. %+v", err))
		}
	default:
		return nil, errors.New(fmt.Sprintf("[CastList] invalid variant %s. %+v", variant, err))
	}
	return
}
func (m *module) AuthorList(username string, count, offset int) (videos []data.Video, err error) {
	var author data.User
	if author, err = m.db.userOrm.GetOneByUsername(username); err != nil {
		return nil, errors.New(fmt.Sprintf("[AuthorList] author not found. %+v", err))
	}
	if videos, err = m.db.videoOrm.GetAllVODByAuthorPaginated(author.ID, count, offset); err != nil {
		return nil, errors.New(fmt.Sprintf("[AuthorList] error retrieving VODs. %+v", err))
	}
	return
}

func (m *module) SearchVideo(query string, _ []string, count, offset int) (videos []data.Video, err error) {
	if videos, err = m.db.videoOrm.Search(query, count, offset); err != nil {
		return nil, errors.New(fmt.Sprintf("[SearchVideo] error searching videos. %+v", err))
	}
	return
}

func (m *module) VideoDetails(hash string) (video data.Video, err error) {
	var comments []data.Comment
	if video, err = m.db.videoOrm.GetOneByHash(hash); err != nil {
		return data.Video{}, errors.New(fmt.Sprintf("[VideoDetails] video with hash %s not found. %+v", hash, err))
	}
	if comments, err = m.db.commentOrm.GetAllByHash(hash); err != nil {
		return data.Video{}, errors.New(fmt.Sprintf("[VideoDetails] failed getting comment list for %s. %+v", hash, err))
	}
	video.Views++
	video.Comments = comments
	if video.Type == constants.VideoTypeVOD {
		if err = m.db.videoOrm.IncrementViews(hash); err != nil {
			return data.Video{}, errors.New(fmt.Sprintf("[VideoDetails] failed incrementing views of %s. %+v", hash, err))
		}
	}
	return
}

func (m *module) CreateVOD(upload data.VideoUpload, userID primitive.ObjectID) (ID primitive.ObjectID, err error) {
	if ID, err = m.db.videoOrm.InsertVideo(data.VideoInsert{
		Type:        constants.VideoTypeVOD,
		Title:       upload.Title,
		Author:      userID,
		Description: upload.Description,
		Tags:        upload.Tags,
		Views:       0,
		Duration:    0,
		IsLive:      true,
		Resolutions: 0,
		CreatedAt:   time.Now(),
	}); err != nil {
		return primitive.ObjectID{}, errors.New(fmt.Sprintf("[CreateVOD] error inserting video. %+v", err))
	}
	return
}

func (m *module) DeleteVideo(ID, userID primitive.ObjectID) (err error) {
	var user data.User
	var video data.Video
	if user, err = m.db.userOrm.GetOneByID(userID); err != nil {
		return errors.New(fmt.Sprintf("[DeleteVideo] failed retrieving user %s detail. %+v", userID.Hex(), err))
	}
	if video, err = m.db.videoOrm.GetOneByHash(ID.Hex()); err != nil {
		return errors.New(fmt.Sprintf("[DeleteVideo] failed retrieving video %s detail. %+v", userID.Hex(), err))
	}
	if video.Type == constants.VideoTypeLive {
		return errors.New(fmt.Sprintf("[DeleteVideo] live videos cannot be deleted."))
	}
	if user.Username != video.Author.Username {
		return errors.New(fmt.Sprintf("[DeleteVideo] cannot delete others' video."))
	}
	_ = os.RemoveAll(fmt.Sprintf(fmt.Sprintf("%s/%s", config.AppConfig.UploadsDirectory, ID.Hex())))
	_ = os.Remove(fmt.Sprintf("%s/%s/%s.ori", config.AppConfig.UploadsDirectory, constants.ThumbnailRootDir, ID.Hex()))
	_ = os.Remove(fmt.Sprintf("%s/%s/%s.jpg", config.AppConfig.UploadsDirectory, constants.ThumbnailRootDir, ID.Hex()))
	return m.db.videoOrm.DeleteOneByID(ID)
}

func (m *module) UpdateVideo(video data.VideoEdit, userID primitive.ObjectID) (err error) {
	if err = m.db.videoOrm.EditVideo(data.VideoInsert{
		Hash:        video.Hash,
		Title:       video.Title,
		Author:      userID,
		Description: video.Description,
		Tags:        video.Tags,
	}); err != nil {
		return errors.New(fmt.Sprintf("[UpdateVideo] error updating video. %+v", err))
	}
	return
}

func (m *module) CheckUniqueVideoTitle(title string) (err error) {
	return m.db.videoOrm.CheckUnique(title)
}

func (m *module) NormalizeThumbnail(hash string) (err error) {
	return util.NormalizeImage(constants.ThumbnailRootDir, hash, constants.ThumbnailWidth, constants.ThumbnailHeight)
}

func (m *module) LikeVideo(userID primitive.ObjectID, hash string, like bool) (err error) {
	if like {
		_, err = m.db.likeOrm.InsertLike(data.Like{
			Hash:      hash,
			Author:    userID,
			CreatedAt: time.Now(),
		})
	} else {
		err = m.db.likeOrm.RemoveLikeByUserIDHash(userID, hash)
	}
	return
}

func (m *module) CheckUserLikes(hash, username string) (liked bool, err error) {
	var user data.User
	if user, err = m.db.userOrm.GetOneByUsername(username); err != nil {
		return false, errors.New(fmt.Sprintf("[CheckUserLikes] failed to get user by username. %+v", err))
	}
	if _, err = m.db.likeOrm.GetOneByUserIDHash(user.ID, hash); err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, errors.New(fmt.Sprintf("[CheckUserLikes] failed to fetch likes by user. %+v", err))
	}
	return true, nil
}

func (m *module) Subscribe(userID primitive.ObjectID, username string, subscribe bool) (err error) {
	var author data.User
	if author, err = m.db.userOrm.GetOneByUsername(username); err != nil {
		return errors.New(fmt.Sprintf("[Subscribe] failed to get author by username. %+v", err))
	}
	if subscribe {
		_, err = m.db.subscriptionOrm.InsertSubscription(data.Subscription{
			AuthorID:  author.ID,
			UserID:    userID,
			CreatedAt: time.Now(),
		})
	} else {
		err = m.db.subscriptionOrm.RemoveSubscriptionByAuthorIDUserID(author.ID, userID)
	}
	return
}

func (m *module) CheckUserSubscribes(authorID primitive.ObjectID, username string) (subscribed bool, err error) {
	var user data.User
	if user, err = m.db.userOrm.GetOneByUsername(username); err != nil {
		return false, errors.New(fmt.Sprintf("[CheckUserSubscribes] failed to get user by username. %+v", err))
	}
	if _, err = m.db.subscriptionOrm.GetOneByAuthorIDUserID(authorID, user.ID); err != nil {
		if err != mongo.ErrNoDocuments {
			return false, errors.New(fmt.Sprintf("[CheckUserSubscribes] failed to subscription info. %+v", err))
		}
		return false, nil
	}
	return true, nil
}

func (m *module) CommentVideo(userID primitive.ObjectID, hash, content string) (comment data.Comment, err error) {
	var commentID primitive.ObjectID
	var user data.User
	now := time.Now()
	if commentID, err = m.db.commentOrm.InsertComment(data.CommentInsert{
		Hash:      hash,
		Content:   content,
		Author:    userID,
		CreatedAt: now,
	}); err != nil {
		return data.Comment{}, errors.New(fmt.Sprintf("[CommentVideo] failed to insert comment. %+v", err))
	}
	if user, err = m.db.userOrm.GetOneByID(userID); err != nil {
		return data.Comment{}, errors.New(fmt.Sprintf("[CommentVideo] failed to retrieve user info. %+v", err))
	}
	return data.Comment{
		ID:      commentID,
		Hash:    hash,
		Content: content,
		Author: data.UserItem{
			Name:     user.Name,
			Username: user.Username,
		},
		CreatedAt: now,
	}, nil
}
