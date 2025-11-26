package sqs_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	sqspkg "github.com/takeuchi-shogo/golang-learn/app/backend/pkg/aws/sqs"
	sqsmock "github.com/takeuchi-shogo/golang-learn/app/backend/pkg/aws/sqs/mock"
	"go.uber.org/mock/gomock"
)

// TestReceiveMessages は ReceiveMessages メソッドのテストです。
// メッセージの受信が正常に行われることと、エラーハンドリングが適切に行われることを確認します。
func TestReceiveMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		queueURL         string
		options          []sqspkg.ReceiveMessageOptionFunc
		setupMock        func(mockAPI *sqsmock.MockSQSAPI)
		expectedMessages int
		expectedError    bool
	}{
		{
			name:     "正常系: メッセージを正常に受信できる",
			queueURL: "https://sqs.ap-northeast-1.amazonaws.com/123456789012/test-queue",
			options:  nil,
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					ReceiveMessage(gomock.Any(), gomock.Any()).
					Return(&sqs.ReceiveMessageOutput{
						Messages: []types.Message{
							{
								MessageId: aws.String("msg-1"),
								Body:      aws.String("test message 1"),
							},
							{
								MessageId: aws.String("msg-2"),
								Body:      aws.String("test message 2"),
							},
						},
					}, nil)
			},
			expectedMessages: 2,
			expectedError:    false,
		},
		{
			name:     "正常系: WithMaxMessages オプション付きでメッセージを受信できる",
			queueURL: "https://sqs.ap-northeast-1.amazonaws.com/123456789012/test-queue",
			options:  []sqspkg.ReceiveMessageOptionFunc{sqspkg.WithMaxMessages(10)},
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					ReceiveMessage(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
						// オプションが正しく適用されているか確認
						if params.MaxNumberOfMessages != 10 {
							t.Errorf("expected MaxNumberOfMessages to be 10, got %d", params.MaxNumberOfMessages)
						}
						return &sqs.ReceiveMessageOutput{
							Messages: []types.Message{
								{
									MessageId: aws.String("msg-1"),
									Body:      aws.String("test message"),
								},
							},
						}, nil
					})
			},
			expectedMessages: 1,
			expectedError:    false,
		},
		{
			name:     "正常系: WithWaitTimeSeconds オプション付きでメッセージを受信できる",
			queueURL: "https://sqs.ap-northeast-1.amazonaws.com/123456789012/test-queue",
			options:  []sqspkg.ReceiveMessageOptionFunc{sqspkg.WithWaitTimeSeconds(20)},
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					ReceiveMessage(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
						// オプションが正しく適用されているか確認
						if params.WaitTimeSeconds != 20 {
							t.Errorf("expected WaitTimeSeconds to be 20, got %d", params.WaitTimeSeconds)
						}
						return &sqs.ReceiveMessageOutput{
							Messages: []types.Message{
								{
									MessageId: aws.String("msg-1"),
									Body:      aws.String("test message"),
								},
							},
						}, nil
					})
			},
			expectedMessages: 1,
			expectedError:    false,
		},
		{
			name:     "正常系: 複数のオプションを組み合わせてメッセージを受信できる",
			queueURL: "https://sqs.ap-northeast-1.amazonaws.com/123456789012/test-queue",
			options:  []sqspkg.ReceiveMessageOptionFunc{sqspkg.WithMaxMessages(5), sqspkg.WithWaitTimeSeconds(10)},
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					ReceiveMessage(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
						// オプションが正しく適用されているか確認
						if params.MaxNumberOfMessages != 5 {
							t.Errorf("expected MaxNumberOfMessages to be 5, got %d", params.MaxNumberOfMessages)
						}
						if params.WaitTimeSeconds != 10 {
							t.Errorf("expected WaitTimeSeconds to be 10, got %d", params.WaitTimeSeconds)
						}
						return &sqs.ReceiveMessageOutput{
							Messages: []types.Message{
								{
									MessageId: aws.String("msg-1"),
									Body:      aws.String("test message"),
								},
							},
						}, nil
					})
			},
			expectedMessages: 1,
			expectedError:    false,
		},
		{
			name:     "正常系: メッセージが存在しない場合は空のスライスを返す",
			queueURL: "https://sqs.ap-northeast-1.amazonaws.com/123456789012/test-queue",
			options:  nil,
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					ReceiveMessage(gomock.Any(), gomock.Any()).
					Return(&sqs.ReceiveMessageOutput{
						Messages: []types.Message{},
					}, nil)
			},
			expectedMessages: 0,
			expectedError:    false,
		},
		{
			name:     "異常系: SQSクライアントがエラーを返す",
			queueURL: "https://sqs.ap-northeast-1.amazonaws.com/123456789012/test-queue",
			options:  nil,
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					ReceiveMessage(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("SQS service error"))
			},
			expectedMessages: 0,
			expectedError:    true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// gomockのコントローラ作成
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// モッククライアントの作成
			mockAPI := sqsmock.NewMockSQSAPI(ctrl)
			tt.setupMock(mockAPI)

			// sqsClientの作成（テスト用コンストラクタ使用）
			client := sqspkg.NewSQSClientWithAPI(mockAPI)

			// テスト実行
			ctx := context.Background()
			messages, err := client.ReceiveMessages(ctx, tt.queueURL, tt.options...)

			// エラーチェック
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}

			// メッセージ数の検証
			if !tt.expectedError && len(messages) != tt.expectedMessages {
				t.Errorf("expected %d messages, got %d", tt.expectedMessages, len(messages))
			}
		})
	}
}

