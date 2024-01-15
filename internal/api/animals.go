package api

import (
	"github.com/KurbanowS/Animal/internal/app"
	"github.com/KurbanowS/Animal/internal/models"
	"github.com/gin-gonic/gin"
)

func AnimalRoutes(api *gin.RouterGroup) {
	animalRoutes := api.Group("/animals")
	{
		animalRoutes.GET("", AnimalList)
		animalRoutes.POST("", AnimalCreate)
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
