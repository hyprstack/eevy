package listener

import (
	"bytes"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/lambda"

	"github.com/hevnly/eevy/config"
	"github.com/hevnly/eevy/event"
)

type Lambda struct {
	ListenerBase

	Region   string
	Function string
}

func (this *Lambda) Exec(evt event.Event) {

	gLog.Debug("Lambda %s on event %s", this.Function, evt.Id)
	msg := this.GetMessage(evt)

	svc := lambda.New(&aws.Config{Region: this.Region})
	params := &lambda.InvokeAsyncInput{
		FunctionName: aws.String(this.Function),
		InvokeArgs:   bytes.NewReader([]byte(msg)),
	}
	_, err := svc.InvokeAsync(params)

	if awserr := aws.Error(err); awserr != nil {
		gLog.Error(awserr.Error())
		return
	} else if err != nil {
		gLog.Error(err.Error())
		return
	}
}

func (this *Lambda) Init(conf config.Listener) {

	this.Region = conf.Region
	this.Function = conf.Function
	this.Message = conf.Message
}
