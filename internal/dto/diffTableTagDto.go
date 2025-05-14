package dto

type DiffTableTagDto struct {
	TableName         string
	TableLevel        string
	TableSymbol       string
	TableTagColor     string
	TableTagTextColor string
}

func PushDownTag(log *RivalScoreLogDto) *DiffTableTagDto {
	return &DiffTableTagDto{
		TableName:         log.TableName,
		TableLevel:        log.TableLevel,
		TableSymbol:       log.TableSymbol,
		TableTagColor:     log.TableTagColor,
		TableTagTextColor: log.TableTagTextColor,
	}
}
