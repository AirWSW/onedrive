package graphapi

import "time"

type ODataError struct {
	Code                     *string     `json:"code"`
	Message                  *string     `json:"message"`
	RequestID                *string     `json:"request-id"`
	Date                     *time.Time  `json:"date"`
	MicrosoftGraphInnerError *ODataError `json:"innerError,omitempty"`
}

func (e *ODataError) OnError() {

}
