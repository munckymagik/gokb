package main

import "fmt"

type ID = int64
type Currency = int64

// Message defines a variant type for input messages
type Message interface {
	ImplementsMessage()
}
type EndOfStream struct{}
type OrderMessage struct {
	OrderID   ID
	ItemCosts []Currency
}

func (EndOfStream) ImplementsMessage()  {}
func (OrderMessage) ImplementsMessage() {}

// OutMessage defines a variant type for output messages
type OutMessage interface {
	ImplementsOutMessage()
}

type LedgerMessage struct {
	OrderID ID
	Total   Currency
}

func (LedgerMessage) ImplementsOutMessage() {}
func (lm LedgerMessage) String() string {
	return fmt.Sprintf(
		"ledger for order %d with total £%d.%d",
		lm.OrderID,
		lm.Total/100,
		lm.Total%100,
	)
}

// The stream processing abstractions

type SourceProcessor func() Message
type SinkProcessor func(OutMessage)
type Processor func(Message) OutMessage

type OrderSourceProcessor struct {
	offset   int
	messages []OrderMessage
}

// Concrete stream processing types

func NewOrderSourceProcessor(messages []OrderMessage) OrderSourceProcessor {
	return OrderSourceProcessor{-1, messages}
}

func (o *OrderSourceProcessor) Next() Message {
	o.offset++
	if o.offset >= len(o.messages) {
		return EndOfStream{}
	}

	return o.messages[o.offset]
}

func NewStdoutSinkProcessor() SinkProcessor {
	return func(message OutMessage) {
		fmt.Printf("Sinking %v\n", message)
	}
}

func NewOrderToLedgerStreamProcessor() Processor {
	return func(message Message) OutMessage {
		switch value := message.(type) {
		case OrderMessage:
			var total Currency = 0
			for _, amount := range value.ItemCosts {
				total += amount
			}
			return LedgerMessage{
				value.OrderID,
				total,
			}
		default:
			// TODO not great error handling
			panic(fmt.Sprintf("Got unexpected message %v", value))
		}
	}
}

// processStream is the stream procesing runtime
func processStream(next SourceProcessor, process Processor, sink SinkProcessor) {
	for {
		message := next()
		if _, ok := message.(EndOfStream); ok {
			fmt.Println("Got end of stream, terminating processing.")
			break
		}

		outMessage := process(message)
		sink(outMessage)
	}
}

func main() {
	orderMessages := []OrderMessage{
		{1234, []Currency{399, 1275, 1275}},
		{1235, []Currency{1099, 1175}},
		{1236, []Currency{799, 1000}},
	}

	source := NewOrderSourceProcessor(orderMessages)
	processor := NewOrderToLedgerStreamProcessor()
	sink := NewStdoutSinkProcessor()

	processStream(source.Next, processor, sink)
}