// TestDeleteMessage は DeleteMessage メソッドのテストです。
// メッセージの削除が正常に行われることと、エラーハンドリングが適切に行われることを確認します。
func TestDeleteMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		queueURL      string
		options       []sqspkg.DeleteMessageOptionFunc
		setupMock     func(mockAPI *sqsmock.MockSQSAPI)
		expectedError bool
	}{
		{
			name:     "正常系: メッセージを正常に削除できる",
			queueURL: "https://sqs.ap-northeast-1.amazonaws.com/123456789012/test-queue",
			options:  []sqspkg.DeleteMessageOptionFunc{sqspkg.WithReceiptHandle("receipt-handle-123")},
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					DeleteMessage(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
						// オプションが正しく適用されているか確認
						if params.ReceiptHandle == nil || *params.ReceiptHandle != "receipt-handle-123" {
							t.Errorf("expected ReceiptHandle to be 'receipt-handle-123', got %v", params.ReceiptHandle)
						}
						return &sqs.DeleteMessageOutput{}, nil
					})
			},
			expectedError: false,
		},
		{
			name:     "正常系: オプションなしで削除できる",
			queueURL: "https://sqs.ap-northeast-1.amazonaws.com/123456789012/test-queue",
			options:  nil,
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					DeleteMessage(gomock.Any(), gomock.Any()).
					Return(&sqs.DeleteMessageOutput{}, nil)
			},
			expectedError: false,
		},
		{
			name:     "異常系: SQSクライアントがエラーを返す",
			queueURL: "https://sqs.ap-northeast-1.amazonaws.com/123456789012/test-queue",
			options:  []sqspkg.DeleteMessageOptionFunc{sqspkg.WithReceiptHandle("receipt-handle-123")},
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					DeleteMessage(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("SQS service error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// gomockのコントローラ作成
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// モッククライアントの作成
			mockAPI := sqsmock.NewMockSQSAPI(ctrl)
			tt.setupMock(mockAPI)

			// sqsClientの作成
			client := sqspkg.NewSQSClientWithAPI(mockAPI)

			// テスト実行
			ctx := context.Background()
			err := client.DeleteMessage(ctx, tt.queueURL, tt.options...)

			// エラーチェック
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}

// TestSendMessage は SendMessage メソッドのテストです。
// メッセージの送信が正常に行われることと、エラーハンドリングが適切に行われることを確認します。
func TestSendMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		queueURL          string
		options           []sqspkg.SendMessageOptionFunc
		setupMock         func(mockAPI *sqsmock.MockSQSAPI)
		expectedMessageID string
		expectedError     bool
	}{
		{
			name:     "正常系: メッセージを正常に送信できる",
			queueURL: "https://sqs.ap-northeast-1.amazonaws.com/123456789012/test-queue",
			options:  []sqspkg.SendMessageOptionFunc{sqspkg.WithMessageBody("test message body")},
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
						// オプションが正しく適用されているか確認
						if params.MessageBody == nil || *params.MessageBody != "test message body" {
							t.Errorf("expected MessageBody to be 'test message body', got %v", params.MessageBody)
						}
						return &sqs.SendMessageOutput{
							MessageId: aws.String("msg-id-123"),
						}, nil
					})
			},
			expectedMessageID: "msg-id-123",
			expectedError:     false,
		},
		{
			name:     "正常系: オプションなしで送信できる",
			queueURL: "https://sqs.ap-northeast-1.amazonaws.com/123456789012/test-queue",
			options:  nil,
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					Return(&sqs.SendMessageOutput{
						MessageId: aws.String("msg-id-456"),
					}, nil)
			},
			expectedMessageID: "msg-id-456",
			expectedError:     false,
		},
		{
			name:     "正常系: JSON形式のメッセージを送信できる",
			queueURL: "https://sqs.ap-northeast-1.amazonaws.com/123456789012/test-queue",
			options:  []sqspkg.SendMessageOptionFunc{sqspkg.WithMessageBody(`{"task":"test","data":"hello"}`)},
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
						if params.MessageBody == nil || *params.MessageBody != `{"task":"test","data":"hello"}` {
							t.Errorf("expected JSON MessageBody, got %v", params.MessageBody)
						}
						return &sqs.SendMessageOutput{
							MessageId: aws.String("msg-id-789"),
						}, nil
					})
			},
			expectedMessageID: "msg-id-789",
			expectedError:     false,
		},
		{
			name:     "異常系: SQSクライアントがエラーを返す",
			queueURL: "https://sqs.ap-northeast-1.amazonaws.com/123456789012/test-queue",
			options:  []sqspkg.SendMessageOptionFunc{sqspkg.WithMessageBody("test message body")},
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					SendMessage(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("SQS service error"))
			},
			expectedMessageID: "",
			expectedError:     true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// gomockのコントローラ作成
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// モッククライアントの作成
			mockAPI := sqsmock.NewMockSQSAPI(ctrl)
			tt.setupMock(mockAPI)

			// sqsClientの作成
			client := sqspkg.NewSQSClientWithAPI(mockAPI)

			// テスト実行
			ctx := context.Background()
			output, err := client.SendMessage(ctx, tt.queueURL, tt.options...)

			// エラーチェック
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}

			// 出力の検証
			if !tt.expectedError {
				if output == nil {
					t.Error("expected output to be non-nil")
				} else if output.MessageId != nil && *output.MessageId != tt.expectedMessageID {
					t.Errorf("expected MessageId %s, got %s", tt.expectedMessageID, *output.MessageId)
				}
			}
		})
	}
}

