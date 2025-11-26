#!/bin/bash
# =============================================================================
# SQS テストメッセージ送信スクリプト
# =============================================================================
# 目的:
#   LocalStackのSQSキューにテストメッセージを送信する
#
# 使用方法:
#   ./test-sqs-send.sh                     # デフォルトメッセージを送信
#   ./test-sqs-send.sh "カスタムメッセージ"  # カスタムメッセージを送信
#   ./test-sqs-send.sh --batch 5           # 5件のメッセージをバッチ送信
#
# 注意事項:
#   - LocalStackが起動していること
#   - worker-queueが作成されていること
# =============================================================================

set -e

# 設定
ENDPOINT_URL="http://localhost:4566"
QUEUE_URL="http://localhost:4566/000000000000/worker-queue"
REGION="ap-northeast-1"

# カラー出力
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ヘルプ表示
show_help() {
    echo "使用方法: $0 [オプション] [メッセージ]"
    echo ""
    echo "オプション:"
    echo "  --batch N    N件のテストメッセージを送信"
    echo "  --help       このヘルプを表示"
    echo ""
    echo "例:"
    echo "  $0                        # デフォルトメッセージを送信"
    echo "  $0 \"Hello World\"          # カスタムメッセージを送信"
    echo "  $0 --batch 10             # 10件のメッセージを送信"
}

# シングルメッセージ送信
send_single_message() {
    local message_body=$1
    local timestamp=$(date -Iseconds)

    # デフォルトメッセージ
    if [ -z "$message_body" ]; then
        message_body="{\"task\":\"test\",\"message\":\"Hello from test script\",\"timestamp\":\"$timestamp\"}"
    fi

    echo -e "${BLUE}送信中...${NC}"
    echo "Queue URL: $QUEUE_URL"
    echo "Message: $message_body"
    echo ""

    result=$(aws --endpoint-url="$ENDPOINT_URL" sqs send-message \
        --queue-url "$QUEUE_URL" \
        --message-body "$message_body" \
        --region "$REGION" \
        --output json)

    message_id=$(echo "$result" | jq -r '.MessageId')

    echo -e "${GREEN}送信成功!${NC}"
    echo "MessageId: $message_id"
}

# バッチメッセージ送信
send_batch_messages() {
    local count=$1
    local timestamp=$(date -Iseconds)

    echo -e "${BLUE}${count}件のメッセージをバッチ送信中...${NC}"
    echo "Queue URL: $QUEUE_URL"
    echo ""

    for i in $(seq 1 $count); do
        message_body="{\"task\":\"batch_test\",\"id\":$i,\"timestamp\":\"$timestamp\"}"

        result=$(aws --endpoint-url="$ENDPOINT_URL" sqs send-message \
            --queue-url "$QUEUE_URL" \
            --message-body "$message_body" \
            --region "$REGION" \
            --output json)

        message_id=$(echo "$result" | jq -r '.MessageId')
        echo "  [$i/$count] MessageId: $message_id"
    done

    echo ""
    echo -e "${GREEN}${count}件のメッセージを送信完了!${NC}"
}

# キュー状態確認
show_queue_status() {
    echo ""
    echo "=== キュー状態 ==="
    attrs=$(aws --endpoint-url="$ENDPOINT_URL" sqs get-queue-attributes \
        --queue-url "$QUEUE_URL" \
        --attribute-names ApproximateNumberOfMessages ApproximateNumberOfMessagesNotVisible \
        --region "$REGION" \
        --output json)

    messages=$(echo "$attrs" | jq -r '.Attributes.ApproximateNumberOfMessages')
    inflight=$(echo "$attrs" | jq -r '.Attributes.ApproximateNumberOfMessagesNotVisible')

    echo "待機中メッセージ: $messages"
    echo "処理中メッセージ: $inflight"
}

# メイン処理
main() {
    case "${1:-}" in
        --help|-h)
            show_help
            exit 0
            ;;
        --batch)
            if [ -z "${2:-}" ]; then
                echo "エラー: --batch オプションには件数を指定してください"
                exit 1
            fi
            send_batch_messages "$2"
            ;;
        *)
            send_single_message "$1"
            ;;
    esac

    show_queue_status
}

main "$@"
