/*
Copyright 2017 Heptio Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sinks

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/viper"
	"k8s.io/api/core/v1"
)

// StdoutSink is the other basic sink
// By default, Fluentd/ElasticSearch won't index glog formatted lines
// By logging raw JSON to stdout, we will get automated indexing which
// can be queried in Kibana.
type FileSink struct {
	// TODO: create a channel and buffer for scaling
}

// NewStdoutSink will create a new StdoutSink with default options, returned as
// an EventSinkInterface
func NewFileSink() EventSinkInterface {
	return &FileSink{}
}


type EventMessage struct {
  Item string
  Timestamp string
  Reason string
  Message string
  RawEventStream string
}

// UpdateEvents implements the EventSinkInterface
func (gs *FileSink) UpdateEvents(eNew *v1.Event, eOld *v1.Event) {
        filename := viper.GetString("filesinkpath")
	
        var file, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if isError(err) { os.Exit(1) }
        defer file.Close()

	eData := NewEventData(eNew, eOld)


	eJSONBytes, err := json.Marshal(eData)
	if err == nil {
		if isError(err) { os.Exit(1) }
	} else {
		fmt.Fprintf(os.Stderr, "Failed to json serialize event: %v", err)
	}


        msg := &EventMessage{
                Item: eNew.InvolvedObject.Name,
                Reason: eNew.Reason,
                Message: eNew.Message,
                Timestamp: eNew.LastTimestamp.String(),
                RawEventStream: string(eJSONBytes),
        }


	msgJSONBytes, err := json.Marshal(msg)
	if err == nil {
		if isError(err) { os.Exit(1) }
	} else {
		fmt.Fprintf(os.Stderr, "Failed to json serialize event: %v", err)
	}

	_, err = file.WriteString(string(msgJSONBytes))
	if isError(err) { os.Exit(1) }
	err = file.Sync()
        if isError(err) { os.Exit(1) }
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}
