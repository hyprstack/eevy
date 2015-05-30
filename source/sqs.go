package source

import (
	"sync"
	"time"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/sqs"
)

// Listens to an AWS SQS.  When a message is placed onto this queue it is
// converted into a Event and assesed by the listeners.  The message must be
// a JSON string representing an Event
type Sqs struct {
	Base

	svc *sqs.SQS
}

// Satisfies the Source interface.  Beigins listening to an AWS SQS queue.  If no
// message is on the queue it sleeps for a set period of time before trying again
func (s *Sqs) Listen(wg sync.WaitGroup) {

	s.Log.Info("Start listening (sqs:%s)", s.Url)
	for {
		numMsg := s.recieve()
		if numMsg == 0 { // if no message sleep for a period of time
			time.Sleep(5 * time.Second)
		}
	}
	wg.Done()
}

// Check to see if we have a message on the queue
func (s *Sqs) recieve() int {

	if s.svc == nil {
		s.svc = sqs.New(&aws.Config{Region: s.Region})
	}

	params := &sqs.ReceiveMessageInput{
		QueueURL:            aws.String(s.Url),
		MaxNumberOfMessages: aws.Long(10),
	}
	resp, err := s.svc.ReceiveMessage(params)

	if err != nil {
		// A non-service error occurred.
		s.Log.Error(err.Error())
		return 0
	}

	numMsg := len(resp.Messages)
	if numMsg == 0 {
		return 0
	}

	for _, element := range resp.Messages {
		s.processRaw(*element.Body)
		s.remove(element)
	}

	return numMsg
}

// Remove the processed message from the SQS
func (s *Sqs) remove(message *sqs.Message) {

	params := &sqs.DeleteMessageInput{
		QueueURL:      aws.String(s.Url),
		ReceiptHandle: message.ReceiptHandle,
	}
	_, err := s.svc.DeleteMessage(params)

	if err != nil {
		s.Log.Error(err.Error())
	}
}
