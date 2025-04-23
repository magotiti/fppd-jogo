// personagem.go - Funções para movimentação e ações do personagem
package main

import (
	"math/rand"
)

const c_RangeInteracao = 2

type elementoPosicao struct {
	Elemento Elemento
	X, Y     int
}

type itemBau int
const (
	nadaEncontrado itemBau = iota
	encontrouArma
	encontrouChave
	encontrouArmadilha
)


// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(tecla rune, jogo *Jogo) {
    dx, dy := 0, 0
    switch tecla {
    case 'w': dy = -1 // Move para cima
    case 'a': dx = -1 // Move para a esquerda
    case 's': dy = 1  // Move para baixo
    case 'd': dx = 1  // Move para a direita
    }

    nx, ny := jogo.PosX+dx, jogo.PosY+dy
    // Verifica se o movimento é permitido
    if jogoPodeMoverPara(jogo, nx, ny) {
        // Verifica se há uma armadilha na célula de destino
        if jogo.Mapa[ny][nx].simbolo == 'A' {
            jogo.VidaJogador -= 100 // Reduz a vida do jogador em 100
            if jogo.VidaJogador < 0 {
                jogo.VidaJogador = 0
            }
            adicionarMensagem(jogo, "Você pisou em uma armadilha e perdeu 100 de vida!")
        }

        // Realiza a movimentação
        jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
        jogo.PosX, jogo.PosY = nx, ny
    }
}

// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, apenas exibe uma mensagem de status
// Você pode expandir essa função para incluir lógica de interação com objetos
func personagemInteragir(jogo *Jogo) {
	alvo := buscaElementoMaisProximo(jogo)
	if alvo == nil {
		adicionarMensagem(jogo, "Nao ha nada para interagir por perto.")
		return
	}
	switch alvo.Elemento {
	case Inimigo:
		if jogo.TemArma {
			for i := range jogo.Inimigos {
				if jogo.Inimigos[i].X == alvo.X && jogo.Inimigos[i].Y == alvo.Y && jogo.Inimigos[i].Ativo {
					// Dano de 33 na vida
					jogo.Inimigos[i].Vida -= 33
					if jogo.Inimigos[i].Vida <= 0 {
						adicionarMensagem(jogo, "Você atacou e eliminou o inimigo!")
						jogo.Pontuacao += 100
						go func(enemy *inimigo) {
							enemy.canalMapa <- Mensagem{Tipo: "Morreu!"}
						}(&jogo.Inimigos[i])
					}
					break
				}
			}
		} else {
			adicionarMensagem(jogo, "Você precisa de uma arma para atacar o inimigo!")
		}
	case Bau:
		abrirBau(jogo, alvo.X, alvo.Y)
	case Porta:
		if jogo.TemChave {
			adicionarMensagem(jogo, "Você usou uma chave para abrir a porta!")
			jogo.TemChave = false
			jogo.Mapa[alvo.Y][alvo.X] = Vazio
		} else {
			adicionarMensagem(jogo, "Você precisa de uma chave para abrir esta porta.")
		}
	default:
		// Verifica se é um portal ativo
        for i := range portais {
            if portais[i].X == alvo.X && portais[i].Y == alvo.Y && portais[i].Ativo {
                // Teleporta o personagem
                jogo.PosX = portais[i].DestX
                jogo.PosY = portais[i].DestY
                adicionarMensagem(jogo, "Você usou o portal!")
                return
            }
        }
        adicionarMensagem(jogo, "Voce nao pode interagir com esse elemento.")
    }
}


// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
    if jogo.VidaJogador <= 0 {
        return true // Não faz nada se estiver morto
    }
    switch ev.Tipo {
    case "sair":
        return false
    case "interagir":
        personagemInteragir(jogo)
    case "mover":
        personagemMover(ev.Tecla, jogo)
    }
    return true
}

