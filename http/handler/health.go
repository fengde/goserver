package handler

import (
	"net/http"
)

func Health(c *Context) {
	c.OutString(http.StatusOK, "success")
}
