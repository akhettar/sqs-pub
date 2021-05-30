package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSReplayer struct {
	sess *session.Session
	cfg  SQSReplayConfig
}

type SQSReplayConfig struct {
	from             string
	to               string
	deleteFromSource bool
	filters          string
}

func NewSQSReplayer(cfg SQSReplayConfig) *SQSReplayer {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return &SQSReplayer{sess: session, cfg: cfg}
}

func (s *SQSReplayer) replay(ctx context.Context, args []string) error {

	// Fetch the urls for the given queues
	fromQueue, toQueue := s.fetchQueueUrl(s.cfg.from, s.cfg.to)

	// Create report files
	preport := createReportFile("processed.log")
	failedReport := createReportFile("failed.log")
	filtered := createReportFile("filtered.log")

	defer preport.Close()
	defer failedReport.Close()
	defer filtered.Close()

	// Fetch all the messages from the source queue and publish them to the destination queue
	processed := 0
	failureToReadFromQueue := 0
	for {

		msgResult, err := s.read(fromQueue)

		if err != nil {
			log.Printf("failed to read from the queue: %v", err)
			failureToReadFromQueue++
		}

		if len(msgResult.Messages) == 0 || failureToReadFromQueue > 20 {
			log.Println("\n no more messages in the queue to replay")
			break
		}

		for _, msg := range msgResult.Messages {

			if !s.filter(msg.Body) {
				log.Printf("Processing message: %d with id %s\n", processed, *msg.MessageId)

				if err := s.send(toQueue, *msg.Body); err != nil {
					log.Printf("failed to send message: %s", *msg.Body)
					failedReport.WriteString(*msg.MessageId)
					failedReport.WriteString(*msg.Body)
					failedReport.WriteString("\n\n")

				} else {
					preport.WriteString(*msg.MessageId)
					preport.WriteString(*msg.Body)
					preport.WriteString("\n\n")
					if s.cfg.deleteFromSource {
						if err := s.delete(fromQueue, *msg.ReceiptHandle); err != nil {
							log.Printf("failed to delete message with id %s", *msg.MessageId)
						}
					}
				}
				processed++
			} else {
				log.Printf("Message with message id filtered: %s", *msg.MessageId)
				filtered.WriteString(*msg.MessageId)
				filtered.WriteString(*msg.Body)
				filtered.WriteString("\n\n")
			}

		}
	}
	log.Printf("number of messages processed %d", processed)
	log.Printf("number of messages failed to be processed %d", failureToReadFromQueue)
	return nil
}

func (s *SQSReplayer) filter(body *string) bool {
	fmt.Printf("filters %v", s.cfg.filters)
	if len(s.cfg.filters) > 0 {
		for _, text := range strings.Split(s.cfg.filters, ",") {
			if strings.Contains(*body, text) {
				return true
			}
		}
	}
	return false
}

func createReportFile(file string) *os.File {
	report, err := os.Create(file)

	if err != nil {
		log.Fatal(err)
	}
	return report
}

func (s *SQSReplayer) fetchQueueUrl(from, to string) (string, string) {
	svc := sqs.New(s.sess)
	resultFrom, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &from,
	})

	if err != nil {
		log.Fatalf("failed to fetch queue url for given queue: %s, err %v", from, err)
	}

	resultTo, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &to,
	})

	if err != nil {
		log.Fatalf("failed to fetch queue url for given queue: %s, err %v", to, err)
	}

	return *resultFrom.QueueUrl, *resultTo.QueueUrl
}
func (s *SQSReplayer) read(queue string) (*sqs.ReceiveMessageOutput, error) {
	svc := sqs.New(s.sess)
	return svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &queue,
		MaxNumberOfMessages: aws.Int64(5),
	})
}

func (s *SQSReplayer) send(queue string, body string) error {
	svc := sqs.New(s.sess)
	_, err := svc.SendMessage(&sqs.SendMessageInput{

		MessageBody: aws.String(body),
		QueueUrl:    &queue,
	})
	return err
}

func (s *SQSReplayer) delete(queueURL string, receiptHandle string) error {

	svc := sqs.New(s.sess)

	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: &receiptHandle,
	})
	return err
}
