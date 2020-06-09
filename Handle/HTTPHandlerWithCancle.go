package handle

import (
	"bytes"
	"context"
	"net/http"

	"cloud.google.com/go/pubsub"
)

//LogAPI is rest api handler for stream log
type LogAPI struct {
	Topic           *pubsub.Topic
	Context         context.Context
	Cancle          context.CancelFunc
	CancleCondition func(r *http.Request) bool
}

func (l LogAPI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if l.Context.Err() != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	if l.CancleCondition(r) {
		l.Cancle()
		return
	}

	var buff bytes.Buffer
	_, err := buff.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	l.Topic.Publish(l.Context, &pubsub.Message{Attributes: map[string]string{"Path": r.URL.Path}, Data: buff.Bytes()})
	w.WriteHeader(http.StatusOK)
}
