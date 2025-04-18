/////////////////////////////////////////////////////////////////////////////////
//  ██╗███╗   ██╗██╗███╗   ███╗██╗ ██████╗  ██████╗      ██████╗  ██████╗
//  ██║████╗  ██║██║████╗ ████║██║██╔════╝ ██╔═══██╗    ██╔════╝ ██╔═══██╗
//  ██║██╔██╗ ██║██║██╔████╔██║██║██║  ███╗██║   ██║    ██║  ███╗██║   ██║
//  ██║██║╚██╗██║██║██║╚██╔╝██║██║██║   ██║██║   ██║    ██║   ██║██║   ██║
//  ██║██║ ╚████║██║██║ ╚═╝ ██║██║╚██████╔╝╚██████╔╝ ██ ╚██████╔╝╚██████╔╝
//  ╚═╝╚═╝  ╚═══╝╚═╝╚═╝     ╚═╝╚═╝ ╚═════╝  ╚═════╝  ╚═╝ ╚═════╝  ╚═════╝
/////////////////////////////////////////////////////////////////////////////////
//  Classe 	  : inimigo.go
//  Descricao : Define o comportamento do elemento inimigo
//  Criado	  : Thiago Cardoso									  [14/04/2025]
//
//  Checklist : [X] Novos elementos concorrentes (1/3)
//				[X] Comunicação entre elementos por canais (inimigoRotina)
// 				[X] Escuta concorrente de multiplos canais (inimigoRotina)
//				[X] Exclusao mutua
/////////////////////////////////////////////////////////////////////////////////

package main

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

// estrutura do inimigo
type inimigo struct {
	X, Y           int
	Ativo          bool
	Vida		   int
	canalMapa 	   chan Mensagem
	canalInimigos  chan Mensagem
	cor		       termbox.Attribute
}
// estrutura da mensagem
type Mensagem struct {
	Tipo    string
	OrigemX int
	OrigemY int
}
// variaveis globais e constantes
var   RangePersonagem 			   = 10
const c_danoAtaque				   = 33
var   sleepMovimento time.Duration = 500 //(ms)
const c_RangeAlerta	  			   = 15
// canais de comunicacao interna
var   mapaCanais 				   = make(map[int]chan Mensagem) 

/////////////////////////////////////////////////////////////////////////////////
//  Funcao 	   : NovoInimigo
//  Descricao  : Construtor para criar uma nova instância de inimigo
//	Criado 	   : Thiago Cardoso	 								  [14/04/2025]
//	Modificado : 
/////////////////////////////////////////////////////////////////////////////////
func NovoInimigo(x, y int) inimigo {	
	return inimigo{
		X            : x,
		Y            : y,
		Ativo        : true,
		Vida         : 99, // 3 ataques para morrer
		canalMapa    : make(chan Mensagem, 10), // interacoes com o resto do jogo
		canalInimigos: make(chan Mensagem, 10), // interacoes com outros inimigos
		cor 		 : CorVermelho,
	}
}

///////////////////////////////////////////////////////////////////////////////	/
//  Funcao 	   : inicializarInimigos
//  Descricao  : Percorre a lista de inimigos e dispara suas rotinas
//	Criado 	   : Thiago Cardoso	 								  [14/04/2025]
//	Modificado : 
/////////////////////////////////////////////////////////////////////////////////
func inicializarInimigos(jogo *Jogo) {
	for i := range jogo.Inimigos {
		enemy := &jogo.Inimigos[i]
		mapaCanais[i] = enemy.canalInimigos
		go rotinaInimigo(jogo, enemy, i)
	}
}

