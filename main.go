// main.go - Loop principal do jogo
package main

import (
	"os"
)

func main() {
	interfaceIniciar()
	defer interfaceFinalizar()

	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	interfaceDesenharJogo(&jogo)
	inicializarInimigos(&jogo)
	inicializarPortais(&jogo)
	inicializarArmadilhas(&jogo)

	for {
        if jogo.Perdeu() {
            jogo.StatusMsg = "Fim do jogo! Pressione R para reiniciar ou ESC para sair."
            interfaceDesenharJogo(&jogo)
            for {
                evento := interfaceLerEventoTeclado()
                if evento.Tipo == "sair" {
                    return
                }
                if evento.Tipo == "mover" && (evento.Tecla == 'r' || evento.Tecla == 'R') {
					for i := range jogo.Inimigos {
						close(jogo.Inimigos[i].canalInimigos)
						close(jogo.Inimigos[i].canalMapa)
					}
					jogo = jogoNovo()
					if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
						panic(err)
					}
					interfaceDesenharJogo(&jogo)
					inicializarInimigos(&jogo)
					inicializarPortais(&jogo)
					inicializarArmadilhas(&jogo)
					break
				}
            }
            continue
        }

        evento := interfaceLerEventoTeclado()
        if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
            break
        }
        interfaceDesenharJogo(&jogo)
    }
}