//////////////////////////////////////////////////////////////////////
//  Funcao    : buscaElementos
//  Descricao : Retorna todos elementos interagiveis dentro do range
//			    definido por c_RangeInteracao
// 	Criado     : Thiago Cardoso							  [13/04/2025]
//  Modificado : 				
//////////////////////////////////////////////////////////////////////
func buscaElementos(jogo *Jogo) []elementoPosicao {
	var encontrados []elementoPosicao

	for dy := -c_RangeInteracao; dy <= c_RangeInteracao; dy++ {
		for dx := -c_RangeInteracao; dx <= c_RangeInteracao; dx++ {
			if dx != 0 || dy != 0 {
				x := jogo.PosX + dx
				y := jogo.PosY + dy

				if y >= 0 && y < len(jogo.Mapa) && x >= 0 && x < len(jogo.Mapa[y]) {
					e := jogo.Mapa[y][x]
					if e != Vazio && e.interagivel {
						encontrados = append(encontrados, elementoPosicao{
							Elemento: e,
							X:        x,
							Y:        y,
						})
					}
				}
			} else { 
			continue
			}
		}
	}

	return encontrados
}

//////////////////////////////////////////////////////////////////////
//  Funcao     : buscaElementoMaisProximo
//  Descricao  : Itera sobre a lista de elementos no range e retorna
//				 o mais proximo
// 	Criado     : Thiago Cardoso							  [13/04/2025]
//  Modificado : 				
//////////////////////////////////////////////////////////////////////
func buscaElementoMaisProximo (jogo *Jogo) *elementoPosicao {
	var elementos = buscaElementos(jogo)
	if len(elementos) != 0 {
		maisPerto := elementos[0]
		menorDist := (abs(jogo.PosX - maisPerto.X) + abs(jogo.PosY - maisPerto.Y))

		for i := 1; i < len(elementos); i++ {
			e := elementos [i]
			dist := (abs(jogo.PosX - e.X) + abs(jogo.PosY - e.Y))
			if dist < menorDist {
				maisPerto = e
				menorDist = dist
			}
		}
		return &maisPerto	
	}
	return nil	
}

//////////////////////////////////////////////////////////////////////
//  Funcao     : abrirBau
//  Descricao  : Define a interacao com os baus
// 	Criado     : Thiago Cardoso							  [13/04/2025]
//  Modificado : Thiago Cardoso							  [16/04/2025]
//			   - Adiciona disparaAlarme;				
//////////////////////////////////////////////////////////////////////
func abrirBau(jogo *Jogo, x, y int) itemBau {
	encontrou := nadaEncontrado
    jogo.BausAbertos++
    totalBaus := 4 // tem que ajustar se mudar a quantidade de baús no mapa

	temArma 	 := rand.Intn(100) < 33 // 33% chance de achar uma arma
	temChave 	 := rand.Intn(100) < 33 // 33% chance de achar uma chave
	temArmadilha := rand.Intn(100) < 75 // 75% chance de achar uma armadilha (alarme)

	if jogo.BausAbertos == totalBaus && !jogo.ArmaGarantida {
        temArma = true
    }

	if temArma && !jogo.TemArma {
		adicionarMensagem(jogo, "Parabens! Voce encontrou uma arma!")
		jogo.TemArma = true
		jogo.ArmaGarantida = true
		encontrou = encontrouArma
	}
	if temChave && !jogo.TemChave {
        adicionarMensagem(jogo, "Parabens! Voce encontrou uma chave!")
        jogo.TemChave = true
        if encontrou != encontrouArma {
            encontrou = encontrouChave
        }
    }
	if temArmadilha {
		adicionarMensagem(jogo, "Essa nao! Voce disparou um alarme!")
		disparaAlarme(jogo)
	}

	if !temArma && !temChave {
		adicionarMensagem(jogo, "O bau estava vazio...")
	}

	jogo.Mapa[y][x] = Vazio
	return encontrou
}
