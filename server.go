package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/op/go-logging"
	redis "gopkg.in/redis.v4"
)

var log = logging.MustGetLogger("ibiza")
var serverName = "Ibiza v1.0 - A PubSubHubbub hub server"
var copyright = "(C) 2016 Hiranya Samarasekera. https://github.com/hiranya/ibiza"

const errorRequiredFieldMissingHubMode = "Required field missing: 'hub.mode'"
const errorRequiredFieldMissingHubCallback = "Required subscription request field missing: 'hub.callback'"
const errorRequiredFieldMissingHubTopic = "Required subscription request field missing: 'hub.topic'"

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379", //6379
	Password: "",               // no password set
	DB:       0,                // use default DB
})

// Subscription ...
type Subscription struct {
	Topic        string
	CallbackURL  url.URL
	LeaseSeconds int
	Secret       string
	From         string
}

var subscriptionStore *map[string]Subscription

func mainHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":
			{
				fmt.Fprintln(w, serverName)
				fmt.Fprintln(w, copyright)
				return
			}
		case "POST":
			{
				log.Debug("Form POST detected")
				hubMode := r.FormValue("hub.mode")

				if hubMode == "" {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintln(w, errorRequiredFieldMissingHubMode)
					return
				}

				if hubMode == "subscribe" {
					log.Debug("hub.mode = subscribe")
					err := subscribe(r)

					if err != nil {
						w.WriteHeader(http.StatusBadRequest)
						fmt.Fprintln(w, err)
						return
					}

					w.WriteHeader(http.StatusAccepted)

				} else if hubMode == "unsubscribe" {
					log.Debug("hub.mode = unsubscribe")
				} else if hubMode == "publish" {
					log.Debug("hub.mode = publish")
				} else {
					log.Debug("hub.mode command did not recognize")
				}
			}
		}

	})
}

func subscribe(r *http.Request) error {
	// validate for required fields
	hubCallback := r.FormValue("hub.callback")
	hubTopic := r.FormValue("hub.topic")

	var errorBuffer bytes.Buffer

	if hubCallback == "" {
		appendErrorString(&errorBuffer, errorRequiredFieldMissingHubCallback)
	}
	if hubTopic == "" {
		appendErrorString(&errorBuffer, errorRequiredFieldMissingHubTopic)
	}

	if errorBuffer.String() != "" {
		return errors.New(errorBuffer.String())
	}

	topicID := int64(redisClient.ZScore("topics", hubTopic).Val())
	if topicID == 0 {
		topicID = redisClient.Incr("id:topics").Val()
		redisClient.ZAdd("topics", redis.Z{Score: float64(topicID), Member: hubTopic})

		topicMap := make(map[string]string)
		topicMap["topic_url"] = hubTopic
		topicMap["created_timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
		redisClient.HMSet("topic:"+strconv.FormatInt(topicID, 10), topicMap)
	}

	subscriptionScore := redisClient.ZScore("subscriptions", hubCallback).Val()
	if subscriptionScore == 0 {
		subscriptionID := redisClient.Incr("id:subscriptions")
		redisClient.ZAdd("subscriptions", redis.Z{Score: float64(subscriptionID.Val()), Member: hubCallback})

		subscriptionMap := make(map[string]string)
		subscriptionMap["topic_id"] = strconv.FormatInt(topicID, 10)
		subscriptionMap["created_timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
		subscriptionMap["callback_url"] = hubCallback
		redisClient.HMSet("subscription:"+strconv.FormatInt(subscriptionID.Val(), 10), subscriptionMap)
	}

	return nil
}

func appendErrorString(errBuffer *bytes.Buffer, err string) {
	errBuffer.WriteString(err)
	errBuffer.WriteString("\n")
}

func main() {
	log.Info(serverName)
	log.Info(copyright)

	subscriptionStore = new(map[string]Subscription)

	http.HandleFunc("/", mainHandler())
	http.ListenAndServe(":9090", nil)
}
