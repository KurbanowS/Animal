package api

import (
	"strconv"

	"github.com/KurbanowS/Animal/internal/app"
	"github.com/KurbanowS/Animal/internal/models"
	"github.com/gin-gonic/gin"
)

func AnimalRoutes(api *gin.RouterGroup) {
	animalRoutes := api.Group("/animals")
	{
		animalRoutes.GET("", AnimalList)
		animalRoutes.POST("", AnimalCreate)
		animalRoutes.PUT(":id", AnimalUpdate)
		animalRoutes.DELETE("", AnimalDelete)
	}
}

func AnimalList(c *gin.Context) {
	r := models.AnimalFilterRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" || errKey != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
		return
	}
	if r.Limit == nil {
		r.Limit = new(int)
		*r.Limit = 30
	}
	if r.Offset == nil {
		r.Offset = new(int)
		*r.Offset = 0
	}
	animals, total, err := app.AnimalList(r)
	if err != nil {
		handleError(c, err)
		return
	}
	Success(c, gin.H{
		"animals": animals,
		"total":   total,
	})
}

func AnimalCreate(c *gin.Context) {
	r := models.AnimalRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" || errKey != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
	}
	animal, err := app.AnimalCreate(r)
	if err != nil {
		handleError(c, err)
		return
	}
	Success(c, gin.H{
		"animal": animal,
	})
}

func AnimalUpdate(c *gin.Context) {
	r := models.AnimalRequest{}
	if errMsg, errKey := BindAndValidate(c, &r); errMsg != "" {
		handleError(c, app.NewAppError(errMsg, errKey, ""))
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		handleError(c, err)
		return
	}
	idp := uint(id)
	r.ID = &idp

	if id == 0 {
		handleError(c, app.ErrRequired.SetKey("id"))
		return
	}
	animal, err := app.AnimalUpdate(r)
	if err != nil {
		handleError(c, err)
		return
	}
	Success(c, gin.H{
		"animal": animal,
	})
}

func AnimalDelete(c *gin.Context) {
	var ids []string = c.QueryArray("ids")
	if len(ids) == 0 {
		handleError(c, app.ErrRequired.SetKey("ids"))
		return
	}
	animals, err := app.AnimalDelete(ids)
	if err != nil {
		handleError(c, err)
		return
	}
	Success(c, gin.H{
		"animals": animals,
	})
}
