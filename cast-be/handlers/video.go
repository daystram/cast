package handlers

import (
	"errors"
	"fmt"
	"time"

	"gitlab.com/daystram/cast/cast-be/constants"
	data "gitlab.com/daystram/cast/cast-be/datatransfers"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (m *module) VideoList(variant string, count int, offset int) (videos []data.Video, err error) {
	if videos, err = m.db().videoOrm.GetRecent(variant, count, offset); err != nil {
		return nil, errors.New(fmt.Sprintf("[VideoList] error retrieving recent videos. %+v", err))
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

func (m *module) DeleteVideo(ID primitive.ObjectID) (err error) {
	return m.db().videoOrm.DeleteOneByID(ID)
}
