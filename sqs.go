package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSMessageReplayer type
type SQSMessageReplayer struct {
	sess *session.Session
	cfg  SQSMessageReplayConfig
}

// SQSMessageReplayConfig type
type SQSMessageReplayConfig struct {
	from             string
	to               string
	deleteFromSource bool
	dryrun           bool
	filters          string
}

// NewSQSMessageReplayer creates an instace of the SQS Replayer
func NewSQSMessageReplayer() *SQSMessageReplayer {

	// Create Session
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return &SQSMessageReplayer{sess: session, cfg: SQSMessageReplayConfig{}}
}

func (s *SQSMessageReplayer) replay(ctx context.Context, args []string) error {

	// Fetch the urls for the given queues
	fromQueue, toQueue := s.fetchQueueURL(s.cfg.from, s.cfg.to)

	processedBody := []string{}
	failedBody := []string{}
	filteredBody := []string{}

	messages := map[string]*[]string{"processed": &processedBody, "failed": &failedBody, "filtered": &filteredBody}

	// Fetch all the messages from the source queue and publish them to the destination queue
	numOfMessags := s.fetchNumberOfMessages(fromQueue)

	for i := 0; i < numOfMessags; i++ {

		msgResult, err := s.read(fromQueue)

		if err != nil {
			log.Printf("failed to read from the queue: %v", err)
		}

		for _, msg := range msgResult.Messages {

			if !s.filter(msg.Body) {
				log.Printf("Processing message: %d with id %s\n", len(processedBody), *msg.MessageId)

				if err := s.send(toQueue, *msg.Body); err != nil {

					log.Printf("failed to send message: %s", *msg.Body)
					failedBody = append(failedBody, *msg.Body)
				} else {
					processedBody = append(processedBody, *msg.Body)
					if s.cfg.deleteFromSource {
						if err := s.delete(fromQueue, *msg.ReceiptHandle); err != nil {
							log.Printf("failed to delete message with id %s", *msg.MessageId)
						}
					}
				}
			} else {
				log.Printf("Message with message id filtered: %s", *msg.MessageId)
				filteredBody = append(filteredBody, *msg.Body)
			}

		}
	}

	log.Printf("number of messages processed %d", len(processedBody))
	log.Printf("number of messages filtered %d", len(filteredBody))
	log.Printf("number of messages failed %d", len(failedBody))

	// generating the reports
	generateReport(messages)
	return nil
}

func generateReport(messages map[string]*[]string) {
	for n, msg := range messages {
		if len(*msg) > 0 {
			pfile := createReportFile(fmt.Sprintf("%s.log", n))
			defer pfile.Close()
			for _, body := range *msg {
				pfile.WriteString(body)
				pfile.WriteString("\n")
				pfile.WriteString("-----------------------------------------------------\n")
				pfile.WriteString("\n")
			}
			log.Printf("%s.log report generated", n)
		}
	}
}

func (s *SQSMessageReplayer) fetchNumberOfMessages(queue string) int {
	svc := sqs.New(s.sess)
	numOfMessags := "ApproximateNumberOfMessages"
	result, err := svc.GetQueueAttributes(&sqs.GetQueueAttributesInput{QueueUrl: &queue, AttributeNames: []*string{&numOfMessags}})

	if err != nil {
		log.Fatal(err)
	}

	num, err := strconv.Atoi(*result.Attributes[numOfMessags])

	if err != nil {
		log.Fatalf("failed to retrieve the number of messages present in the source queue, %v", err)
	}
	log.Printf("found %d messages in the source queue %s", num, queue)
	return num
}

func (s *SQSMessageReplayer) filter(body *string) bool {
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

func (s *SQSMessageReplayer) fetchQueueURL(from, to string) (string, string) {
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
func (s *SQSMessageReplayer) read(queue string) (*sqs.ReceiveMessageOutput, error) {
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

func (s *SQSMessageReplayer) send(queue string, body string) error {
	if !s.cfg.dryrun {
		svc := sqs.New(s.sess)
		_, err := svc.SendMessage(&sqs.SendMessageInput{

			MessageBody: aws.String(body),
			QueueUrl:    &queue,
		})
		return err
	}
	return nil
}

func (s *SQSMessageReplayer) delete(queueURL string, receiptHandle string) error {
	if !s.cfg.dryrun {
		svc := sqs.New(s.sess)

		_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      &queueURL,
			ReceiptHandle: &receiptHandle,
		})
		return err
	}
	return nil
}
