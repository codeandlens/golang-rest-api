package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GoHello(s *gin.Context) {
	s.IndentedJSON(http.StatusOK, "Hello Go !")
}
