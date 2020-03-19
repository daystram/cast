package handlers

import (
	"errors"
	"fmt"
	"gitlab.com/daystram/cast/cast-be/constants"

	data "gitlab.com/daystram/cast/cast-be/datatransfers"
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
