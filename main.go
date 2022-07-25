package main

import (
	"encoding/json"
	"gmeheut/scrapping/scrap"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Search string `json:"search"`
}

func getSalary(c *gin.Context) {

	jsonData, err := c.GetRawData()

	if err != nil {
		return
	}
	var s Body
	err = json.Unmarshal(jsonData, &s)

	c.IndentedJSON(http.StatusOK, scrap.FetchData(s.Search))

}

func main() {

	// router := gin.Default()
	// router.POST("/salary", getSalary)

	// router.Run("localhost:4000")

	scrap.FetchData("nextjs")
}
