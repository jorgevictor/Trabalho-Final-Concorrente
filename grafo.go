package main

import (
    "bufio"
    "os"
    "strconv"
    "strings"
)

// Estrutura do Grafo
type Grafo struct {
    nos map[int][]int
}

// Função para Criar um Novo Grafo
func NovoGrafo() *Grafo {
    return &Grafo{nos: make(map[int][]int)}
}

// Adiciona uma aresta ao grafo (conexão bidirecional)
func (g *Grafo) AdicionarAresta(u, v int) {
    g.nos[u] = append(g.nos[u], v)
    g.nos[v] = append(g.nos[v], u)
}

// Remove uma aresta do grafo
func (g *Grafo) RemoverAresta(u, v int) {
    g.nos[u] = remover(g.nos[u], v)
    g.nos[v] = remover(g.nos[v], u)
}

// Função auxiliar para remover um valor de um slice
func remover(slice []int, valor int) []int {
    for i, v := range slice {
        if v == valor {
            return append(slice[:i], slice[i+1:]...)
        }
    }
    return slice
}

// Função para carregar um grafo a partir de um arquivo .edges
func CarregarGrafoEdges(arquivo string) *Grafo {
    g := NovoGrafo()

    file, err := os.Open(arquivo)
    if err != nil {
        panic("Erro ao abrir o arquivo: " + err.Error())
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        linha := scanner.Text()

        // Divide a linha em dois números (u e v)
        partes := strings.Fields(linha)
        if len(partes) != 2 {
            continue // Ignora linhas inválidas
        }

        u, err1 := strconv.Atoi(partes[0])
        v, err2 := strconv.Atoi(partes[1])
        if err1 != nil || err2 != nil {
            continue // Ignora linhas que não contêm números válidos
        }

        g.AdicionarAresta(u, v) // Adiciona a aresta ao grafo
    }

    if err := scanner.Err(); err != nil {
        panic("Erro ao ler o arquivo: " + err.Error())
    }

    return g
}

func CarregarGrafoTxt(caminhoArquivo string) *Grafo {
    g := NovoGrafo()

    arquivo, err := os.Open(caminhoArquivo)
    if err != nil {
        panic(err)
    }
    defer arquivo.Close()

    scanner := bufio.NewScanner(arquivo)
    for scanner.Scan() {
        linha := scanner.Text()
        // Divide a linha em dois números (nó origem e nó destino)
        partes := strings.Fields(linha)
        if len(partes) != 2 {
            continue
        }

        u, err1 := strconv.Atoi(partes[0])
        v, err2 := strconv.Atoi(partes[1])
        if err1 != nil || err2 != nil {
            continue
        }

        // Adiciona a aresta no grafo
        g.AdicionarAresta(u, v)
    }

    if err := scanner.Err(); err != nil {
        panic(err)
    }

    return g
}