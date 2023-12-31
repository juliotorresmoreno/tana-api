package mmlu

import (
	"errors"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juliotorresmoreno/tana-api/db"
	"github.com/juliotorresmoreno/tana-api/models"
	"github.com/juliotorresmoreno/tana-api/utils"
)

var tablename = models.Mmlu{}.TableName()

type MMLURouter struct {
}

func SetupAPIRoutes(r *gin.RouterGroup) {
	ai := &MMLURouter{}
	r.GET("", ai.find)
	r.GET("/:id", ai.findOne)
}

type Mmlus []Mmlu

func (h *MMLURouter) find(c *gin.Context) {
	mmlus := &Mmlus{}
	err := Find(mmlus)
	if err != nil {
		utils.Response(c, err)
		return
	}
	c.JSON(200, mmlus)
}

func Find(dest interface{}) error {
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer")
	}

	conn := db.DefaultClient

	if tx := conn.Table(tablename).Find(dest); tx.Error != nil {
		return utils.StatusInternalServerError
	}

	return nil
}

type Mmlu struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Feeling     string `json:"feeling"`
	PhotoURL    string `json:"photo_url"`
}

func (h *MMLURouter) findOne(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	mmlu := &Mmlu{}
	err := FindOne(id, mmlu)

	if err != nil {
		utils.Response(c, err)
		return
	}

	c.JSON(200, mmlu)
}

func FindOne(id int, dest interface{}) error {
	destType := reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer")
	}

	conn := db.DefaultClient
	if tx := conn.Table(tablename).First(dest, "id = ?", id); tx.Error != nil {
		return utils.StatusInternalServerError
	}

	return nil
}
