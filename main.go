// main.go - Loop principal do jogo
package main

import (
	"os"
	"time"
)

func main() {
	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// (adicionado) Thread de atualizacao periodica 
	// Atualiza a interface periodicamente para exibir movimento dos inimigos
	go func() {
		for {
			time.Sleep(100 * time.Millisecond) // Taxa de atualização da interface
			interfaceDesenharJogo(&jogo)
		}
	}()

	// Thread que inicializa o comportamento dos inimigos 
	go func() {
		for {
			iniciarMovimentoInimigos(&jogo)
			time.Sleep(1 * time.Second)
		}
	}()
	
	// Loop principal de entrada
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			break
		}
		interfaceDesenharJogo(&jogo)
	}
}