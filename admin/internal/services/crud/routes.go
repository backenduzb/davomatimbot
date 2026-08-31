package crud

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterCRUDRoutes[T any, S any, R any](router *gin.RouterGroup, path string, db *gorm.DB, middleware gin.HandlerFunc) {
    ctrl := NewCRUDController[T, S, R](db)
    
    if middleware == nil {
        middleware = func(c *gin.Context) {
            c.Next()
        }
    }

    router.GET(path, middleware, ctrl.List)
    router.POST(path, middleware, ctrl.Create)
    router.GET(path+"/:id", middleware, ctrl.Retrieve)
    router.PUT(path+"/:id", middleware, ctrl.Update)
    router.PATCH(path+"/:id", middleware, ctrl.PartialUpdate)
    router.DELETE(path+"/:id", middleware, ctrl.Delete)
}