package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/astaxie/beego"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/daystram/cast/cast-be/config"
	"github.com/daystram/cast/cast-be/constants"
	data "github.com/daystram/cast/cast-be/datatransfers"
	"github.com/daystram/cast/cast-be/util"
)

func (m *module) CastList(variant string, count int, offset int, userID ...string) (videos []data.Video, err error) {
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

func (m *module) CreateVOD(upload data.VideoUpload, controller beego.Controller, userID string) (ID primitive.ObjectID, err error) {
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
	// Retrieve video and thumbnail
	video, _, _ := controller.GetFile("video")
	if _, err = m.s3.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(config.AppConfig.S3Bucket),
		Key:         aws.String(fmt.Sprintf("%s/%s/video.mp4", constants.VideoRootDir, ID.Hex())),
		Body:        video,
		ContentType: aws.String("video/mp4"),
	}); err != nil {
		_ = m.DeleteVideo(ID, userID)
		return primitive.ObjectID{}, fmt.Errorf("[CreateVOD] Failed saving video file. %+v", err)
	}
	var result bytes.Buffer
	thumbnail, _, _ := controller.GetFile("thumbnail")
	if result, err = util.NormalizeImage(thumbnail, constants.ThumbnailWidth, constants.ThumbnailHeight); err != nil {
		_ = m.DeleteVideo(ID, userID)
		return primitive.ObjectID{}, fmt.Errorf("[CreateVOD] Failed normalizing thumbnail image. %+v", err)
	}
	if _, err = m.s3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(config.AppConfig.S3Bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s.jpg", constants.ThumbnailRootDir, ID.Hex())),
		Body:   bytes.NewReader(result.Bytes()),
	}); err != nil {
		_ = m.DeleteVideo(ID, userID)
		return primitive.ObjectID{}, fmt.Errorf("[CreateVOD] Failed saving thumbnail image. %+v", err)
	}
	return
}

func (m *module) DeleteVideo(ID primitive.ObjectID, userID string) (err error) {
	var user data.User
	var video data.Video
	if user, err = m.db.userOrm.GetOneByID(userID); err != nil {
		return errors.New(fmt.Sprintf("[DeleteVideo] failed retrieving user %s detail. %+v", userID, err))
	}
	if video, err = m.db.videoOrm.GetOneByHash(ID.Hex()); err != nil {
		return errors.New(fmt.Sprintf("[DeleteVideo] failed retrieving video %s detail. %+v", userID, err))
	}
	if video.Type == constants.VideoTypeLive {
		return errors.New(fmt.Sprintf("[DeleteVideo] live videos cannot be deleted."))
	}
	if user.Username != video.Author.Username {
		return errors.New(fmt.Sprintf("[DeleteVideo] cannot delete others' video."))
	}
	// Delete files
	objects := []*s3.ObjectIdentifier{
		{Key: aws.String(fmt.Sprintf("%s/%s.ori", constants.ThumbnailRootDir, ID.Hex()))},
		{Key: aws.String(fmt.Sprintf("%s/%s.jpg", constants.ThumbnailRootDir, ID.Hex()))},
	}
	if list, err := m.s3.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(config.AppConfig.S3Bucket),
		Prefix: aws.String(fmt.Sprintf("%s/%s/", constants.VideoRootDir, ID.Hex())),
	}); err == nil {
		for _, item := range list.Contents {
			objects = append(objects, &s3.ObjectIdentifier{
				Key: item.Key,
			})
		}
	}
	_, _ = m.s3.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(config.AppConfig.S3Bucket),
		Delete: &s3.Delete{
			Objects: objects,
		},
	})
	return m.db.videoOrm.DeleteOneByID(ID)
}

func (m *module) UpdateVideo(video data.VideoEdit, controller beego.Controller, userID string) (err error) {
	if err = m.db.videoOrm.EditVideo(data.VideoInsert{
		Hash:        video.Hash,
		Title:       video.Title,
		Author:      userID,
		Description: video.Description,
		Tags:        video.Tags,
	}); err != nil {
		return errors.New(fmt.Sprintf("[UpdateVideo] error updating video. %+v", err))
	}
	// Retrieve thumbnail
	var thumbnail multipart.File
	if thumbnail, _, err = controller.GetFile("thumbnail"); err!= nil {
		if err == http.ErrMissingFile {
			return nil
		} else {
			return fmt.Errorf("[UpdateVideo] Failed retrieving thumbnail image. %+v\n", err)
		}
	}
	var result bytes.Buffer
	if result, err = util.NormalizeImage(thumbnail, constants.ThumbnailWidth, constants.ThumbnailHeight); err != nil {
		return fmt.Errorf("[UpdateVideo] Failed normalizing thumbnail image. %+v", err)
	}
	if _, err = m.s3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(config.AppConfig.S3Bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s.jpg", constants.ThumbnailRootDir, video.Hash)),
		Body:   bytes.NewReader(result.Bytes()),
	}); err != nil {
		return fmt.Errorf("[UpdateVideo] Failed saving thumbnail image. %+v", err)
	}
	return
}

func (m *module) CheckUniqueVideoTitle(title string) (err error) {
	return m.db.videoOrm.CheckUnique(title)
}

func (m *module) LikeVideo(userID string, hash string, like bool) (err error) {
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

func (m *module) Subscribe(userID string, username string, subscribe bool) (err error) {
	var user data.User
	var author data.User
	if user, err = m.db.userOrm.GetOneByID(userID); err != nil {
		return errors.New(fmt.Sprintf("[Subscribe] failed to get user by ID. %+v", err))
	}
	if author, err = m.db.userOrm.GetOneByUsername(username); err != nil {
		return errors.New(fmt.Sprintf("[Subscribe] failed to get author by username. %+v", err))
	}
	if subscribe {
		_, err = m.db.subscriptionOrm.InsertSubscription(data.Subscription{
			AuthorID:  author.ID,
			UserID:    userID,
			CreatedAt: time.Now(),
		})
		m.PushNotification(author.ID, data.NotificationOutgoing{
			Message:   fmt.Sprintf("%s just subscribed!", user.Username),
			Username:  user.Username,
			CreatedAt: time.Now(),
		})
	} else {
		err = m.db.subscriptionOrm.RemoveSubscriptionByAuthorIDUserID(author.ID, userID)
	}
	return
}

func (m *module) CheckUserSubscribes(authorID string, username string) (subscribed bool, err error) {
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

func (m *module) CommentVideo(userID string, hash, content string) (comment data.Comment, err error) {
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
			Username: user.Username,
			Name:     user.Name,
		},
		CreatedAt: now,
	}, nil
}
