package main

import "fmt"

// ----------------------------------------------------------------------------
// Imagine the following is a generic stream processing package

// Note: using an interface as a marker where the set of types are structurally
// unlike but form a distinct set. Kind of like variant types.

type SourceMessage interface {
	// See https://golang.org/doc/faq#guarantee_satisfies_interface
	ImplementsSourceMessage()
}
type SinkMessage interface {
	ImplementsSinkMessage()
}

type EndOfStream struct{}

func (EndOfStream) ImplementsSourceMessage() {}

// Note: using function signatures as "interfaces of one". We
// can pass either free functions or methods bound to stateful objects.

type SourceProcessor func() SourceMessage
type SinkProcessor func(SinkMessage)
type Processor func(SourceMessage) (SinkMessage, error)

// Note: the types defined above are the duck-types for this
// function.
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

type OrderMessage struct {
	OrderID   ID
	ItemCosts []Currency
}

func (OrderMessage) ImplementsSourceMessage() {}

type LedgerMessage struct {
	OrderID ID
	Total   Currency
}

func (LedgerMessage) ImplementsSinkMessage() {}

func (lm LedgerMessage) String() string {
	return fmt.Sprintf(
		"ledger for order %d with total £%d.%d",
		lm.OrderID,
		lm.Total/100,
		lm.Total%100,
	)
}

type OrderSourceProcessor struct {
	offset   int
	messages []OrderMessage
}

func NewOrderSourceProcessor(messages []OrderMessage) SourceProcessor {
	p := OrderSourceProcessor{-1, messages}

	// Note: returns a bound method value
	return p.Next
}

// Note: Next satisfies the SourceProcessor signature
func (o *OrderSourceProcessor) Next() SourceMessage {
	o.offset++
	if o.offset >= len(o.messages) {
		return EndOfStream{}
	}

	return o.messages[o.offset]
}

func NewStdoutSinkProcessor() SinkProcessor {
	// Note: high order functions
	return func(message SinkMessage) {
		fmt.Printf("Sinking %v\n", message)
	}
}

func NewOrderToLedgerStreamProcessor() Processor {
	// Note: high order functions
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

	// Note: you can pass either functions or bound method values
	// for function typed arguments
	ProcessStream(
		source,    // this is a bound method value
		processor, // this is an anonymous function
		sink,
	)
}
