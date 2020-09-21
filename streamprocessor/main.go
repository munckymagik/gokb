package main

import "fmt"

// ----------------------------------------------------------------------------
// Imagine the following is a generic stream processing package

// Thing to note 01: using an interface as a marker. This allows us to
// group unlike variants without having to wrap them in `interface{}`. Kind of
// like variant types.
type SourceMessage interface {
	ImplementsSourceMessage()
}
type SinkMessage interface {
	ImplementsSinkMessage()
}

// Thing to note 02: using function signatures as an "interface of one". We
// can pass either free functions or methods bound to stateful objects.
type SourceProcessor func() SourceMessage
type SinkProcessor func(SinkMessage)
type Processor func(SourceMessage) (SinkMessage, error)

type EndOfStream struct{}

func (EndOfStream) ImplementsSourceMessage() {}

// Thing to note 03: making the client code the owner of the interfaces rather
// than the concrete types. Leads to many smaller interfaces, and simple
// mocks in tests.
func ProcessStream(next SourceProcessor, process Processor, sink SinkProcessor) {
	for {
		message := next()
		if _, ok := message.(EndOfStream); ok {
			fmt.Println("Got end of stream, terminating processing.")
			break
		}

		outMessage, err := process(message)
		if err != nil {
			fmt.Printf("Error processing message: %v, error was: %v", message, err)
			continue
		}

		sink(outMessage)
	}
}

// ----------------------------------------------------------------------------
// Imagine the following is a package that contains business stream processing
// logic

type ID = int64
type Currency = int64

// OrderMessage is our input message type
type OrderMessage struct {
	OrderID   ID
	ItemCosts []Currency
}

// A no-op that tags OrderMessage as a SourceMessage
func (OrderMessage) ImplementsSourceMessage() {}

// LedgerMessage is our output message type
type LedgerMessage struct {
	OrderID ID
	Total   Currency
}

// A no-op that tags LedgerMessage as a SinkMessage
func (LedgerMessage) ImplementsSinkMessage() {}

func (lm LedgerMessage) String() string {
	return fmt.Sprintf(
		"ledger for order %d with total £%d.%d",
		lm.OrderID,
		lm.Total/100,
		lm.Total%100,
	)
}

// OrderSourceProcessor knows how to stream OrderMessages
type OrderSourceProcessor struct {
	offset   int
	messages []OrderMessage
}

func NewOrderSourceProcessor(messages []OrderMessage) OrderSourceProcessor {
	return OrderSourceProcessor{-1, messages}
}

// Next satisfies the SourceProcessor signature
func (o *OrderSourceProcessor) Next() SourceMessage {
	o.offset++
	if o.offset >= len(o.messages) {
		return EndOfStream{}
	}

	return o.messages[o.offset]
}

// NewStdoutSinkProcessor knows how to sink all message types to STDOUT.
// It returns a function that satisfies the SinkProcessor signature.
func NewStdoutSinkProcessor() SinkProcessor {
	return func(message SinkMessage) {
		fmt.Printf("Sinking %v\n", message)
	}
}

// NewOrderToLedgerStreamProcessor creates a Processor that knows how to
// produce LedgerMessage sink messages from OrderMessage source messages.
func NewOrderToLedgerStreamProcessor() Processor {
	return func(message SourceMessage) (SinkMessage, error) {
		switch value := message.(type) {
		case OrderMessage:
			var total Currency = 0
			for _, amount := range value.ItemCosts {
				total += amount
			}
			newMessage := LedgerMessage{
				value.OrderID,
				total,
			}

			return newMessage, nil
		default:
			return nil, fmt.Errorf("Got unexpected message %v", value)
		}
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

	// Thing to note 04: you can pass either functions or bound method values
	// for function typed arguments
	ProcessStream(
		source.Next, // passing a bound method value
		processor,   // passing a function
		sink,        // ditto
	)
}
