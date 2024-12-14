package cmd

import (
    "time"
    "github.com/briandowns/spinner"
)

func newSpinner(message string) *spinner.Spinner {
    s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
    s.Suffix = " " + message
    return s
}

func newProgressBar(total int, message string) *spinner.Spinner {
    s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
    s.Suffix = " " + message
    return s
}
