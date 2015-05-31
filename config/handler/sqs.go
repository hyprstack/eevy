package handler

type Sqs struct {
	HandlerBase
}

func (this *Sqs) GetRegion() string {
	return this.getString("region")
}

func (this *Sqs) GetUrl() string {
	return this.getString("url")
}
