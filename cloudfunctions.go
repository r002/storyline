package cloudfunctions

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/r002/storyline-api/fbservices"
)

// UpdateMemberMetricsTestGet is an HTTP Cloud Function.
func UpdateMemberMetricsTestGet(w http.ResponseWriter, r *http.Request) {
	// s := fbservices.UpdateMemberMetrics()
	fmt.Fprint(w, ">> Testing GCF UpdateMemberMetricsTestGet")
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// Subscribes to topic and performs `member metrics update` on every message received
func UpdateMemberMetrics(ctx context.Context, m PubSubMessage) error {
	payload := string(m.Data)
	log.Printf(">> GCF: UpdateMemberMetrics: payload, %s", payload)
	s := fbservices.DoNightlyMetricsUpdate()
	log.Println(">> Job triggered: UpdateMemberMetrics", s)
	return nil
}
