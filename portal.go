package main

import (
	"time"
)

type portal struct {
    X, Y       int
    Ativo      bool
    canalMapa  chan Mensagem
    tempoAtivo time.Duration
    DestX, DestY int // posição de destino do portal
}

func NovoPortal(x, y, destX, destY int) portal {
    return portal{
        X:          x,
        Y:          y,
        Ativo:      true,
        canalMapa:  make(chan Mensagem, 1),
        tempoAtivo: 10 * time.Second,
        DestX:      destX,
        DestY:      destY,
    }
}

func rotinaPortal(jogo *Jogo, p *portal) {
    for {
        select {
        case msg := <-p.canalMapa:
            if msg.Tipo == "Ativar" {
                p.Ativo = true
                jogo.Mapa[p.Y][p.X] = Elemento{'O', CorAzulClaro, CorPadrao, false, true}
                jogo.Mapa[p.DestY][p.DestX] = Elemento{'O', CorAzulClaro, CorPadrao, false, true}
                interfaceDesenharJogo(jogo)
            }
        case <-time.After(p.tempoAtivo):
            if p.Ativo {
                p.Ativo = false
                jogo.Mapa[p.Y][p.X] = Vazio
                jogo.Mapa[p.DestY][p.DestX] = Vazio
                interfaceDesenharJogo(jogo)
            }
        }
    }
}