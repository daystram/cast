package handlers

import (
	"errors"
	"fmt"
	"image"
	"os"
	"time"

	"gitlab.com/daystram/cast/cast-be/config"
	"gitlab.com/daystram/cast/cast-be/constants"
	data "gitlab.com/daystram/cast/cast-be/datatransfers"

	"github.com/disintegration/imaging"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *module) FreshList(variant string, count int, offset int) (videos []data.Video, err error) {
	if videos, err = m.db().videoOrm.GetRecent(variant, count, offset); err != nil {
		return nil, errors.New(fmt.Sprintf("[FreshList] error retrieving recent videos. %+v", err))
	}
	return
}
func (m *module) AuthorList(username string, count, offset int) (videos []data.Video, err error) {
	var author data.User
	if author, err = m.db().userOrm.GetOneByUsername(username); err != nil {
		return nil, errors.New(fmt.Sprintf("[AuthorList] author not found. %+v", err))
	}
	if videos, err = m.db().videoOrm.GetAllVODByAuthorPaginated(author.ID, count, offset); err != nil {
		return nil, errors.New(fmt.Sprintf("[AuthorList] error retrieving VODs. %+v", err))
	}
	return
}

func (m *module) Search(query string, tags []string) (videos []data.Video, err error) {
	return nil, nil
}

func (m *module) VideoDetails(hash string) (video data.Video, err error) {
	if video, err = m.db().videoOrm.GetOneByHash(hash); err != nil {
		return data.Video{}, errors.New(fmt.Sprintf("[VideoDetails] video with hash %s not found. %+v", hash, err))
	}
	video.Views++
	if video.Type == constants.VideoTypeVOD {
		if err = m.db().videoOrm.IncrementViews(hash); err != nil {
			return data.Video{}, errors.New(fmt.Sprintf("[VideoDetails] failed incrementing views of %s. %+v", hash, err))
		}
	}
	return
}

func (m *module) CreateVOD(upload data.VideoUpload, userID primitive.ObjectID) (ID primitive.ObjectID, err error) {
	if ID, err = m.db().videoOrm.InsertVideo(data.VideoInsert{
		Type:        constants.VideoTypeVOD,
		Title:       upload.Title,
		Author:      userID,
		Description: upload.Description,
		Views:       0,
		Duration:    0,
		IsLive:      false,
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
	if user, err = m.db().userOrm.GetOneByID(userID); err != nil {
		return errors.New(fmt.Sprintf("[DeleteVideo] failed retrieving user %s detail. %+v", userID.Hex(), err))
	}
	if video, err = m.db().videoOrm.GetOneByHash(ID.Hex()); err != nil {
		return errors.New(fmt.Sprintf("[DeleteVideo] failed retrieving video %s detail. %+v", userID.Hex(), err))
	}
	if video.Type == constants.VideoTypeLive {
		return errors.New(fmt.Sprintf("[DeleteVideo] live videos cannot be deleted."))
	}
	if user.Username != video.Author.Username {
		return errors.New(fmt.Sprintf("[DeleteVideo] cannot delete others' video."))
	}
	_ = os.RemoveAll(fmt.Sprintf(fmt.Sprintf("%s/%s", config.AppConfig.UploadsDirectory, ID.Hex())))
	_ = os.Remove(fmt.Sprintf("%s/thumbnail/%s.jpg", config.AppConfig.UploadsDirectory, ID.Hex()))
	return m.db().videoOrm.DeleteOneByID(ID)
}

func (m *module) UpdateVideo(video data.VideoEdit, userID primitive.ObjectID) (err error) {
	if err = m.db().videoOrm.EditVideo(data.VideoInsert{
		Hash:        video.Hash,
		Title:       video.Title,
		Author:      userID,
		Description: video.Description,
	}); err != nil {
		return errors.New(fmt.Sprintf("[UpdateVideo] error updating video. %+v", err))
	}
	return
}

func (m *module) CheckUniqueVideoTitle(title string) (err error) {
	return m.db().videoOrm.CheckUnique(title)
}

func (m *module) NormalizeThumbnail(ID primitive.ObjectID) (err error) {
	var reader *os.File
	if reader, err = os.Open(fmt.Sprintf("%s/thumbnail/%s.ori", config.AppConfig.UploadsDirectory, ID.Hex())); err != nil {
		return errors.New(fmt.Sprintf("[NormalizeThumbnail] failed to read original image. %+v", err))
	}
	original, _, err := image.Decode(reader)
	if err != nil {
		return
	}
	normalized := imaging.Fill(original, constants.ThumbnailWidth, constants.ThumbnailHeight, imaging.Center, imaging.Lanczos)
	if err = imaging.Save(normalized, fmt.Sprintf("%s/thumbnail/%s.jpg", config.AppConfig.UploadsDirectory, ID.Hex())); err != nil {
		return errors.New(fmt.Sprintf("[NormalizeThumbnail] failed to normalize image. %+v", err))
	}
	reader.Close()
	_ = os.Remove(fmt.Sprintf("%s/thumbnail/%s.ori", config.AppConfig.UploadsDirectory, ID.Hex()))
	return
}
