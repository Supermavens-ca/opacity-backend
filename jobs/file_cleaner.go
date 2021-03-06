package jobs

import (
	"time"

	"github.com/opacity/storage-node/models"
	"github.com/opacity/storage-node/utils"
)

type fileCleaner struct {
}

var olderThanOffset = -1 * 24 * time.Hour

func (f fileCleaner) Name() string {
	return "fileCleaner"
}

func (f fileCleaner) ScheduleInterval() string {
	return "@every 15m"
}

func (f fileCleaner) Run() {
	utils.SlackLog("running " + f.Name())

	files, err := models.DeleteUploadsOlderThan(time.Now().Add(olderThanOffset))
	utils.LogIfError(err, nil)

	if len(files) == 0 {
		return
	}

	for _, file := range files {
		utils.LogIfError(models.DeleteCompletedUploadIndexes(file.FileID), nil)
	}

	var ids []string
	for _, file := range files {
		ids = append(ids, models.GetFileMetadataKey(file.FileID))
	}
	err = utils.DeleteDefaultBucketObjects(ids)
	utils.LogIfError(err, nil)
}

func (f fileCleaner) Runnable() bool {
	return models.DB != nil
}
