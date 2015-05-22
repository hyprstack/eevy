package source

import (
	"sync"
	"time"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/sqs"
)

type Sqs struct {
	Base

	svc *sqs.SQS
}

func (s *Sqs) Listen(wg sync.WaitGroup) {

	s.AppLog.Info("Start listening (sqs:%s)", s.Url)
	for {
		numMsg := s.recieve()
		if numMsg == 0 {
			time.Sleep(5 * time.Second)
		}
	}
	wg.Done()
}

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
		s.AppLog.Error(err.Error())
		return 0
	}

	numMsg := len(resp.Messages)
	if numMsg == 0 {
		return 0
	}

	s.AppLog.Info("Recieved %d messages", numMsg)
	for _, element := range resp.Messages {
		s.processRaw(*element.Body)
		s.remove(element)
	}

	return numMsg
}

func (s *Sqs) remove(message *sqs.Message) {

	s.AppLog.Info("Deleteing: '%s'", *message.Body)

	params := &sqs.DeleteMessageInput{
		QueueURL:      aws.String(s.Url),
		ReceiptHandle: message.ReceiptHandle,
	}
	_, err := s.svc.DeleteMessage(params)

	if err != nil {
		s.AppLog.Error(err.Error())
	}
}
