package main

import (
	"sort"
    "sync"
)

func girvanNewman(g *Grafo, concorrente bool) [][]int {
    var comunidades [][]int

    for {
        var intermediacao map[[2]int]float64
        if concorrente {
            intermediacao = calcularIntermediacaoArestasConcorrente(g)
        } else {
            intermediacao = calcularIntermediacaoArestas(g)
        }

        if len(intermediacao) == 0 {
            break
        }

        maiorAresta := encontrarMaiorAresta(intermediacao)
        g.RemoverAresta(maiorAresta[0], maiorAresta[1])
        comunidades = detectarComunidades(g)

        if len(comunidades) > 1 {
            break
        }
    }

    return ordenarComunidades(comunidades)
}

func calcularIntermediacaoArestas(g *Grafo) map[[2]int]float64 {
    intermediacao := make(map[[2]int]float64)
    for no := range g.nos {
        intermediacaoPorOrigem(g, no, intermediacao)
    }
    return intermediacao
}

func calcularIntermediacaoArestasConcorrente(g *Grafo) map[[2]int]float64 {
    intermediacao := make(map[[2]int]float64)
    var mu sync.Mutex

    wg := sync.WaitGroup{}
    for no := range g.nos {
        wg.Add(1)
        go func(no int) {
            defer wg.Done()
            intermediacaoLocal := make(map[[2]int]float64)
            intermediacaoPorOrigem(g, no, intermediacaoLocal)
            mu.Lock()
            for aresta, valor := range intermediacaoLocal {
                intermediacao[aresta] += valor
            }
            mu.Unlock()
        }(no)
    }

    wg.Wait()
    return intermediacao
}

func intermediacaoPorOrigem(g *Grafo, origem int, intermediacao map[[2]int]float64) {
    caminhos := make(map[int]int)
    distancias := make(map[int]int)
    predecessores := make(map[int][]int)
    pilha := []int{}
    fila := []int{origem}

    for no := range g.nos {
        distancias[no] = -1
    }
    distancias[origem] = 0
    caminhos[origem] = 1

    for len(fila) > 0 {
        atual := fila[0]
        fila = fila[1:]
        pilha = append(pilha, atual)

        for _, vizinho := range g.nos[atual] {
            if distancias[vizinho] == -1 {
                fila = append(fila, vizinho)
                distancias[vizinho] = distancias[atual] + 1
            }
            if distancias[vizinho] == distancias[atual]+1 {
                caminhos[vizinho] += caminhos[atual]
                predecessores[vizinho] = append(predecessores[vizinho], atual)
            }
        }
    }

    dependencias := make(map[int]float64)
    for len(pilha) > 0 {
        no := pilha[len(pilha)-1]
        pilha = pilha[:len(pilha)-1]

        for _, predecessor := range predecessores[no] {
            aresta := criarAresta(predecessor, no)
            peso := float64(caminhos[predecessor]) / float64(caminhos[no])
            dependencias[predecessor] += peso * (1 + dependencias[no])
            intermediacao[aresta] += dependencias[no]
        }

        if no != origem {
            dependencias[no] += 1.0
        }
    }
}

func criarAresta(u, v int) [2]int {
    if u < v {
        return [2]int{u, v}
    }
    return [2]int{v, u}
}

func encontrarMaiorAresta(intermediacao map[[2]int]float64) [2]int {
    var maiorAresta [2]int
    maiorValor := -1.0
    for aresta, valor := range intermediacao {
        if valor > maiorValor {
            maiorAresta = aresta
            maiorValor = valor
        }
    }
    return maiorAresta
}

func detectarComunidades(g *Grafo) [][]int {
    visitados := make(map[int]bool)
    var comunidades [][]int

    for no := range g.nos {
        if !visitados[no] {
            comunidade := []int{}
            dfs(g, no, &comunidade, visitados)
            if len(comunidade) > 0 {
                comunidades = append(comunidades, comunidade)
            }
        }
    }
    return comunidades
}

func dfs(g *Grafo, no int, comunidade *[]int, visitados map[int]bool) {
    visitados[no] = true
    *comunidade = append(*comunidade, no)
    for _, vizinho := range g.nos[no] {
        if !visitados[vizinho] {
            dfs(g, vizinho, comunidade, visitados)
        }
    }
}

func ordenarComunidades(comunidades [][]int) [][]int {
    for _, comunidade := range comunidades {
        sort.Ints(comunidade)
    }
    sort.Slice(comunidades, func(i, j int) bool {
        return comunidades[i][0] < comunidades[j][0]
    })
    return comunidades
}
