package handler

type Lambda struct {
	HandlerBase
}

func (this *Lambda) GetFunction() string {
	return this.getString("function")
}

func (this *Lambda) GetRegion() string {
	return this.getString("region")
}
