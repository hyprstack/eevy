package handler

import (
	"fmt"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awserr"
	"github.com/awslabs/aws-sdk-go/aws/awsutil"
	"github.com/awslabs/aws-sdk-go/service/lambda"

	"github.com/hevnly/eevy/event"
	"github.com/hevnly/eevy/handler/config"
)

// Executes an AWS Lambda function when relevant event triggered
type Lambda struct {
	HandlerBase

	Config config.Lambda
}

// Satisfies the Listener interface and calls the Lambda function
func (this *Lambda) Exec(evt event.Event) {

	this.Log.Handler(this, &evt)

	msg := magicString(this.Config.GetMessage(), evt)
	reg := magicString(this.Config.GetRegion(), evt)
	fun := magicString(this.Config.GetFunction(), evt)

	svc := lambda.New(&aws.Config{Region: reg})
	params := &lambda.InvokeInput{
		FunctionName: aws.String(fun), // Required
		Payload:      []byte(msg),
	}
	resp, err := svc.Invoke(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {

			msg := fmt.Sprintf("%s %s", awsErr.Code(), awsErr.Message())
			this.Log.HandlerError(this, msg, &evt)
			if _, ok := err.(awserr.RequestFailure); ok {
			}
		} else {
			this.Log.HandlerError(this, err.Error(), &evt)
		}
		return
	}
	this.Log.Debug("LAMBDA OUT %s %s \"%s\"", evt.GetId(), this.GetName(), awsutil.StringValue(resp))
}

func (this *Lambda) GetType() string {

	return this.GetConfig().GetType()
}

func (this *Lambda) GetConfig() config.Handler {
	return this.Config
}
