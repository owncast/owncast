package webhooks

func SendStreamStatusEvent(eventType EventType) {
	SendEventToWebhooks(WebhookEvent{Type: eventType})
}
