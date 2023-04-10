package prometheus

import (
	"context"
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/logger"
	"go.unistack.org/micro/v3/server"
)

var (
	// default metric prefix
	DefaultMetricPrefix = "micro_"
	// default label prefix
	DefaultLabelPrefix = "micro_"

	clientOpsCounter           *prometheus.CounterVec
	clientTimeCounterSummary   *prometheus.SummaryVec
	clientTimeCounterHistogram *prometheus.HistogramVec

	serverOpsCounter           *prometheus.CounterVec
	serverTimeCounterSummary   *prometheus.SummaryVec
	serverTimeCounterHistogram *prometheus.HistogramVec

	publishOpsCounter           *prometheus.CounterVec
	publishTimeCounterSummary   *prometheus.SummaryVec
	publishTimeCounterHistogram *prometheus.HistogramVec

	subscribeOpsCounter           *prometheus.CounterVec
	subscribeTimeCounterSummary   *prometheus.SummaryVec
	subscribeTimeCounterHistogram *prometheus.HistogramVec

	mu sync.Mutex
)

type Options struct {
	Name    string
	Version string
	ID      string
	Context context.Context
}

type Option func(*Options)

func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

func ServiceName(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func ServiceVersion(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

func ServiceID(id string) Option {
	return func(o *Options) {
		o.ID = id
	}
}

func registerServerMetrics(ctx context.Context) {
	mu.Lock()
	defer mu.Unlock()

	if serverOpsCounter == nil {
		serverOpsCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%sserver_request_total", DefaultMetricPrefix),
				Help: "Requests processed, partitioned by endpoint and status",
			},
			[]string{
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "name"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "version"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "id"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "endpoint"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "status"),
			},
		)
	}

	if serverTimeCounterSummary == nil {
		serverTimeCounterSummary = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: fmt.Sprintf("%sserver_latency_microseconds", DefaultMetricPrefix),
				Help: "Request latencies in microseconds, partitioned by endpoint",
			},
			[]string{
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "name"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "version"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "id"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "endpoint"),
			},
		)
	}

	if serverTimeCounterHistogram == nil {
		serverTimeCounterHistogram = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: fmt.Sprintf("%sserver_request_duration_seconds", DefaultMetricPrefix),
				Help: "Request time in seconds, partitioned by endpoint",
			},
			[]string{
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "name"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "version"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "id"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "endpoint"),
			},
		)
	}

	for _, collector := range []prometheus.Collector{serverOpsCounter, serverTimeCounterSummary, serverTimeCounterHistogram} {
		if err := prometheus.DefaultRegisterer.Register(collector); err != nil {
			// if already registered, skip fatal
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				logger.Fatal(ctx, err.Error())
			}
		}
	}

}

func registerPublishMetrics(ctx context.Context) {
	mu.Lock()
	defer mu.Unlock()

	if publishOpsCounter == nil {
		publishOpsCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%spublish_message_total", DefaultMetricPrefix),
				Help: "Messages sent, partitioned by endpoint and status",
			},
			[]string{
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "name"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "version"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "id"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "endpoint"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "status"),
			},
		)
	}

	if publishTimeCounterSummary == nil {
		publishTimeCounterSummary = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: fmt.Sprintf("%spublish_message_latency_microseconds", DefaultMetricPrefix),
				Help: "Message latencies in microseconds, partitioned by endpoint",
			},
			[]string{
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "name"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "version"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "id"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "endpoint"),
			},
		)
	}

	if publishTimeCounterHistogram == nil {
		publishTimeCounterHistogram = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: fmt.Sprintf("%spublish_message_duration_seconds", DefaultMetricPrefix),
				Help: "Message publish time in seconds, partitioned by endpoint",
			},
			[]string{
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "name"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "version"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "id"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "endpoint"),
			},
		)
	}

	for _, collector := range []prometheus.Collector{publishOpsCounter, publishTimeCounterSummary, publishTimeCounterHistogram} {
		if err := prometheus.DefaultRegisterer.Register(collector); err != nil {
			// if already registered, skip fatal
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				logger.Fatal(ctx, err.Error())
			}
		}
	}

}

