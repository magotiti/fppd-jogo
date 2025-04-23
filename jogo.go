// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"os"
	"sync"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo     rune
	cor         Cor
	corFundo    Cor
	tangivel    bool // Indica se o elemento bloqueia passagem
	interagivel bool
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa           [][]Elemento // grade 2D representando o mapa
	PosX, PosY     int          // posição atual do personagem
	VidaJogador	   int			// vida atual do personagem
	UltimoVisitado Elemento     // elemento que estava na posição do personagem antes de mover
	StatusMsg      [3]string       // mensagem para a barra de status
	TemChave       bool         // flag que indica se o personagem possui a chave no inventário
	TemArma        bool         // flag que indica se o personagem possui uma arma no inventário
	Inimigos       []inimigo    // coleção para armazenar os inimigos ativos
	BausAbertos int
    ArmaGarantida bool
	Pontuacao int
}

// Elementos visuais do jogo
var (
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true, false}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true, false}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false, false}
	Bau        = Elemento{'⌂', CorAmarelo, CorPadrao, false, true}
	Porta      = Elemento{'╬', CorAzulClaro, CorFundoParede, true, true}
	Chave      = Elemento{'✒', CorAmarelo, CorPadrao, false, true}
	Arma       = Elemento{'⚔', CorCinzaEscuro, CorPadrao, false, true}
)

var mapaLeituraLock sync.Mutex

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	return Jogo{UltimoVisitado: Vazio}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()
	mapaLeituraLock.Lock()
	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Parede.simbolo:
				e = Parede
			case Inimigo.simbolo:
				e = Vazio
				jogo.Inimigos = append(jogo.Inimigos, inimigo{
					X			  : x,
					Y			  : y,
					Ativo 	      : true,
					Vida          : 99,
					canalMapa     : make(chan Mensagem, 4),
					canalInimigos : make(chan Mensagem, 4),
				})
			case Vegetacao.simbolo:
				e = Vegetacao
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY, jogo.VidaJogador = x, y, 999
			case Bau.simbolo:
				e = Bau
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	mapaLeituraLock.Unlock()
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}
	if jogo.Mapa[y][x].tangivel {
		return false
	}
	return true
}

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy
	elemento := jogo.Mapa[y][x]
	jogo.Mapa[y][x] = jogo.UltimoVisitado
	jogo.UltimoVisitado = jogo.Mapa[ny][nx]
	jogo.Mapa[ny][nx] = elemento
}

//////////////////////////////////////////////////////////////////////
//  Funcao     : disparaAlarme
//  Descricao  : Registra o evento e dispara uma mensagem para o canal de inimigos
// 	Criado     : Thiago Cardoso							  [16/04/2025]
//  Modificado : 				
//////////////////////////////////////////////////////////////////////
func disparaAlarme(jogo *Jogo) {
	msg := Mensagem{
		Tipo    : "Alarme!",
		OrigemX : jogo.PosX,
		OrigemY : jogo.PosY,
	}

	for i := range jogo.Inimigos {
		select {
		case jogo.Inimigos[i].canalMapa <- msg:
		default:
		}
	}
}

func adicionarMensagem(jogo *Jogo, msg string) {
    jogo.StatusMsg[0] = jogo.StatusMsg[1]
    jogo.StatusMsg[1] = jogo.StatusMsg[2]
    jogo.StatusMsg[2] = msg
}
