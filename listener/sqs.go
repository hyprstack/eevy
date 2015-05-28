package listener

import (
	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/sqs"

	"github.com/hevnly/eevy/event"
)

type SqsConfig interface {
	ListenerConfig

	GetUrl() string
	GetRegion() string
}

// Places a message on an AWS SQS when relavent event is triggered
type Sqs struct {
	ListenerBase

	Config SqsConfig
}

// Satisfies the Listener interface and places the event on an AWS SQS
func (this *Sqs) Exec(evt event.Event) {

	gLog.Debug("SQS %s on event %s", this.Config.GetUrl(), evt.Id)

	url := magicString(this.Config.GetUrl(), evt)
	reg := magicString(this.Config.GetRegion(), evt)
	msg := magicString(this.Config.GetMessage(), evt)
	gLog.Debug("debug: %s", url)
	svc := sqs.New(&aws.Config{Region: reg})
	params := &sqs.SendMessageInput{
		MessageBody: aws.String(msg),
		QueueURL:    aws.String(url),
	}
	_, err := svc.SendMessage(params)

	if err != nil {
		gLog.Error(err.Error())
		return
	}
}
