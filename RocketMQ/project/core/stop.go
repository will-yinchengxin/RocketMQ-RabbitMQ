package core

import "rock/rocket"

func Stop() {
	rocket.CloseProducer()
}
