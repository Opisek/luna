package handlers

import (
	"luna-backend/api/internal/util"
	"luna-backend/files"

	"github.com/gin-gonic/gin"
)

func GetFile(c *gin.Context) {
	u := util.GetUtil(c)

	fileId, tr := util.GetId(c, "file")
	if tr != nil {
		u.Error(tr)
		return
	}

	file := files.GetDatabaseFile(fileId)

	_, tr = file.GetContent(u.Tx.Queries())
	if tr != nil {
		u.Error(tr)
		return
	}

	u.ResponseWithFile(file)
}
