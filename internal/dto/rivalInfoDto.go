package dto

import "github.com/Catizard/lampghost_wails/internal/entity"

type RivalInfoDto struct {
	ID        uint
	Name      string
	PlayCount int
	MainUser  bool

	DiffTableHeader *DiffTableHeaderDto
}

func NewRivalInfoDto(rival *entity.RivalInfo) *RivalInfoDto {
	return &RivalInfoDto{
		ID:        rival.ID,
		Name:      rival.Name,
		PlayCount: rival.PlayCount,
		MainUser:  rival.MainUser,
	}
}

func (rival *RivalInfoDto) Entity() *entity.RivalInfo {
	return &entity.RivalInfo{
		Name:      rival.Name,
		PlayCount: rival.PlayCount,
		MainUser:  rival.MainUser,
	}
}

func NewRivalInfoDtoWithDiffTable(rival *entity.RivalInfo, header *DiffTableHeaderDto) *RivalInfoDto {
	ret := NewRivalInfoDto(rival)
	ret.DiffTableHeader = header
	return ret
}
