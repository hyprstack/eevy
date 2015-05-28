package listener

import (
	"bytes"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/lambda"

	"github.com/hevnly/eevy/event"
)

type LambdaConfig interface {
	ListenerConfig

	GetFunction() string
	GetRegion() string
}

// Executes an AWS Lambda function when relevant event triggered
type Lambda struct {
	ListenerBase

	Config LambdaConfig
}

// Satifies the Listener interface and calls the Lambda function
func (this *Lambda) Exec(evt event.Event) {

	gLog.Debug("Lambda %s on event %s", this.Config.GetFunction(), evt.Id)

	msg := magicString(this.Config.GetMessage(), evt)
	reg := magicString(this.Config.GetRegion(), evt)
	fun := magicString(this.Config.GetFunction(), evt)

	svc := lambda.New(&aws.Config{Region: reg})
	params := &lambda.InvokeAsyncInput{
		FunctionName: aws.String(fun),
		InvokeArgs:   bytes.NewReader([]byte(msg)),
	}
	_, err := svc.InvokeAsync(params)

	if err != nil {
		gLog.Error(err.Error())
		return
	}
}
