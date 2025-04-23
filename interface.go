// interface.go - Interface gráfica do jogo usando termbox
// O código abaixo implementa a interface gráfica do jogo usando a biblioteca termbox-go.
// A biblioteca termbox-go é uma biblioteca de interface de terminal que permite desenhar
// elementos na tela, capturar eventos do teclado e gerenciar a aparência do terminal.

package main

import (
	"github.com/nsf/termbox-go"
)

type Cor = termbox.Attribute

const (
	CorPadrao     Cor = termbox.ColorDefault
	CorCinzaEscuro    = termbox.ColorDarkGray
	CorVermelho       = termbox.ColorRed
	CorVerde          = termbox.ColorGreen
	CorParede         = termbox.ColorBlack | termbox.AttrBold | termbox.AttrDim
	CorFundoParede    = termbox.ColorDarkGray
	CorTexto          = termbox.ColorDarkGray
	CorAmarelo 		  = termbox.ColorYellow
	CorAzulClaro      = termbox.ColorLightBlue
	CorRoxo		      = termbox.ColorLightMagenta
)

type EventoTeclado struct {
	Tipo  string
	Tecla rune
}

func interfaceIniciar() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
}

func interfaceFinalizar() {
	termbox.Close()
}

func interfaceLerEventoTeclado() EventoTeclado {
	ev := termbox.PollEvent()
	if ev.Type != termbox.EventKey {
		return EventoTeclado{}
	}
	if ev.Key == termbox.KeyEsc {
		return EventoTeclado{Tipo: "sair"}
	}
	if ev.Ch == 'e' {
		return EventoTeclado{Tipo: "interagir"}
	}
	return EventoTeclado{Tipo: "mover", Tecla: ev.Ch}
}

func interfaceDesenharJogo(jogo *Jogo) {
	interfaceLimparTela()

	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			interfaceDesenharElemento(x, y, elem)
		}
	}

	interfaceDesenharElemento(jogo.PosX, jogo.PosY, Personagem)

	interfaceDesenharBarraDeStatus(jogo)

	interfaceAtualizarTela()
}

func interfaceLimparTela() {
	termbox.Clear(CorPadrao, CorPadrao)
}

func interfaceAtualizarTela() {
	termbox.Flush()
}

func interfaceDesenharElemento(x, y int, elem Elemento) {
	termbox.SetCell(x, y, elem.simbolo, elem.cor, elem.corFundo)
}

func interfaceDesenharBarraDeStatus(jogo *Jogo) {
    for i, c := range jogo.StatusMsg {
        termbox.SetCell(i, len(jogo.Mapa)+1, c, CorTexto, CorPadrao)
    }

    vida := "Vida: [ "
    totalBlocos := 10
    vidaPorBloco := 999 / totalBlocos // 999 é a vida máxima inicial
    blocosCheios := jogo.VidaJogador / vidaPorBloco
    if blocosCheios > totalBlocos {
        blocosCheios = totalBlocos
    }
    for i := 0; i < totalBlocos; i++ {
        if i < blocosCheios {
            vida += "█"
        } else {
            vida += "░"
        }
        vida += " "
    }
    vida += "]"

    for i, c := range vida {
        termbox.SetCell(i, len(jogo.Mapa)+2, c, CorVerde, CorPadrao)
    }

    itens := "Itens: "
    if jogo.TemArma {
        itens += "🔫"
    }
    if jogo.TemChave {
        itens += "🔑"
    }
    if !jogo.TemArma && !jogo.TemChave {
        itens += "Nenhum"
    }

    for i, c := range itens {
        termbox.SetCell(i, len(jogo.Mapa)+3, c, CorTexto, CorPadrao)
    }

    msg := "Use WASD para mover e E para interagir. ESC para sair."
    for i, c := range msg {
        termbox.SetCell(i, len(jogo.Mapa)+5, c, CorTexto, CorPadrao)
    }
}