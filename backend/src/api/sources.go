package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getSources(c *gin.Context) {
	user, _ := c.Get("user")
	fmt.Println("sources request received by " + fmt.Sprintf("%v", user))

	sources := make([]string, 0)

	//for _, source := range util.Sources {
	//	sources = append(sources, source.String())
	//}

	c.JSON(http.StatusOK, sources)
}

func putSource(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func patchSource(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func deleteSource(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
