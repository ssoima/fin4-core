package bysonhandlers

import (
	"fin4-core/server/datatype"
	"fin4-core/server/routermiddleware"
	"github.com/gin-gonic/gin"
)

//InjectHandlers injects asset handlers into the application
func InjectHandlers(sc datatype.ServiceContainer, rg *gin.RouterGroup) {
	authenticator := routermiddleware.SessionMustAuth()
	rg.GET("/test", authenticator, GetTokens(sc))
/*	rg.GET("/tokens/:tokenID", FindTokenForUser(sc))
	rg.POST("/create-token", authenticator, CreateToken(sc))
	rg.POST("/toggle-token-like", authenticator, ToggleTokenLike(sc))
	rg.POST("/create-claim", authenticator, CreateClaim(sc))
	rg.POST("/approve-claim", authenticator, ApproveClaim(sc))*/
}
