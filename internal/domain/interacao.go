package domain

import "time"

type Interacao struct {
	De, Para string
	Jogo     Jogo
	Curtiu   bool
	Em       time.Time
}
