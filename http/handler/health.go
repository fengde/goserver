package handler

import "net/http"

func Health(c *Context) {
	c.String(http.StatusOK, "i am ok")
}
