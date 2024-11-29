package main

import (
    "fmt"
    "time"
)

func medirTempoExecucao(nome string, f func()) time.Duration {
    inicio := time.Now()
    f()
    duracao := time.Since(inicio)
    fmt.Printf("    Tempo: %s\n", duracao)
    return duracao
}

func imprimirComunidades(comunidades [][]int) {
    fmt.Printf("\n    Total de Comunidades: %d\n", len(comunidades))
    for i, comunidade := range comunidades {
        fmt.Printf("    Comunidade %d: %v (Tamanho: %d)\n",
            i+1, comunidade, len(comunidade))
    }
}