// TestGetQueueURL は GetQueueURL メソッドのテストです。
// キューURLの取得が正常に行われることと、エラーハンドリングが適切に行われることを確認します。
func TestGetQueueURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		queueName     string
		setupMock     func(mockAPI *sqsmock.MockSQSAPI)
		expectedURL   string
		expectedError bool
	}{
		{
			name:      "正常系: キューURLを正常に取得できる",
			queueName: "test-queue",
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					GetQueueUrl(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error) {
						// パラメータが正しく設定されているか確認
						if params.QueueName == nil || *params.QueueName != "test-queue" {
							t.Errorf("expected QueueName to be 'test-queue', got %v", params.QueueName)
						}
						return &sqs.GetQueueUrlOutput{
							QueueUrl: aws.String("https://sqs.ap-northeast-1.amazonaws.com/123456789012/test-queue"),
						}, nil
					})
			},
			expectedURL:   "https://sqs.ap-northeast-1.amazonaws.com/123456789012/test-queue",
			expectedError: false,
		},
		{
			name:      "正常系: 別のキュー名でURLを取得できる",
			queueName: "production-queue",
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					GetQueueUrl(gomock.Any(), gomock.Any()).
					Return(&sqs.GetQueueUrlOutput{
						QueueUrl: aws.String("https://sqs.ap-northeast-1.amazonaws.com/123456789012/production-queue"),
					}, nil)
			},
			expectedURL:   "https://sqs.ap-northeast-1.amazonaws.com/123456789012/production-queue",
			expectedError: false,
		},
		{
			name:      "異常系: SQSクライアントがエラーを返す",
			queueName: "test-queue",
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					GetQueueUrl(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("SQS service error"))
			},
			expectedURL:   "",
			expectedError: true,
		},
		{
			name:      "異常系: キューが存在しない場合にエラーを返す",
			queueName: "non-existent-queue",
			setupMock: func(mockAPI *sqsmock.MockSQSAPI) {
				mockAPI.EXPECT().
					GetQueueUrl(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("queue does not exist"))
			},
			expectedURL:   "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// gomockのコントローラ作成
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// モッククライアントの作成
			mockAPI := sqsmock.NewMockSQSAPI(ctrl)
			tt.setupMock(mockAPI)

			// sqsClientの作成
			client := sqspkg.NewSQSClientWithAPI(mockAPI)

			// テスト実行
			ctx := context.Background()
			url, err := client.GetQueueURL(ctx, tt.queueName)

			// エラーチェック
			if (err != nil) != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}

			// URLの検証
			if !tt.expectedError && url != tt.expectedURL {
				t.Errorf("expected URL %s, got %s", tt.expectedURL, url)
			}
		})
	}
}

