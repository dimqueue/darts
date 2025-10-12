//go:build !dev
// +build !dev

package swagger

import "github.com/gin-gonic/gin"

// Swagger disabled in production
func SetupSwagger(router *gin.Engine) {

}
