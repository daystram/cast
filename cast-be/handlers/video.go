package handlers

import (
	"errors"
	"fmt"
	"time"

	"gitlab.com/daystram/cast/cast-be/constants"
	data "gitlab.com/daystram/cast/cast-be/datatransfers"

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
	if videos, err = m.db().videoOrm.GetAllVODByAuthor(author.ID, count, offset); err != nil {
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
