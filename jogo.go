// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"os"
	"sync"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo   	rune
	cor       	Cor
	corFundo  	Cor
	tangivel  	bool // Indica se o elemento bloqueia passagem
	interagivel bool
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa            [][]Elemento // grade 2D representando o mapa
	PosX, PosY      int          // posição atual do personagem
	UltimoVisitado  Elemento     // elemento que estava na posição do personagem antes de mover
	StatusMsg       string       // mensagem para a barra de status
	TemChave 	    bool		 // (adicionado) flag que indica se o personagem possui a chave no inventario
	TemArma 		bool		 // (adicionado) flag que indica se o personagem possui uma arma no inventario
	Inimigos        []inimigo	 // (adicionado) colecao para armazenar os inimigos ativos
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

type inimigo struct {
	X, Y  int
	Ativo bool
}
var mapaLocks [][]sync.Mutex // thread por celula do mapa

func inicializarMutexes(jogo *Jogo) {
	mapaLocks = make([][]sync.Mutex, len(jogo.Mapa))
	for y := range jogo.Mapa {
		mapaLocks[y] = make([]sync.Mutex, len(jogo.Mapa[y]))
	}
}

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo{UltimoVisitado: Vazio}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

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
				jogo.Inimigos = append(jogo.Inimigos, inimigo{X: x, Y: y, Ativo: true})
			case Vegetacao.simbolo:
				e = Vegetacao
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y // registra a posição inicial do personagem 
			case Bau.simbolo:
				e = Bau;
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	inicializarMutexes(jogo)

	return nil
}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy

	elemento := jogo.Mapa[y][x]         // personagem
	jogo.Mapa[y][x] = jogo.UltimoVisitado // limpa posição antiga
	jogo.UltimoVisitado = jogo.Mapa[ny][nx] // guarda o que havia na nova posição
	jogo.Mapa[ny][nx] = elemento          // move personagem para nova posição
}


