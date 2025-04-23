package main

import (
	"time"
)

type armadilha struct {
    X, Y      int
    Ativada   bool
    canalMapa chan Mensagem
}

func NovaArmadilha(x, y int) armadilha {
    return armadilha{
        X:        x,
        Y:        y,
        Ativada:  false,
        canalMapa: make(chan Mensagem, 1),
    }
}

func rotinaArmadilha(jogo *Jogo, a *armadilha) {
    for {
        select {
        case msg := <-a.canalMapa:
            if msg.Tipo == "Ativar" && !a.Ativada {
                a.Ativada = true
                jogo.Mapa[a.Y][a.X] = Elemento{'A', CorVermelho, CorPadrao, false, true}
                interfaceDesenharJogo(jogo)
                go func() {
                    time.Sleep(5 * time.Second)
                    a.Ativada = false
                    jogo.Mapa[a.Y][a.X] = Vazio
                    interfaceDesenharJogo(jogo)
                }()
            }
        default:
            if abs(jogo.PosX-a.X) <= 1 && abs(jogo.PosY-a.Y) <= 1 && !a.Ativada {
                a.canalMapa <- Mensagem{Tipo: "Ativar"}
                disparaAlarme(jogo)
            }
        }
        time.Sleep(100 * time.Millisecond)
    }
}