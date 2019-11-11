package entity

import "bufio"

type DownstreamMessage interface {
	Write(writer *bufio.Writer)
}

type UpstreamMessage interface {
}
