// main.go - Loop principal do jogo
package main

import (
	"os"
)

var portais []portal

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

	portais = []portal{
		NovoPortal(12, 10, 65, 2), // portal leva de (12, 10) para (65, 2)
	}
	for i := range portais {
		if portais[i].Ativo {
			// desenha o portal na posição de entrada
			jogo.Mapa[portais[i].Y][portais[i].X] = Elemento{'O', CorAzulClaro, CorPadrao, false, true}
			// desenha o portal na posição de saída
			jogo.Mapa[portais[i].DestY][portais[i].DestX] = Elemento{'O', CorAzulClaro, CorPadrao, false, true}
		}
		go rotinaPortal(&jogo, &portais[i])
	}

    armadilhas := []armadilha{
        NovaArmadilha(5, 5),
        NovaArmadilha(12, 18),
    }
    for i := range armadilhas {
        go rotinaArmadilha(&jogo, &armadilhas[i])
    }

	interfaceDesenharJogo(&jogo)
	
	inicializarInimigos(&jogo)

	for {
        if jogo.VidaJogador <= 0 {
            // Exibe mensagem final
            jogo.StatusMsg = "Fim do jogo! Pressione R para reiniciar ou ESC para sair."
            interfaceDesenharJogo(&jogo)
            // Espera até o usuário pressionar R ou ESC
            for {
                evento := interfaceLerEventoTeclado()
                if evento.Tipo == "sair" {
                    return
                }
                if evento.Tipo == "mover" && (evento.Tecla == 'r' || evento.Tecla == 'R') {
					// Reinicia o jogo
					jogo = jogoNovo()
					jogo.Inimigos = nil // Limpa inimigos
					jogo.Mapa = nil     // Limpa o mapa!
					if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
						panic(err)
					}
					interfaceDesenharJogo(&jogo)
					inicializarInimigos(&jogo)
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