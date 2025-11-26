#!/bin/bash
# =============================================================================
# LocalStack SQS 初期化スクリプト
# =============================================================================
# 目的:
#   LocalStack起動時にSQSキューを自動作成する
#
# 作成されるキュー:
#   - worker-queue: メインのワーカーキュー
#   - worker-queue-dlq: Dead Letter Queue (処理失敗メッセージ用)
#
# 使用方法:
#   このスクリプトはLocalStackの初期化ディレクトリに配置され、
#   コンテナ起動時に自動実行される
#
# 注意事項:
#   - awslocalコマンドはLocalStack内で使用可能なAWS CLI
#   - 本番環境では使用しないこと
# =============================================================================

set -e

echo "=========================================="
echo "LocalStack SQS 初期化を開始します..."
echo "=========================================="

# 設定
REGION="ap-northeast-1"
QUEUE_NAME="worker-queue"
DLQ_NAME="worker-queue-dlq"

# Dead Letter Queue (DLQ) を作成
# 処理に3回失敗したメッセージはここに移動される
echo "[1/3] Dead Letter Queue を作成中: $DLQ_NAME"
awslocal sqs create-queue \
    --queue-name "$DLQ_NAME" \
    --region "$REGION" \
    --attributes '{
        "MessageRetentionPeriod": "1209600"
    }'

# DLQのARNを取得
DLQ_URL=$(awslocal sqs get-queue-url --queue-name "$DLQ_NAME" --region "$REGION" --query 'QueueUrl' --output text)
DLQ_ARN="arn:aws:sqs:$REGION:000000000000:$DLQ_NAME"

echo "  DLQ URL: $DLQ_URL"
echo "  DLQ ARN: $DLQ_ARN"

# メインのワーカーキューを作成
# - VisibilityTimeout: 30秒 (メッセージ処理中は他のワーカーから見えない)
# - MessageRetentionPeriod: 345600秒 (4日間保持)
# - ReceiveMessageWaitTimeSeconds: 20秒 (ロングポーリング)
# - RedrivePolicy: 3回失敗したらDLQへ
echo "[2/3] メインキューを作成中: $QUEUE_NAME"
awslocal sqs create-queue \
    --queue-name "$QUEUE_NAME" \
    --region "$REGION" \
    --attributes "{
        \"VisibilityTimeout\": \"30\",
        \"MessageRetentionPeriod\": \"345600\",
        \"ReceiveMessageWaitTimeSeconds\": \"20\",
        \"RedrivePolicy\": \"{\\\"deadLetterTargetArn\\\":\\\"$DLQ_ARN\\\",\\\"maxReceiveCount\\\":\\\"3\\\"}\"
    }"

QUEUE_URL=$(awslocal sqs get-queue-url --queue-name "$QUEUE_NAME" --region "$REGION" --query 'QueueUrl' --output text)
echo "  Queue URL: $QUEUE_URL"

# 作成されたキューを一覧表示
echo "[3/3] 作成されたキュー一覧:"
awslocal sqs list-queues --region "$REGION"

echo ""
echo "=========================================="
echo "LocalStack SQS 初期化が完了しました!"
echo "=========================================="
echo ""
echo "使用方法:"
echo "  メッセージ送信:"
echo "    aws --endpoint-url=http://localhost:4566 sqs send-message \\"
echo "      --queue-url $QUEUE_URL \\"
echo "      --message-body '{\"task\":\"test\"}' \\"
echo "      --region $REGION"
echo ""
echo "  メッセージ受信:"
echo "    aws --endpoint-url=http://localhost:4566 sqs receive-message \\"
echo "      --queue-url $QUEUE_URL \\"
echo "      --region $REGION"
echo ""
