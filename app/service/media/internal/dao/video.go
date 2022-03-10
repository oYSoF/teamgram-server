// Copyright 2022 Teamgram Authors
//  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: teamgramio (teamgram.io@gmail.com)
//

package dao

import (
	"context"

	"github.com/gogo/protobuf/types"

	"github.com/teamgram/proto/mtproto"
	"github.com/teamgram/teamgram-server/app/service/media/internal/dal/dataobject"
)

func (m *Dao) GetVideoSizeList(ctx context.Context, sizeId int64) (sizes []*mtproto.VideoSize) {
	sizeDOList, _ := m.VideoSizesDAO.SelectListByVideoSizeId(ctx, sizeId)
	if len(sizeDOList) >= 0 {
		sizes = make([]*mtproto.VideoSize, 0, len(sizeDOList))
		for i := 0; i < len(sizeDOList); i++ {
			size := &sizeDOList[i]
			videoSize := mtproto.MakeTLVideoSize(&mtproto.VideoSize{
				Type:         size.SizeType,
				W:            size.Width,
				H:            size.Height,
				Size2:        size.FileSize,
				VideoStartTs: nil,
			}).To_VideoSize()
			if size.VideoStartTs > 0 {
				videoSize.VideoStartTs = &types.DoubleValue{Value: size.VideoStartTs}
			}
			sizes = append(sizes, videoSize)
		}
	}
	return
}

func (m *Dao) SaveVideoSizeV2(ctx context.Context, szId int64, szList []*mtproto.VideoSize) error {
	if len(szList) == 0 {
		return nil
	}

	for _, sz := range szList {
		szDO := &dataobject.VideoSizesDO{
			VideoSizeId: szId,
			SizeType:    sz.Type,
			// VolumeId:     sz.GetLocation().GetVolumeId(),
			// LocalId:      sz.GetLocation().GetLocalId(),
			// Secret:       sz.GetLocation().GetSecret(),
			Width:        sz.W,
			Height:       sz.H,
			FileSize:     sz.Size2,
			VideoStartTs: sz.GetVideoStartTs().GetValue(),
			FilePath:     "",
		}
		if _, _, err := m.VideoSizesDAO.Insert(ctx, szDO); err != nil {
			return err
		}
	}

	return nil
}