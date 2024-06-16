package loader

import (
	"fmt"

	"github.com/Catizard/lampghost/internal/common/source"
	"github.com/Catizard/lampghost/internal/data"
	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/data/rival"
	"github.com/Catizard/lampghost/internal/data/score"
	"github.com/Catizard/lampghost/internal/service"
	"github.com/Catizard/lampghost/internal/sqlite"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
)

// TODO: OrajaDataLoader is obviously not a good name if there was multiple log loader for oraja exist.
// Note: orajaDataLoader is designed to be stateless, so we can expose it directly
var OrajaDataLoader RivalDataLoader = newOrajaDataLoader()

type orajaDataLoader struct {
}

func newOrajaDataLoader() *orajaDataLoader {
	return &orajaDataLoader{}
}

func (l *orajaDataLoader) Interest(r *rival.RivalInfo) bool {
	return r.SongDataPath.Valid && r.ScoreLogPath.Valid
}

// The default OrajaDataLoader loads 2 files: songdata.db and scorelog.db
func (l *orajaDataLoader) Load(r *rival.RivalInfo, filter null.Value[data.Filter]) ([]*score.CommonScoreLog, error) {
	if !l.Interest(r) {
		return nil, fmt.Errorf("[OrajaDataLoader] cannot load")
	}

	// 1) Loads scorelog
	// Directly read from scorelog table
	rawLogs, err := sqlite.DirectlyLoadTable[score.ScoreLog](r.ScoreLogPath.String, "scorelog")
	if err != nil {
		return nil, err
	}

	// Convert raw data to common form
	logs := make([]*score.CommonScoreLog, 0)
	for _, rawLog := range rawLogs {
		logs = append(logs, score.NewCommonScoreLogFromOraja(rawLog))
	}

	// 2) Loads songdata.db
	rawSongs, err := sqlite.DirectlyLoadTable[score.SongData](r.SongDataPath.String, "song")
	if err != nil {
		return nil, err
	}

	// 3) Assign md5 to log based on songdata.db
	// Note: Unassigned logs are ignored, because we cannot do anything further on it.
	sha256MapsToMD5 := make(map[string]string)
	md5MapsToSha256 := make(map[string]string)
	for _, v := range rawSongs {
		sha256MapsToMD5[v.Sha256] = v.Md5
		md5MapsToSha256[v.Md5] = v.Sha256
	}

	// Workaround: Courses are also related on md5
	// Note: We need to know a hash is whether generated from course or single song
	coursesHashSet := make(map[string]struct{})
	courses, _, err := service.CourseInfoService.FindCourseInfoList(difftable.CourseInfoFilter{})
	if err != nil {
		return nil, err
	}
	for _, course := range courses {
		sha256 := ""
		valid := true
		for _, v := range course.Md5 {
			if hash, ok := md5MapsToSha256[v]; ok {
				sha256 += hash
			} else {
				valid = false
				break
			}
		}
		if !valid {
			log.Warnf("[%s] is skipped because of lack of data", course.Name)
			continue
		}
		// Hack on sha256MapsToMD5
		sha256MapsToMD5[sha256] = course.Md5s
		// Mark it as a course hash
		coursesHashSet[sha256] = struct{}{}
	}

	finalLogs := make([]*score.CommonScoreLog, 0)
	for _, log := range logs {
		logHash := log.Sha256.ValueOrZero()
		if md5, ok := sha256MapsToMD5[logHash]; ok {
			// Warn: Any modification on 'log' wouldn't take affect
			log.Md5 = null.StringFrom(md5)
			if _, isCourse := coursesHashSet[logHash]; isCourse {
				log.LogType = source.Course
			} else {
				log.LogType = source.Song
			}
			finalLogs = append(finalLogs, log)
			// TODO: remove me!
			if !finalLogs[len(finalLogs)-1].Md5.Valid {
				panic("panic: no md5")
			}
		}
	}

	return finalLogs, nil
}
