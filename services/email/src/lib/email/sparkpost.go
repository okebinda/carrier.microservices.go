package mail

import (
	"os"
	"strconv"
	"time"

	sp "github.com/SparkPost/gosparkpost"
)

// SparkPostExchange defines a SparkPost service
type SparkPostExchange struct {
	Client sp.Client
}

// Init initializes the SparkPost service
func Init(ex *SparkPostExchange) error {

	// get SparkPost configuration from ENV
	APIKey := os.Getenv("SPARKPOST_API_KEY")
	APIBaseURL := os.Getenv("SPARKPOST_BASE_URL")
	APIVersion, _ := strconv.Atoi(os.Getenv("SPARKPOST_API_VERSION"))

	// create new SparkPost client
	cfg := &sp.Config{
		BaseUrl:    APIBaseURL,
		ApiKey:     APIKey,
		ApiVersion: APIVersion,
	}
	return ex.Client.Init(cfg)
}

// Send sends an email through the service
func (ex *SparkPostExchange) Send(email *Email) error {

	// create recipient list
	recipients := []sp.Recipient{}
	for _, address := range email.Recipients {
		recipient := sp.Recipient{
			Address:          address,
			SubstitutionData: email.Substitutions,
		}
		recipients = append(recipients, recipient)
	}

	// send email
	email.LastAttemptAt = time.Now()
	tx := &sp.Transmission{
		Recipients: recipients,
		Content: map[string]interface{}{
			"template_id": email.Template,
		},
	}
	id, res, err := ex.Client.Send(tx)
	if err != nil {
		return err
	}

	txResults := res.Results.(map[string]interface{})
	email.ID = id
	email.Accepted = int(txResults["total_accepted_recipients"].(float64))
	email.Rejected = int(txResults["total_rejected_recipients"].(float64))

	return nil
}
