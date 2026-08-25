package browser

import "testing"

func TestParseAndValidateSubscription(t *testing.T) {
	valid := `{"endpoint":"https://93.184.216.34/push","keys":{"auth":"auth","p256dh":"key"}}`
	if _, err := ParseAndValidateSubscription(valid); err != nil {
		t.Fatalf("expected public HTTPS endpoint to be accepted: %v", err)
	}

	for _, subscription := range []string{
		`{"endpoint":"http://93.184.216.34/push"}`,
		`{"endpoint":"https://127.0.0.1/push"}`,
		`{"endpoint":"https://169.254.169.254/latest"}`,
		`{"endpoint":"https://100.64.0.1/push"}`,
		`{"endpoint":"https://[2002:a00:1::]/push"}`,
		`{"endpoint":"https://does-not-exist.invalid/push"}`,
		`not json`,
	} {
		t.Run(subscription, func(t *testing.T) {
			if _, err := ParseAndValidateSubscription(subscription); err == nil {
				t.Fatalf("expected subscription %q to be rejected", subscription)
			}
		})
	}
}
