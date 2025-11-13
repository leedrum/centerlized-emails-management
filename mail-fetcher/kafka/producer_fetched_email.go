package kafka

func PublishFetchedEmailProducer(producer *Producer, tenant, s3_key, file string) {
	data := map[string]interface{}{
		"tenant_id": tenant,
		"s3_key":    s3_key,
		"file":      file,
	}
	producer.Publish("email_processor.fechted_email", data)
}
