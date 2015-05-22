package listener

import (
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/sqs"

	"github.com/hevnly/eevy/config"
	"github.com/hevnly/eevy/event"
)

type Sqs struct {
	ListenerBase

	Url    string
	Region string
}

func (this *Sqs) Exec(evt event.Event) {

	gLog.Debug("SQS %s on event %s", this.Url, evt.Id)
	msg := this.GetMessage(evt)

	svc := sqs.New(&aws.Config{Region: this.Region})
	params := &sqs.SendMessageInput{
		MessageBody: aws.String(msg),
		QueueURL:    aws.String(this.Url),
	}
	_, err := svc.SendMessage(params)

	if err != nil {
		return
	}
}

func (this *Sqs) Init(conf config.Listener) {

	this.Region = conf.Region
	this.Url = conf.Url
	this.Message = conf.Message
}
