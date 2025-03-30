package slack

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendToSlack(t *testing.T) {
	// テスト用のWebhookサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストメソッドの確認
		assert.Equal(t, "POST", r.Method)
		// Content-Typeの確認
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// レスポンスを返す
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// テストケース
	tests := []struct {
		name    string
		message string
		wantErr bool
	}{
		{
			name:    "正常なメッセージ",
			message: "テストメッセージ",
			wantErr: false,
		},
		{
			name:    "空のメッセージ",
			message: "",
			wantErr: false,
		},
		{
			name:    "長いメッセージ",
			message: "テストメッセージ".repeat(1000),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendToSlack(server.URL, tt.message)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