func registerSubscribeMetrics(ctx context.Context) {
	mu.Lock()
	defer mu.Unlock()

	if subscribeOpsCounter == nil {
		subscribeOpsCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%ssubscribe_message_total", DefaultMetricPrefix),
				Help: "Messages processed, partitioned by endpoint and status",
			},
			[]string{
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "name"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "version"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "id"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "endpoint"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "status"),
			},
		)
	}

	if subscribeTimeCounterSummary == nil {
		subscribeTimeCounterSummary = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: fmt.Sprintf("%ssubscribe_message_latency_microseconds", DefaultMetricPrefix),
				Help: "Message processing latencies in microseconds, partitioned by endpoint",
			},
			[]string{
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "name"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "version"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "id"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "endpoint"),
			},
		)
	}

	if subscribeTimeCounterHistogram == nil {
		subscribeTimeCounterHistogram = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: fmt.Sprintf("%ssubscribe_message_duration_seconds", DefaultMetricPrefix),
				Help: "Request time in seconds, partitioned by endpoint",
			},
			[]string{
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "name"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "version"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "id"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "endpoint"),
			},
		)
	}

	for _, collector := range []prometheus.Collector{subscribeOpsCounter, subscribeTimeCounterSummary, subscribeTimeCounterHistogram} {
		if err := prometheus.DefaultRegisterer.Register(collector); err != nil {
			// if already registered, skip fatal
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				logger.Fatal(ctx, err.Error())
			}
		}
	}

}

func registerClientMetrics(ctx context.Context) {
	mu.Lock()
	defer mu.Unlock()

	if clientOpsCounter == nil {
		clientOpsCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("%srequest_total", DefaultMetricPrefix),
				Help: "Requests processed, partitioned by endpoint and status",
			},
			[]string{
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "name"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "version"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "id"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "endpoint"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "status"),
			},
		)
	}

	if clientTimeCounterSummary == nil {
		clientTimeCounterSummary = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: fmt.Sprintf("%slatency_microseconds", DefaultMetricPrefix),
				Help: "Request latencies in microseconds, partitioned by endpoint",
			},
			[]string{
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "name"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "version"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "id"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "endpoint"),
			},
		)
	}

	if clientTimeCounterHistogram == nil {
		clientTimeCounterHistogram = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: fmt.Sprintf("%srequest_duration_seconds", DefaultMetricPrefix),
				Help: "Request time in seconds, partitioned by endpoint",
			},
			[]string{
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "name"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "version"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "id"),
				fmt.Sprintf("%s%s", DefaultLabelPrefix, "endpoint"),
			},
		)
	}

	for _, collector := range []prometheus.Collector{clientOpsCounter, clientTimeCounterSummary, clientTimeCounterHistogram} {
		if err := prometheus.DefaultRegisterer.Register(collector); err != nil {
			// if already registered, skip fatal
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				logger.Fatal(ctx, err.Error())
			}
		}
	}

}

type wrapper struct {
	options  Options
	callFunc client.CallFunc
	client.Client
}

func NewClientWrapper(opts ...Option) client.Wrapper {
	options := Options{Context: context.Background()}
	for _, o := range opts {
		o(&options)
	}

	registerClientMetrics(options.Context)
	registerPublishMetrics(options.Context)

	return func(c client.Client) client.Client {
		handler := &wrapper{
			options: options,
			Client:  c,
		}

		return handler
	}
}

func NewCallWrapper(opts ...Option) client.CallWrapper {
	options := Options{Context: context.Background()}
	for _, o := range opts {
		o(&options)
	}

	registerClientMetrics(options.Context)

	return func(fn client.CallFunc) client.CallFunc {
		handler := &wrapper{
			options:  options,
			callFunc: fn,
		}

		return handler.CallFunc
	}
}

