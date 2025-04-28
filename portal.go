package main

import (
	"math/rand"
	"time"
)

type portal struct {
	X, Y        int
	Ativo       bool
	canalMapa   chan Mensagem
	tempoAtivo  time.Duration
	DestX, DestY int
}

func NovoPortal(x, y int, jogo *Jogo) portal {
	min := 15
	max := 30
	destX, destY := encontrarDestinoValido(jogo)
	return portal{
		X:          x,
		Y:          y,
		Ativo:      true,
		canalMapa:  make(chan Mensagem, 1),
		tempoAtivo: time.Duration(rand.Intn(max-min+1)+min) * time.Second,
		DestX:      destX,
		DestY:      destY,
	}
}

func inicializarPortais(jogo *Jogo) {
for i := range jogo.Portais {
    portal := &jogo.Portais[i]
    go rotinaPortal(jogo, portal)
    }
}

func rotinaPortal(jogo *Jogo, p *portal) {
	// elemento visual do portal
    elemDestino := Elemento{'╬', CorRoxo, CorPadrao, false, true}
	for {
		select {
		case msg := <-p.canalMapa:
            if (msg.Tipo == "Teleporte!") { 
                mapaLeituraLock.Lock()
                jogo.PosX = p.DestX
                jogo.PosY = p.DestY
                p.Ativo   = true
                mapaLeituraLock.Unlock()
            }
		// após p.tempoAtivo, alterna entre remover e redesenhar o portal
		case <-time.After(p.tempoAtivo):
			mapaLeituraLock.Lock()
			if p.Ativo {
				// desativa: limpa as duas posições
				jogo.Mapa[p.Y][p.X] = Vazio
				jogo.Mapa[p.DestY][p.DestX] = Vazio
				p.Ativo = false
			} else {
				// reativa: redesenha nas mesmas coordenadas
				jogo.Mapa[p.Y][p.X] = Portal
				jogo.Mapa[p.DestY][p.DestX] = elemDestino
				p.Ativo = true
			}

			mapaLeituraLock.Unlock()
			interfaceDesenharJogo(jogo)
		}
	}
}

func encontrarDestinoValido(jogo *Jogo) (int, int) {
	for {
		y := rand.Intn(len(jogo.Mapa))
		x := rand.Intn(len(jogo.Mapa[y]))

		// Verifica se a posição é válida
		if !jogo.Mapa[y][x].tangivel {
			return x, y
		}
	}
}
