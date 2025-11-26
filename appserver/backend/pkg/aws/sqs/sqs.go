package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

//go:generate mockgen -source=$GOFILE -destination=${GOPACKAGE}mock/$GOFILE -package=${GOPACKAGE}mock

// SQS はSQSクライアントの外部公開用インターフェースです。
// メッセージの送受信、削除、キューURL取得の機能を提供します。
type SQS interface {
	// ReceiveMessages はSQSからメッセージを受信するメソッドです。
	// オプションで最大メッセージ数や待ち時間を指定できます。
	ReceiveMessages(ctx context.Context, queueURL string, options ...ReceiveMessageOptionFunc) ([]types.Message, error)

	// DeleteMessage はSQSのメッセージを削除するメソッドです。
	DeleteMessage(ctx context.Context, queueURL string, options ...DeleteMessageOptionFunc) error

	// SendMessage はSQSにメッセージを送信するメソッドです。
	// オプションでメッセージの内容や属性を指定できます。
	SendMessage(ctx context.Context, queueURL string, options ...SendMessageOptionFunc) (*sqs.SendMessageOutput, error)

	// GetQueueURL はSQSのキューURLを取得するメソッドです。
	// キュー名を指定すると、キューURLを取得できます。
	GetQueueURL(ctx context.Context, queueName string) (string, error)
}

// SQSAPI はAWS SDK v2のSQSクライアントが実装すべき内部用インターフェースです。
// テスト時にモッククライアントを注入するために使用します。
type SQSAPI interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
	GetQueueUrl(ctx context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)
}

// sqsClient はSQSインターフェースの実装です。
type sqsClient struct {
	client SQSAPI
}

// NewSQSClient はSQSクライアントを生成するコンストラクタです。
// 引数:
//   - client: AWS SDK v2のSQSクライアント (*sqs.Client)
//
// 戻り値:
//   - SQS: SQSインターフェースの実装
func NewSQSClient(client *sqs.Client) SQS {
	return &sqsClient{client: client}
}

// NewSQSClientWithAPI はテスト用のコンストラクタです。
// SQSAPIインターフェースを受け取り、モッククライアントの注入を可能にします。
func NewSQSClientWithAPI(client SQSAPI) SQS {
	return &sqsClient{client: client}
}

// ReceiveMessageOptionFunc はReceiveMessagesのオプション関数型です。
type ReceiveMessageOptionFunc func(*sqs.ReceiveMessageInput)

// WithMaxMessages は一度に受信するメッセージの最大数を設定します。
// 引数:
//   - maxMessages: 最大メッセージ数 (1-10)
func WithMaxMessages(maxMessages int32) ReceiveMessageOptionFunc {
	return func(input *sqs.ReceiveMessageInput) {
		input.MaxNumberOfMessages = maxMessages
	}
}

// WithWaitTimeSeconds はロングポーリングの待機時間を設定します。
// 引数:
//   - waitTimeSeconds: 待機時間（秒）(0-20)
func WithWaitTimeSeconds(waitTimeSeconds int32) ReceiveMessageOptionFunc {
	return func(input *sqs.ReceiveMessageInput) {
		input.WaitTimeSeconds = waitTimeSeconds
	}
}

func (s *sqsClient) ReceiveMessages(ctx context.Context, queueURL string, options ...ReceiveMessageOptionFunc) ([]types.Message, error) {
	input := &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueURL),
	}
	for _, option := range options {
		option(input)
	}

	output, err := s.client.ReceiveMessage(ctx, input)
	if err != nil {
		return nil, err
	}
	return output.Messages, nil
}

// DeleteMessageOptionFunc はDeleteMessageのオプション関数型です。
type DeleteMessageOptionFunc func(*sqs.DeleteMessageInput)

// WithReceiptHandle は削除するメッセージのレシートハンドルを設定します。
// 引数:
//   - receiptHandle: メッセージ受信時に取得したレシートハンドル
func WithReceiptHandle(receiptHandle string) DeleteMessageOptionFunc {
	return func(input *sqs.DeleteMessageInput) {
		input.ReceiptHandle = aws.String(receiptHandle)
	}
}

func (s *sqsClient) DeleteMessage(ctx context.Context, queueURL string, options ...DeleteMessageOptionFunc) error {
	input := &sqs.DeleteMessageInput{
		QueueUrl: aws.String(queueURL),
	}
	for _, option := range options {
		option(input)
	}
	_, err := s.client.DeleteMessage(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

// SendMessageOptionFunc はSendMessageのオプション関数型です。
type SendMessageOptionFunc func(*sqs.SendMessageInput)

// WithMessageBody は送信するメッセージの本文を設定します。
// 引数:
//   - messageBody: メッセージ本文
func WithMessageBody(messageBody string) SendMessageOptionFunc {
	return func(input *sqs.SendMessageInput) {
		input.MessageBody = aws.String(messageBody)
	}
}

func (s *sqsClient) SendMessage(ctx context.Context, queueURL string, options ...SendMessageOptionFunc) (*sqs.SendMessageOutput, error) {
	input := &sqs.SendMessageInput{
		QueueUrl: aws.String(queueURL),
	}
	for _, option := range options {
		option(input)
	}
	output, err := s.client.SendMessage(ctx, input)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (s *sqsClient) GetQueueURL(ctx context.Context, queueName string) (string, error) {
	input := &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	}
	output, err := s.client.GetQueueUrl(ctx, input)
	if err != nil {
		return "", err
	}
	return *output.QueueUrl, nil
}
