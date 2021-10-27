package main

import (
	"context"
	"math"
	"os"
	"strconv"
	"time"

	emailService "carrier.microservices.go/src/lib/email"
	"carrier.microservices.go/src/lib/store"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

// EmailQueue ...
func EmailQueue(ctx context.Context, cloudWatchEvent events.CloudWatchEvent) {
	var emailExchange emailService.EmailExchange
	var err error

	logger.Debugf("CloudWatch event: EmailQueue: %+v", cloudWatchEvent)

	limit := 25
	counter := 0
	attemptLimit := 5
	continueLoop := true
	exchangeInitialized := false

	// get email repository
	emailRepository := NewEmailRepository(store.NewDynamoDBTable(db, os.Getenv("EMAILS_TABLE")))

	// main loop
	for continueLoop && counter < limit {

		counter++

		// get exchange from context if not initialized
		if emailExchange == nil {
			emailExchange = &emailService.SparkPostExchange{}
			err = emailExchange.Init()
			if err != nil {
				logger.Errorf("Cannot create email exchange: %s\n", err)
				return
			}
			logger.Debugw("Initialized email exchange")
			exchangeInitialized = true
		}

		// get email from top of queue and try to send
		if exchangeInitialized {

			// retrieve the first queued email (as list)
			emails, err := emailRepository.List(
				1,
				1,
				map[string]interface{}{
					"index": os.Getenv("EMAIL_QUEUE_INDEX"),
					"query": "send_status = :send_status",
					"expressionAttributeValues": map[string]*dynamodb.AttributeValue{
						":send_status": {
							N: aws.String(strconv.Itoa(EmailStatusQueued)),
						},
					},
				},
			)
			if err != nil {
				logger.Errorf("List queued emails error: %v", err)
				return
			}

			// if have emails in queue send
			if len(emails) > 0 {

				var emailIDs []uuid.UUID
				for _, email := range emails {
					emailIDs = append(emailIDs, email.ID)
				}
				logger.Debugf("Queued emails: (%d) %v", len(emailIDs), emailIDs)

				// loop over queued emails and send
				for _, email := range emails {

					// is it time to send this email?
					if time.Now().After(email.Queued) {

						// set email status to processing
						err = emailRepository.Update(email, store.ChangeSet{"send_status": EmailStatusProcessing})
						if err != nil {
							logger.Errorf("Unable to update email: %v", err)
						}

						// send email
						if sent := SendEmail(emailExchange, email, emailRepository); !sent {

							if email.Attempts >= attemptLimit {

								// failed too many times, do not attempt again
								err = emailRepository.Update(email, store.ChangeSet{
									"send_status": EmailStatusFailed,
									"queued":      time.Time{},
								})
								if err != nil {
									logger.Errorf("Unable to update email: %v", err)
								}
							} else {

								// update `queued` attribute with new date to push it back in the queue
								err = emailRepository.Update(email, store.ChangeSet{
									"queued": nextAttemptDate(email),
								})
								if err != nil {
									logger.Errorf("Unable to update email: %v", err)
								}
							}
						}
					} else {
						continueLoop = false // remaining emails in queue are not scheduled yet
						break
					}
				}
			} else {
				continueLoop = false // no more emails in queue, end loop
			}
		} else {
			continueLoop = false // email exchange not initialized unexpectedly
		}
	}
}

// nextAttemptDate generates the the next time to attempt a send using a backoff algorithm
func nextAttemptDate(email *Email) time.Time {
	return email.Queued.Add(time.Minute * time.Duration(math.Pow(2, float64(email.Attempts))))
}
