package option

import "time"

type ConnectorOption struct {
	Type string
	NoDelay bool
	Deadline time.Duration
}
