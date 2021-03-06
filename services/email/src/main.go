package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"carrier.microservices.go/src/lib/store"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	chiproxy "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger
var adapter *chiproxy.ChiLambda
var db *dynamodb.DynamoDB

func init() {
	var err error

	// connect to datastore
	db, err = store.CreateConnection(os.Getenv("AWS_REGION"), os.Getenv("DYNAMODB_ENDPOINT"))
	if err != nil {
		log.Fatalf("Database connection error: %s", err)
	}

	// create router
	r := chi.NewRouter()

	// add middleware
	r.Use(LogRequest)
	r.Use(Authorize)
	r.Use(EmailRepositoryCtx)
	r.Use(EmailExchangeCtx)

	// add routes
	r.Route("/email/{emailID}", func(r chi.Router) {
		r.Use(EmailCtx)
		r.Get("/", GetEmail)
		r.Put("/", UpdateEmail)
		r.Delete("/", DeleteEmail)
	})
	r.Get("/emails", GetEmails)
	r.Post("/emails", PostEmails)

	adapter = chiproxy.New(r)
}

// APIGatewayHandler is the lambda handler invoked by API Gateway events
func APIGatewayHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// initialize logger
	lc, _ := lambdacontext.FromContext(ctx)
	logger = sugaredLogger(lc.AwsRequestID)
	defer logger.Sync()

	// serve request
	c, err := adapter.ProxyWithContext(ctx, request)
	return c, err
}

// CloudWatchHandler is the lambda handler invoked by CloudWatch events
func CloudWatchHandler(ctx context.Context, cloudWatchEvent events.CloudWatchEvent) {

	// initialize logger
	lc, _ := lambdacontext.FromContext(ctx)
	logger = sugaredLogger(lc.AwsRequestID)
	defer logger.Sync()

	// run job
	EmailQueue(ctx, cloudWatchEvent)
}

// sugaredLogger initializes the zap sugar logger
func sugaredLogger(requestID string) *zap.SugaredLogger {
	logConfig := []byte(fmt.Sprintf(`{
			"level": "%s",
			"encoding": "%s",
			"outputPaths": ["stdout"],
			"errorOutputPaths": ["stderr"],
			"encoderConfig": {
				"messageKey": "message",
				"levelKey": "level",
				"levelEncoder": "lowercase"
			}
		}`,
		os.Getenv("LOG_LEVEL"),
		os.Getenv("LOG_ENCODING"),
	))

	var cfg zap.Config
	if err := json.Unmarshal(logConfig, &cfg); err != nil {
		panic(err)
	}
	zapLogger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync()

	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	return zapLogger.
		With(zap.Field{Key: "request_id", Type: zapcore.StringType, String: requestID}).
		With(zap.Field{Key: "timestamp", Type: zapcore.StringType, String: time.Now().Format("2006-01-02T15:04:05.999999999Z07:00")}).
		Sugar()
}

// successResponse generates a success (200) response
func successResponse(w http.ResponseWriter, code int, fields interface{}) {
	if fields == nil {
		generateResponse(w, code, nil)
	} else {
		body, err := json.Marshal(fields)
		if err != nil {
			logger.Errorf("Marshalling error: %s", err)
			serverErrorResponse(w)
		}
		generateResponse(w, code, body)
	}
}

// userErrorResponse generates a user error (400) response
func userErrorResponse(w http.ResponseWriter, code int, errorMessage string) {
	body, err := json.Marshal(map[string]interface{}{
		"error": errorMessage,
	})
	if err != nil {
		logger.Errorf("Marshalling error: %s", err)
		serverErrorResponse(w)
	}
	generateResponse(w, code, body)
}

// serverErrorResponse generates a server error (500) response
func serverErrorResponse(w http.ResponseWriter) {
	generateResponse(w, 500, []byte("{\"error\":\"Server error\"}"))
}

// generateResponse generates an HTTP JSON Lambda response to return to the user
func generateResponse(w http.ResponseWriter, statusCode int, body []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		logger.Errorf("Error writing response: %s", err)
	}
}

// Handler handles AWS Lambda events
type Handler struct{}

// Invoke handles a generic AWS Lambda invocation event
func (h Handler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {

	// parse event payload so we can determine what type of event this is
	var eventData interface{}
	if err := json.Unmarshal(payload, &eventData); err != nil {
		logger.Errorf("Unmarshalling error: %s", err)
		return nil, nil
	}
	eventDataMap := eventData.(map[string]interface{})

	// event data shape looks like APIGatewayProxyRequest
	_, pathExists := eventDataMap["path"]
	_, httpMethodExists := eventDataMap["httpMethod"]
	if pathExists && httpMethodExists {
		event := new(events.APIGatewayProxyRequest)
		if err := json.Unmarshal(payload, event); err != nil {
			logger.Errorf("Unmarshalling error: %s", err)
			return nil, nil
		}
		response, err := APIGatewayHandler(ctx, *event)
		responseBytes, err := json.Marshal(response)
		if err != nil {
			logger.Errorf("Marshalling error: %s", err)
			return nil, nil
		}
		return responseBytes, err
	}

	// event data shape looks like CloudWatchEvent
	_, IDExists := eventDataMap["id"]
	_, sourceExists := eventDataMap["source"]
	if IDExists && sourceExists {
		event := new(events.CloudWatchEvent)
		if err := json.Unmarshal(payload, event); err != nil {
			logger.Errorf("Unmarshalling error: %s", err)
			return nil, nil
		}
		CloudWatchHandler(ctx, *event)
		return nil, nil
	}

	return nil, nil
}

func main() {
	lambda.StartHandler(Handler{})
}
