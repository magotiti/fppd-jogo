package main

import (
	"time"
)

type armadilha struct {
	X, Y       int
	Ativada    bool
	canalMapa  chan Mensagem
}

func NovaArmadilha(x, y int) armadilha {
	return armadilha{
		X:         x,
		Y:         y,
		Ativada:   false,
		canalMapa: make(chan Mensagem, 1),
	}
}

func inicializarArmadilhas(jogo *Jogo) {
    for i := range jogo.Armadilhas {
        armadilha := &jogo.Armadilhas[i]
        go rotinaArmadilha(jogo, armadilha)
        }
    }

func rotinaArmadilha(jogo *Jogo, a *armadilha) {
	for {
		select {
		case msg := <-a.canalMapa:
			if msg.Tipo == "Ativar!" {
				mapaLeituraLock.Lock()
				a.Ativada = true
				jogo.Mapa[a.Y][a.X] = Armadilha
				mapaLeituraLock.Unlock()

				interfaceDesenharJogo(jogo)

				go func() {
					time.Sleep(5 * time.Second)

					mapaLeituraLock.Lock()
					a.Ativada = false
					jogo.Mapa[a.Y][a.X] = Vazio
					mapaLeituraLock.Unlock()

					interfaceDesenharJogo(jogo)
				}()
			}
		default:
			if abs(jogo.PosX-a.X) <= 1 && abs(jogo.PosY-a.Y) <= 1 {
				mapaLeituraLock.Lock()
				ativada := a.Ativada
				mapaLeituraLock.Unlock()

				if !ativada {
					select {
					case a.canalMapa <- Mensagem{Tipo: "Ativar"}:
					default:
					}
					disparaAlarme(jogo)
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
