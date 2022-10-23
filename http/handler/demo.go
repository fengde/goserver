package handler

type DemoRequest struct {
	User string `form:"user"` // uri参数
	Name string `json:"name"` // json参数
	Age  string `form:"age"`  // form参数
}

func Demo(c *Context, r *DemoRequest) error {
	return nil
}