// TestWithMaxMessages は WithMaxMessages オプション関数のテストです。
func TestWithMaxMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		maxMessages int32
		expected    int32
	}{
		{
			name:        "正常系: MaxMessagesに1を設定できる",
			maxMessages: 1,
			expected:    1,
		},
		{
			name:        "正常系: MaxMessagesに10を設定できる",
			maxMessages: 10,
			expected:    10,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			input := &sqs.ReceiveMessageInput{}
			option := sqspkg.WithMaxMessages(tt.maxMessages)
			option(input)

			if input.MaxNumberOfMessages != tt.expected {
				t.Errorf("expected MaxNumberOfMessages to be %d, got %d", tt.expected, input.MaxNumberOfMessages)
			}
		})
	}
}

// TestWithWaitTimeSeconds は WithWaitTimeSeconds オプション関数のテストです。
func TestWithWaitTimeSeconds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		waitTimeSeconds int32
		expected        int32
	}{
		{
			name:            "正常系: WaitTimeSecondsに0を設定できる",
			waitTimeSeconds: 0,
			expected:        0,
		},
		{
			name:            "正常系: WaitTimeSecondsに20を設定できる",
			waitTimeSeconds: 20,
			expected:        20,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			input := &sqs.ReceiveMessageInput{}
			option := sqspkg.WithWaitTimeSeconds(tt.waitTimeSeconds)
			option(input)

			if input.WaitTimeSeconds != tt.expected {
				t.Errorf("expected WaitTimeSeconds to be %d, got %d", tt.expected, input.WaitTimeSeconds)
			}
		})
	}
}

// TestWithReceiptHandle は WithReceiptHandle オプション関数のテストです。
func TestWithReceiptHandle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		receiptHandle string
		expected      string
	}{
		{
			name:          "正常系: ReceiptHandleに通常の値を設定できる",
			receiptHandle: "receipt-handle-123",
			expected:      "receipt-handle-123",
		},
		{
			name:          "正常系: ReceiptHandleに長い値を設定できる",
			receiptHandle: "very-long-receipt-handle-with-many-characters-1234567890",
			expected:      "very-long-receipt-handle-with-many-characters-1234567890",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			input := &sqs.DeleteMessageInput{}
			option := sqspkg.WithReceiptHandle(tt.receiptHandle)
			option(input)

			if input.ReceiptHandle == nil {
				t.Error("expected ReceiptHandle to be non-nil")
			} else if *input.ReceiptHandle != tt.expected {
				t.Errorf("expected ReceiptHandle to be %s, got %s", tt.expected, *input.ReceiptHandle)
			}
		})
	}
}

// TestWithMessageBody は WithMessageBody オプション関数のテストです。
func TestWithMessageBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		messageBody string
		expected    string
	}{
		{
			name:        "正常系: MessageBodyに通常のメッセージを設定できる",
			messageBody: "test message body",
			expected:    "test message body",
		},
		{
			name:        "正常系: MessageBodyにJSON形式のメッセージを設定できる",
			messageBody: `{"key": "value", "number": 123}`,
			expected:    `{"key": "value", "number": 123}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			input := &sqs.SendMessageInput{}
			option := sqspkg.WithMessageBody(tt.messageBody)
			option(input)

			if input.MessageBody == nil {
				t.Error("expected MessageBody to be non-nil")
			} else if *input.MessageBody != tt.expected {
				t.Errorf("expected MessageBody to be %s, got %s", tt.expected, *input.MessageBody)
			}
		})
	}
}

// TestNewSQSClientWithAPI は NewSQSClientWithAPI コンストラクタのテストです。
func TestNewSQSClientWithAPI(t *testing.T) {
	t.Parallel()

	t.Run("正常系: NewSQSClientWithAPIでクライアントを作成できる", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAPI := sqsmock.NewMockSQSAPI(ctrl)
		client := sqspkg.NewSQSClientWithAPI(mockAPI)

		if client == nil {
			t.Error("expected client to be non-nil")
		}
	})
}
