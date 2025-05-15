package dto

type DiffTableTagDto struct {
	Md5               string
	TableName         string
	TableLevel        string
	TableSymbol       string
	TableTagColor     string
	TableTagTextColor string
}

func PushDownTag(log *RivalScoreLogDto) *DiffTableTagDto {
	return &DiffTableTagDto{
		Md5:               log.Md5,
		TableName:         log.TableName,
		TableLevel:        log.TableLevel,
		TableSymbol:       log.TableSymbol,
		TableTagColor:     log.TableTagColor,
		TableTagTextColor: log.TableTagTextColor,
	}
}
