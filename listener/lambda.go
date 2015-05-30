package listener

import (
	"bytes"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/lambda"

	"github.com/hevnly/eevy/event"
	"github.com/hevnly/eevy/listener/config"
)

// Executes an AWS Lambda function when relevant event triggered
type Lambda struct {
	ListenerBase

	Config config.Lambda
}

// Satifies the Listener interface and calls the Lambda function
func (this *Lambda) Exec(evt event.Event) {

	this.Log.Listener(this, &evt)

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
		this.Log.Error(err.Error())
		return
	}
}

func (this *Lambda) GetType() string {

	return this.GetConfig().GetType()
}

func (this *Lambda) GetConfig() config.Listener {
	return this.Config
}
