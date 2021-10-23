package main

import (
	"context"
	"os"

	emailService "carrier.microservices.go/src/lib/email"
	"carrier.microservices.go/src/lib/store"
	"github.com/aws/aws-lambda-go/events"
)

// EmailQueue ...
func EmailQueue(ctx context.Context, cloudWatchEvent events.CloudWatchEvent) {
	var emailExchange emailService.EmailExchange

	logger.Debugf("CloudWatch event: EmailQueue: %+v", cloudWatchEvent)

	// get email repository
	emailRepository := NewEmailRepository(store.NewDynamoDBTable(db, os.Getenv("EMAILS_TABLE")))

	// retrieve a list of queued emails
	emails, err := emailRepository.List(1, 50, map[string]string{"index": os.Getenv("EMAIL_QUEUE_INDEX")})
	if err != nil {
		logger.Errorf("List queued emails error: %v", err)
		return
	}

	logger.Debugf("Queued emails: %v", emails)

	// if have emails in queue send
	if len(emails) > 0 {
		emailExchange = &emailService.SparkPostExchange{}
		err = emailExchange.Init()
		if err != nil {
			logger.Errorf("Cannot create email exchange: %s\n", err)
			return
		}

		logger.Debugw("Initialized email exchange")

		// loop over queued emails and send
		for _, email := range emails {
			SendEmail(emailExchange, email, emailRepository)
		}
	}
}
