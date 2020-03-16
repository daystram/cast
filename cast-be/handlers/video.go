package handlers

import (
	"errors"
	"fmt"

	data "gitlab.com/daystram/cast/cast-be/datatransfers"
)

func (m *module) GetVideo(variant string, count int, offset int) (videos []data.Video, err error) {
	if videos, err = m.db().videoOrm.GetRecent(variant, count, offset); err != nil {
		return nil, errors.New(fmt.Sprintf("[VideoHandler] error retrieving recent videos. %+v", err))
	}
	return
}

func (m *module) Search(query string, tags []string) (videos []data.Video, err error) {
	return nil, nil
}

func (m *module) VODDetails(hash string) (video data.Video, err error) {
	return data.Video{}, nil
}

func (m *module) LiveDetails(username string) (video data.Video, err error) {
	return data.Video{}, nil
}
