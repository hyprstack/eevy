package listener

import (
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/sqs"

	"github.com/hevnly/eevy/config"
	"github.com/hevnly/eevy/event"
)

// Places a message on an AWS SQS when relavent event is triggered
type Sqs struct {
	ListenerBase

	Url    string
	Region string
}

// Satisfies the Listener interface and places the event on an AWS SQS
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
