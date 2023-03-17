package queue

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SendMessageClient interface {
	SendMessage(context.Context, *sqs.SendMessageInput, ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

type SQSError struct {
	Message string
	Faults  []string
}

type SendMessageOutput struct {
	MessageID      *string
	SequenceNumber *string
}

func (se SQSError) Error() string {
	sb := strings.Builder{}

	_, _ = sb.WriteString(se.Message)
	_, _ = sb.WriteString(": faults [ ")
	_, _ = sb.WriteString(strings.Join(se.Faults, ", "))
	_, _ = sb.WriteString(" ]")
	return sb.String()
}

type EnqueueInput struct {
	QueueURL       string
	Body           *string
	Attributes     map[string]string
	MessageGroupID *string
}

func (ei *EnqueueInput) Validate() error {
	errors := make([]string, 0)
	if len(ei.QueueURL) == 0 {
		errors = append(errors, "missing queue url")
	}

	if ei.Body == nil || len(*ei.Body) == 0 {
		errors = append(errors, "sqs message body do not have a value")
	}

	if strings.HasSuffix(ei.QueueURL, ".fifo") {
		if ei.MessageGroupID == nil || len(*ei.MessageGroupID) == 0 {
			errors = append(errors, "message group id must be set for fifo queues")
		}
	} else {
		if ei.MessageGroupID != nil && len(*ei.MessageGroupID) > 0 {
			errors = append(errors, "message group id can only be used with fifo queues")
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return &SQSError{
		Message: "validation of sqs message failed",
		Faults:  errors,
	}
}

func (ei *EnqueueInput) SetMessageGroupID(id string) {
	ei.MessageGroupID = new(string)
	*ei.MessageGroupID = id
}

func (ei *EnqueueInput) GetMessageGroupID() *string {
	if ei.MessageGroupID == nil || len(*ei.MessageGroupID) == 0 {
		return nil
	}
	return ei.MessageGroupID
}

func (ei *EnqueueInput) SetBody(input string) {
	ei.Body = new(string)
	*ei.Body = input
}

func (ei *EnqueueInput) SetBodyBytes(input []byte) {
	ei.Body = new(string)
	*ei.Body = string(input)
}

func (ei *EnqueueInput) GetBody() *string {
	if ei.Body != nil && len(*ei.Body) == 0 {
		return nil
	}

	return ei.Body
}

func (ei *EnqueueInput) GetURL() *string {
	return &ei.QueueURL
}

func NewQueueInput(queueURL, msg string) *EnqueueInput {
	input := &EnqueueInput{
		QueueURL: queueURL,
		Body:     &msg,
	}

	return input
}

func (ei *EnqueueInput) GetMessageAttributes() map[string]types.MessageAttributeValue {
	if len(ei.Attributes) == 0 {
		return nil
	}

	attrs := make(map[string]types.MessageAttributeValue)

	for k, v := range ei.Attributes {
		attrs[k] = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(v),
		}
	}

	return attrs
}

func (ei *EnqueueInput) SetMessageAttributes(key, val string) {
	if len(key) == 0 && len(val) == 0 {
		return
	}

	if ei.Attributes == nil {
		ei.Attributes = make(map[string]string)
	}

	ei.Attributes[key] = val
}

func Enqueue(ctx context.Context, client SendMessageClient, input *EnqueueInput) (*SendMessageOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	msg := sqs.SendMessageInput{
		QueueUrl:          input.GetURL(),
		MessageBody:       input.GetBody(),
		MessageGroupId:    input.GetMessageGroupID(),
		MessageAttributes: input.GetMessageAttributes(),
	}

	output, err := client.SendMessage(ctx, &msg)
	if err != nil {
		return nil, err
	}

	return &SendMessageOutput{
		MessageID:      output.MessageId,
		SequenceNumber: output.SequenceNumber,
	}, nil
}
