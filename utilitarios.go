package main

// aponta direcao da coordenada
func direcao(orig, dest int) int {
	if orig < dest {
		return 1
	} else if orig > dest {
		return -1
	}
	return 0
}

// retorna operacao absolute()
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}