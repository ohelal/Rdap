package errors

import (
	"fmt"
	"strings"
)

type ErrorChain struct {
	Root   *Error
	Chain  []*Error
	Causes map[string]int
}

func AnalyzeErrorChain(err error) *ErrorChain {
	chain := &ErrorChain{
		Chain:  make([]*Error, 0),
		Causes: make(map[string]int),
	}

	if e, ok := err.(*Error); ok {
		chain.Root = e
		chain.Chain = append(chain.Chain, e)

		// Analyze error causes
		curr := e.Err
		for curr != nil {
			if wrapped, ok := curr.(*Error); ok {
				chain.Chain = append(chain.Chain, wrapped)
				chain.Causes[wrapped.Category.String()]++
				curr = wrapped.Err
			} else {
				chain.Causes["unknown"]++
				break
			}
		}
	}

	return chain
}

func (ec *ErrorChain) Summary() string {
	var summary strings.Builder
	summary.WriteString("Error Chain Summary:\n")

	for i, err := range ec.Chain {
		summary.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, err.Category.String(), err.Message))
	}

	summary.WriteString("\nCause Distribution:\n")
	for cause, count := range ec.Causes {
		summary.WriteString(fmt.Sprintf("%s: %d\n", cause, count))
	}

	return summary.String()
}
