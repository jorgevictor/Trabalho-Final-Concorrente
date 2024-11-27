package main

import (
    "fmt"
    "time"
)

func medirTempoExecucao(nome string, f func()) time.Duration {
    inicio := time.Now()
    f()
    duracao := time.Since(inicio)
    fmt.Printf("\n%s levou %s\n", nome, duracao)
    return duracao
}

func imprimirComunidades(comunidades [][]int) {
    fmt.Printf("    Total de Comunidades: %d\n", len(comunidades))
    for i, comunidade := range comunidades {
        fmt.Printf("        Comunidade %d: %v \n(Tamanho: %d)\n\n",
            i+1, comunidade, len(comunidade))
    }
}
