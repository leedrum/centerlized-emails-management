# Mail Fetcher Service
- Can be scaled by tenant (check the `env` file)

## 1. Kafka
- Receive message fetch email from `kafka` -> fetch email -> store on `s3` (implement later).
- Currently, fetch email by `tenants` in `env`
- Receive message token update from `kafka` -> update data into postgres
- Produce the message after fetching email and storing on `s3`
### Consumer
- Fetch email
- Token update: topic is `mail_fetcher.token_update`
### Producer
- fetched email: topic is `email_processor.fechted_email`

## 2. TODO
- run for all tenant if config is `all`
- migration
- API document
