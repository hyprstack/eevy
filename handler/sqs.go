package handler

import (
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/sqs"

	"github.com/hevnly/eevy/event"
	"github.com/hevnly/eevy/handler/config"
)

// Places a message on an AWS SQS when relavent event is triggered
type Sqs struct {
	HandlerBase

	Config config.Sqs
}

// Satisfies the Listener interface and places the event on an AWS SQS
func (this *Sqs) Exec(evt event.Event) {

	this.Log.Handler(this, &evt)

	url := magicString(this.Config.GetUrl(), evt)
	reg := magicString(this.Config.GetRegion(), evt)
	msg := magicString(this.Config.GetMessage(), evt)
	svc := sqs.New(&aws.Config{Region: reg})
	params := &sqs.SendMessageInput{
		MessageBody: aws.String(msg),
		QueueURL:    aws.String(url),
	}
	_, err := svc.SendMessage(params)

	if err != nil {
		this.Log.Error(err.Error())
		return
	}
}

func (this *Sqs) GetType() string {

	return this.GetConfig().GetType()
}

func (this *Sqs) GetConfig() config.Handler {
	return this.Config
}