/////////////////////////////////////////////////////////////////////////////////
//  Funcao 	   : rotinaInimigo
//  Descricao  : Define e propaga o comportamento do inimigo de acordo com seu
//				 estado atual; consome sinais vindos dos canais de comunicacao
//	Criado 	   : Thiago Cardoso	 								  [16/04/2025]
//	Modificado : Thiago Cardoso									  [18/04/2025]
//			   - Implementa novo canal que recebe mensagens externas;
//			   - Adiciona suporte para novos tipos de mensagem; 
/////////////////////////////////////////////////////////////////////////////////
func rotinaInimigo(jogo *Jogo, enemy *inimigo, id int) {
	for {
		movido := false
		morreu := false
		// nao interage com inimigos desativados (mortos)
		if !enemy.Ativo {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		select {
		// CANAL EXTERNO:
		// ocorrencia de sinal quando uma armadilha eh encontrada (mensagem ja propagada)
		case msg := <-enemy.canalMapa:
			switch msg.Tipo {
			case "Alarme!" :
				RangePersonagem = 25  // os inimigos captam o personagem muito mais longe,
				sleepMovimento  = 250 // e se movem muito mais rapido tambem
				enemy.cor       = CorRoxo
				// cria uma goroutine para restaurar os valores apos 15 segundos
				go func() {
					time.Sleep(15 * time.Second)
					RangePersonagem = 10   
					sleepMovimento  = 500
					enemy.cor       = CorVermelho
				}()
				// movimento aleatorio padrao
				dx := rand.Intn(3) - 1
				dy := rand.Intn(3) - 1
				movido =  enemy.inimigoMover(jogo, dx, dy)
			case "Morreu!" :
				morreu = enemy.morrerERespawnar(jogo)
			case "Spawnou!":
				// ativa o inimigo novamente
				enemy.Ativo = true
				enemy.Vida = 99
				// ressurge enfurecido, se autoenvia uma mensagem
				msg := Mensagem{ Tipo: "Alarme!", OrigemX: enemy.X, OrigemY: enemy.Y} 
				enemy.canalMapa <- msg
			} 

		// CANAL INTERNO:
		// ocorrencia de sinal quando o personagem eh avistado por um inimgo
		case msg := <-enemy.canalInimigos:
			if msg.Tipo == "Encontrado!" &&
				// outros inimigos estao no range do alerta 
				abs(enemy.X - msg.OrigemX) <= c_RangeAlerta &&
				abs(enemy.Y - msg.OrigemY) <= c_RangeAlerta {
				// -> direciona todos para as coordenadas do personagem
				dx := direcao(enemy.X, msg.OrigemX)
				dy := direcao(enemy.Y, msg.OrigemY)
				movido = enemy.inimigoMover(jogo, dx, dy)
			}
		default:
			if enemy.personagemAVista(jogo) {
				msg := Mensagem{Tipo: "Encontrado!", OrigemX: enemy.X, OrigemY: enemy.Y}
				// <- propaga um sinal no canal de todos os outros 
				for outroId, outroCanal := range mapaCanais {
					if outroId != id {
						select {
						case outroCanal <- msg:
						default:
						}
					}
				}
				dx := direcao(enemy.X, jogo.PosX)
				dy := direcao(enemy.Y, jogo.PosY)
				movido =  enemy.inimigoMover(jogo, dx, dy)
				// atacar se o personagem estiver proximo
				enemy.atacarSeProximo(jogo)
			} else {
				// movimento aleatorio padrao 
				dx := rand.Intn(3) - 1
				dy := rand.Intn(3) - 1
				movido =  enemy.inimigoMover(jogo, dx, dy)
			}
		}
		if movido || morreu {
			mapaLeituraLock.Lock()
			interfaceDesenharJogo(jogo)
			mapaLeituraLock.Unlock()
		}
		time.Sleep(sleepMovimento * time.Millisecond)
	}
}

/////////////////////////////////////////////////////////////////////////////////
//  Funcao 	   : inimigoMover
//  Descricao  : Realiza interacao de movimento do inimigo e trata as
//				 concorrencias com outros elementos semelhantes
//	Criado 	   : Thiago Cardoso	 								  [14/04/2025]
//	Modificado : 
/////////////////////////////////////////////////////////////////////////////////
func (enemy *inimigo) inimigoMover(jogo *Jogo, dx, dy int) bool {
	nx, ny := enemy.X + dx, enemy.Y + dy
	if ny < 0 || ny >= len(jogo.Mapa) || nx < 0 || nx >= len(jogo.Mapa[ny]) {
		return false
	}
	mapaLeituraLock.Lock()
	defer mapaLeituraLock.Unlock()
	original := jogo.Mapa[enemy.Y][enemy.X]
	destino  := jogo.Mapa[ny][nx]
	// inimigos nao atravessam objetos tangiveis e nem entram na vegetacao
	if jogoPodeMoverPara(jogo, nx, ny) && !destino.tangivel && destino != Vegetacao {
		if original == Inimigo {
			original = Vazio
		}
		jogo.Mapa[enemy.Y][enemy.X] = original
		jogo.Mapa[ny][nx] 			= Inimigo
		enemy.X = nx
		enemy.Y = ny
		return true
	}

	return false
}

/////////////////////////////////////////////////////////////////////////////////
//  Metodo 	   : personagemAVista
//  Descricao  : Verifica se o jogador esta no campo de visao dos inimigos
//	Criado 	   : Thiago Cardoso	 								  [14/04/2025]
//	Modificado : 
/////////////////////////////////////////////////////////////////////////////////
func (enemy *inimigo) personagemAVista(jogo *Jogo) bool {
	if abs(jogo.PosX - enemy.X) <= RangePersonagem &&
	   abs(jogo.PosY - enemy.Y) <= RangePersonagem {

		direcaoX := direcao(enemy.X, jogo.PosX)
		direcaoY := direcao(enemy.Y, jogo.PosY)
		x, y := enemy.X, enemy.Y 
		for x != jogo.PosX || y != jogo.PosY {
			if x != jogo.PosX {
				x += direcaoX
			}
			if y != jogo.PosY {
				y += direcaoY
			}
			if jogo.Mapa[y][x] == Parede {
				return false
			}
		}
		return true
	}
	return false
}

/////////////////////////////////////////////////////////////////////////////////
//  Metodo 	   : morrerERespawnar
//  Descricao  : Finaliza o inimigo e faz o respawn em um ponto fixo do mapa
//	Criado 	   : Thiago Cardoso	 								  [18/04/2025]
//	Modificado : 
/////////////////////////////////////////////////////////////////////////////////
func (enemy *inimigo) morrerERespawnar(jogo *Jogo) bool {
	mapaLeituraLock.Lock()
	defer mapaLeituraLock.Unlock()

	if enemy.X >= 0 && enemy.Y >= 0 &&
		enemy.Y < len(jogo.Mapa) && enemy.X < len(jogo.Mapa[enemy.Y]) {
		if jogo.Mapa[enemy.Y][enemy.X] == Inimigo {
			jogo.Mapa[enemy.Y][enemy.X] = Vazio
		}
	}

	enemy.Ativo = false
	enemy.X, enemy.Y = -1, -1
	enemyCopia := enemy
	jogoCopia := jogo

	// dispara uma rotina de ressurgimento
	go func() {
		time.Sleep(20 * time.Second) // inimigo ressurge apos 20 sec
		mapaLeituraLock.Lock()
		defer mapaLeituraLock.Unlock()

		if enemyCopia != nil || jogoCopia != nil {
			enemyCopia.X = 5
			enemyCopia.Y = 5
			enemyCopia.Ativo = true
			enemyCopia.Vida = 99

			if enemyCopia.Y >= 0 && enemyCopia.Y < len(jogoCopia.Mapa) &&
				enemyCopia.X >= 0 && enemyCopia.X < len(jogoCopia.Mapa[enemyCopia.Y]) {
				jogoCopia.Mapa[enemyCopia.Y][enemyCopia.X] = Inimigo
			}
			select {
			// propaga para o canal para retomar a sua rotina
			case enemyCopia.canalMapa <- Mensagem{Tipo: "Spawnou!"}:
			default:
		}
	}
	}()
	return true
}

/////////////////////////////////////////////////////////////////////////////////
//  Metodo 	   : atacarSeProximo
//  Descricao  : Verifica se eh possivel atacar o jogador e realiza a acao
//	Criado 	   : Thiago Cardoso	 								  [18/04/2025]
//	Modificado : 
/////////////////////////////////////////////////////////////////////////////////
func (enemy *inimigo) atacarSeProximo(jogo *Jogo) {
	// verifica se está a 1 celula de distancia
	if abs(enemy.X-jogo.PosX) <= 1 && abs(enemy.Y-jogo.PosY) <= 1 {
		jogo.VidaJogador -= c_danoAtaque
		if jogo.VidaJogador < 0 {
			jogo.VidaJogador = 0
		}
	}
}