func (w *wrapper) CallFunc(ctx context.Context, addr string, req client.Request, rsp interface{}, opts client.CallOptions) error {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())

	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		clientTimeCounterSummary.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint).Observe(us)
		clientTimeCounterHistogram.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint).Observe(v)
	}))
	defer timer.ObserveDuration()

	err := w.callFunc(ctx, addr, req, rsp, opts)
	if err == nil {
		clientOpsCounter.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint, "success").Inc()
	} else {
		clientOpsCounter.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint, "failure").Inc()
	}

	return err

}

func (w *wrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())

	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		clientTimeCounterSummary.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint).Observe(us)
		clientTimeCounterHistogram.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint).Observe(v)
	}))
	defer timer.ObserveDuration()

	err := w.Client.Call(ctx, req, rsp, opts...)
	if err == nil {
		clientOpsCounter.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint, "success").Inc()
	} else {
		clientOpsCounter.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint, "failure").Inc()
	}

	return err
}

func (w *wrapper) Stream(ctx context.Context, req client.Request, opts ...client.CallOption) (client.Stream, error) {
	endpoint := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())

	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		clientTimeCounterSummary.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint).Observe(us)
		clientTimeCounterHistogram.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint).Observe(v)
	}))
	defer timer.ObserveDuration()

	stream, err := w.Client.Stream(ctx, req, opts...)
	if err == nil {
		clientOpsCounter.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint, "success").Inc()
	} else {
		clientOpsCounter.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint, "failure").Inc()
	}

	return stream, err
}

func (w *wrapper) Publish(ctx context.Context, p client.Message, opts ...client.PublishOption) error {
	endpoint := p.Topic()

	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		publishTimeCounterSummary.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint).Observe(us)
		publishTimeCounterHistogram.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint).Observe(v)
	}))
	defer timer.ObserveDuration()

	err := w.Client.Publish(ctx, p, opts...)
	if err == nil {
		publishOpsCounter.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint, "success").Inc()
	} else {
		publishOpsCounter.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint, "failure").Inc()
	}

	return err
}

func NewHandlerWrapper(opts ...Option) server.HandlerWrapper {
	options := Options{Context: context.Background()}
	for _, o := range opts {
		o(&options)
	}
	registerServerMetrics(options.Context)

	handler := &wrapper{
		options: options,
	}

	return handler.HandlerFunc
}

func (w *wrapper) HandlerFunc(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		endpoint := req.Endpoint()

		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			us := v * 1000000 // make microseconds
			serverTimeCounterSummary.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint).Observe(us)
			serverTimeCounterHistogram.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint).Observe(v)
		}))
		defer timer.ObserveDuration()

		err := fn(ctx, req, rsp)
		if err == nil {
			serverOpsCounter.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint, "success").Inc()
		} else {
			serverOpsCounter.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint, "failure").Inc()
		}

		return err
	}
}

func NewSubscriberWrapper(opts ...Option) server.SubscriberWrapper {
	options := Options{Context: context.Background()}
	for _, o := range opts {
		o(&options)
	}

	registerSubscribeMetrics(options.Context)

	handler := &wrapper{
		options: options,
	}

	return handler.SubscriberFunc
}

func (w *wrapper) SubscriberFunc(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Message) error {
		endpoint := msg.Topic()

		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			us := v * 1000000 // make microseconds
			subscribeTimeCounterSummary.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint).Observe(us)
			subscribeTimeCounterHistogram.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint).Observe(v)
		}))
		defer timer.ObserveDuration()

		err := fn(ctx, msg)
		if err == nil {
			subscribeOpsCounter.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint, "success").Inc()
		} else {
			subscribeOpsCounter.WithLabelValues(w.options.Name, w.options.Version, w.options.ID, endpoint, "failure").Inc()
		}

		return err
	}
}
