package main

type Grafo struct {
    nos map[int][]int
}

func NovoGrafo() *Grafo {
    return &Grafo{nos: make(map[int][]int)}
}

func (g *Grafo) AdicionarAresta(u, v int) {
    g.nos[u] = append(g.nos[u], v)
    g.nos[v] = append(g.nos[v], u)
}

func (g *Grafo) RemoverAresta(u, v int) {
    g.nos[u] = remover(g.nos[u], v)
    g.nos[v] = remover(g.nos[v], u)
}

func remover(slice []int, valor int) []int {
    for i, v := range slice {
        if v == valor {
            return append(slice[:i], slice[i+1:]...)
        }
    }
    return slice
}
