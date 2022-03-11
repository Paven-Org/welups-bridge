package userRouter

import "github.com/gin-gonic/gin"

func Config(router gin.IRouter) {
	gr := router.Group("/u")
	gr.POST("/login", loginHandler)
	gr.POST("/logout" /*,authenMW*/, logoutHandler)
	gr.POST("/passwd" /*,authenMW*/, passwdHandler)
}

func loginHandler(c *gin.Context) {

}

func logoutHandler(c *gin.Context) {

}

func passwdHandler(c *gin.Context) {

}
