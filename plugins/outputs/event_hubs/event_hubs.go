//go:generate ../../../tools/readme_config_includer/generator
package event_hubs

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/serializers"
)

//go:embed sample.conf
var sampleConfig string

type EventHubs struct {
	Log              telegraf.Logger `toml:"-"`
	ConnectionString string          `toml:"connection_string"`
	Timeout          config.Duration `toml:"timeout"`
	PartitionKey     string          `toml:"partition_key"`
	MaxMessageSize   uint64          `toml:"max_message_size"`

	client       *azeventhubs.ProducerClient
	batchOptions azeventhubs.EventDataBatchOptions
	serializer   serializers.Serializer
}

const (
	defaultRequestTimeout = time.Second * 30
)

func (*EventHubs) SampleConfig() string {
	return sampleConfig
}

func (e *EventHubs) Init() error {
	if e.MaxMessageSize > 0 {
		e.batchOptions.MaxBytes = e.MaxMessageSize
	}

	if e.PartitionKey != "" {
		e.Log.Warn("The PartitionKey config option is current ignored in this build")
	}

	return nil
}

func (e *EventHubs) Connect() error {
	var err error
	producerOptions := azeventhubs.ProducerClientOptions{}
	// EventHub string is always emtpy, meaning it must be provided as part of
	// the connection string (e.g. EntityPath=<entity path>;)
	e.client, err = azeventhubs.NewProducerClientFromConnectionString(e.ConnectionString, "", &producerOptions)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventHubs) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(e.Timeout))
	defer cancel()

	err := e.client.Close(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (e *EventHubs) SetSerializer(serializer serializers.Serializer) {
	e.serializer = serializer
}

func (e *EventHubs) Write(metrics []telegraf.Metric) error {
	events := make([]*azeventhubs.EventData, 0, len(metrics))
	for _, metric := range metrics {
		payload, err := e.serializer.Serialize(metric)
		if err != nil {
			e.Log.Debugf("Could not serialize metric: %v", err)
			continue
		}

		event := azeventhubs.EventData{Body: payload}
		events = append(events, &event)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(e.Timeout))
	defer cancel()

	batch, err := e.client.NewEventDataBatch(ctx, &e.batchOptions)
	if err != nil {
		return err
	}

	for i := 0; i < len(events); i++ {
		if err = batch.AddEventData(events[i], nil); err != nil {
			return fmt.Errorf("failed to and event %s: %w", events[i], err)
		}
	}

	if batch.NumEvents() > 0 {
		if err := e.client.SendEventDataBatch(ctx, batch, nil); err != nil {
			return fmt.Errorf("failed to send batch: %w", err)
		}
	}

	return nil
}

func init() {
	outputs.Add("event_hubs", func() telegraf.Output {
		return &EventHubs{
			Timeout: config.Duration(defaultRequestTimeout),
		}
	})
}
