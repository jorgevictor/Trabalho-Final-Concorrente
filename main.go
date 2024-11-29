package main

import (
    "fmt"
    "strings"
	"math/rand"
	"time"
	"runtime"
    "io/ioutil"
	"path/filepath"
)


func main() {
	runtime.GOMAXPROCS(4) // Define o uso de threads
	rand.Seed(time.Now().UnixNano()) // Inicializa o gerador de números aleatórios
    fmt.Println("Teste de Grafos:")

    testes := []struct {
        nome  string
        grafo *Grafo
    }{
        {"Ciclo pequeno com conexão", func() *Grafo {
            g := NovoGrafo()
            g.AdicionarAresta(1, 2)
            g.AdicionarAresta(2, 3)
            g.AdicionarAresta(3, 4)
            g.AdicionarAresta(4, 1)
            g.AdicionarAresta(3, 5)
            return g
        }()},
        {"Grafo desconexo pequeno", func() *Grafo {
            g := NovoGrafo()
            g.AdicionarAresta(1, 2)
            g.AdicionarAresta(3, 4)
            return g
        }()},
        {"Grafo linear pequeno", func() *Grafo {
            g := NovoGrafo()
            g.AdicionarAresta(1, 2)
            g.AdicionarAresta(2, 3)
            g.AdicionarAresta(3, 4)
            return g
        }()},
        {"Grafo completo pequeno", func() *Grafo {
            g := NovoGrafo()
            g.AdicionarAresta(1, 2)
            g.AdicionarAresta(1, 3)
            g.AdicionarAresta(1, 4)
            g.AdicionarAresta(2, 3)
            g.AdicionarAresta(2, 4)
            g.AdicionarAresta(3, 4)
            return g
        }()},
        {"Grafo raro pequeno", func() *Grafo {
            g := NovoGrafo()
            g.AdicionarAresta(1, 2)
            g.AdicionarAresta(3, 4)
            g.AdicionarAresta(5, 6)
            return g
        }()},
        {"Grafo denso pequeno", func() *Grafo {
            g := NovoGrafo()
            for i := 1; i <= 5; i++ {
                for j := i + 1; j <= 5; j++ {
                    g.AdicionarAresta(i, j)
                }
            }
            return g
        }()},
        {"Grafo pequeno com comunidades", func() *Grafo {
            g := NovoGrafo()
            // Comunidade 1
            g.AdicionarAresta(1, 2)
            g.AdicionarAresta(2, 3)
            g.AdicionarAresta(3, 1)
            // Comunidade 2
            g.AdicionarAresta(4, 5)
            g.AdicionarAresta(5, 6)
            g.AdicionarAresta(6, 4)
            // Conexão entre comunidades
            g.AdicionarAresta(3, 4)
            return g
        }()},
        {"Grafo escalado médio", func() *Grafo {
            g := NovoGrafo()
            for i := 1; i <= 100; i++ {
                g.AdicionarAresta(i, i+1)
                if i%10 == 0 && i < 100 {
                    g.AdicionarAresta(i, i+10)
                }
            }
            return g
        }()},
    }

    for _, teste := range testes {
        fmt.Printf("\n\n")
        fmt.Println(strings.Repeat("-", 40))
        fmt.Printf("%s\n", teste.nome)
        fmt.Println(strings.Repeat("-", 40))

        seqGrafo := teste.grafo
        conGrafo := teste.grafo

        fmt.Printf("\n  Sequencial: ")
        medirTempoExecucao("Sequencial", func() {
            comunidades := girvanNewman(seqGrafo, false) // Detecta comunidades
            imprimirComunidades(comunidades)            // Imprime as comunidades
        })

        fmt.Printf("\n  Concorrente: ")
        medirTempoExecucao("Concorrente", func() {
            comunidades := girvanNewman(conGrafo, true) // Detecta comunidades
            imprimirComunidades(comunidades)           // Imprime as comunidades
        })
    }


    fmt.Println("\n\nTeste de Desempenho para grafos grandes:")


    grafoEscaladoGrande := func() *Grafo {
        g := NovoGrafo()
        for i := 1; i <= 1000; i++ {
            g.AdicionarAresta(i, i+1)
            if i%5 == 0 && i < 1000 {
                g.AdicionarAresta(i, i+5)
            }
        }
        return g
    }()

    fmt.Println("\n  Comunidades Detectadas no grafo escalado grande:")

    fmt.Printf("  Sequencial: ")
    seqTempoEscalado := medirTempoExecucao("Sequencial em Grafo escalado grande", func() {
        comunidades := girvanNewman(grafoEscaladoGrande, false) // Detecta comunidades
        imprimirComunidades(comunidades)               // Imprime as comunidades
    })

    // Reseta o grafo grande
    grafoEscaladoGrande = func() *Grafo {
        g := NovoGrafo()
        for i := 1; i <= 1000; i++ {
            g.AdicionarAresta(i, i+1)
            if i%5 == 0 && i < 1000 {
                g.AdicionarAresta(i, i+5)
            }
        }
        return g
    }()

    fmt.Printf("\n  Concorrente: ")
    conTempoEscalado := medirTempoExecucao("Concorrente em Grafo Escalado Grande", func() {
        comunidades := girvanNewman(grafoEscaladoGrande, true) // Detecta comunidades
        imprimirComunidades(comunidades)              // Imprime as comunidades
    })

    fmt.Printf("\n\n  Diferença de Tempo (grafo escalado): Sequencial levou %s, Concorrente levou %s\n", seqTempoEscalado, conTempoEscalado)

	grafoCompletoGrande := func() *Grafo {
		g := NovoGrafo()
		n := 100 // Número de nós

		// Adiciona todas as arestas
		for i := 1; i <= n; i++ {
			for j := i + 1; j <= n; j++ {
				g.AdicionarAresta(i, j)
			}
		}

		return g
	}()

	fmt.Println("\n  Comunidades Detectadas no Grafo Completo Grande:")

	// Teste Sequencial
	fmt.Printf("  Sequencial: ")
	seqTempoCompleto := medirTempoExecucao("Sequencial em Grafo Completo Grande", func() {
		comunidades := girvanNewman(grafoCompletoGrande, false) // Detecta comunidades
		imprimirComunidades(comunidades)                       // Imprime as comunidades
	})

	// Reseta o grafo completo grande
	grafoCompletoGrande = func() *Grafo {
		g := NovoGrafo()
		n := 100 // Número de nós

		// Adiciona todas as arestas
		for i := 1; i <= n; i++ {
			for j := i + 1; j <= n; j++ {
				g.AdicionarAresta(i, j)
			}
		}

		return g
	}()

	// Teste Concorrente
	fmt.Printf("\n  Concorrente: ")
	conTempoCompleto := medirTempoExecucao("Concorrente em Grafo Completo Grande", func() {
		comunidades := girvanNewman(grafoCompletoGrande, true) // Detecta comunidades
		imprimirComunidades(comunidades)                      // Imprime as comunidades
	})

	fmt.Printf("\n\n  Diferença de Tempo (grafo completo): Sequencial levou %s, Concorrente levou %s\n", seqTempoCompleto, conTempoCompleto)

    grafoComComunidades := func() *Grafo {
		g := NovoGrafo()
		n := 1000         // Total de nós
		k := 10           // Número de comunidades
		tamanho := n / k  // Tamanho de cada comunidade

		// Conectar nós dentro de cada comunidade
		for c := 0; c < k; c++ {
			inicio := c * tamanho + 1
			fim := (c+1)*tamanho

			// Conexões densas dentro da comunidade
			for i := inicio; i <= fim; i++ {
				for j := i + 1; j <= fim; j++ {
					g.AdicionarAresta(i, j)
				}
			}
		}

		// Conectar algumas arestas entre comunidades (bem esparsas)
		for c := 0; c < k-1; c++ {
			no1 := c * tamanho + rand.Intn(tamanho) + 1 // Nó aleatório na comunidade atual
			no2 := (c+1)*tamanho - rand.Intn(tamanho)  // Nó aleatório na próxima comunidade
			g.AdicionarAresta(no1, no2)
		}

		return g
	}()


	fmt.Println("\n  Comunidades Detectadas no Grafo Com Comunidades:")

	// Teste Sequencial
	fmt.Printf("  Sequencial: ")
	seqTempoComComunidades := medirTempoExecucao("Sequencial em Grafo Com Comunidades", func() {
		comunidades := girvanNewman(grafoComComunidades, false) // Detecta comunidades
		imprimirComunidades(comunidades)                       // Imprime as comunidades
	})

	// Reseta o grafo com comunidades
	grafoComComunidades = func() *Grafo {
		g := NovoGrafo()
		n := 1000         // Total de nós
		k := 10           // Número de comunidades
		tamanho := n / k  // Tamanho de cada comunidade

		// Conectar nós dentro de cada comunidade
		for c := 0; c < k; c++ {
			inicio := c * tamanho + 1
			fim := (c+1)*tamanho

			// Conexões densas dentro da comunidade
			for i := inicio; i <= fim; i++ {
				for j := i + 1; j <= fim; j++ {
					g.AdicionarAresta(i, j)
				}
			}
		}

		// Conectar algumas arestas entre comunidades (bem esparsas)
		for c := 0; c < k-1; c++ {
			no1 := c * tamanho + rand.Intn(tamanho) + 1 // Nó aleatório na comunidade atual
			no2 := (c+1)*tamanho - rand.Intn(tamanho)  // Nó aleatório na próxima comunidade
			g.AdicionarAresta(no1, no2)
		}

		return g
	}()



	// Teste Concorrente
	fmt.Printf("\n  Concorrente: ")
	conTempoComComunidades := medirTempoExecucao("Concorrente em Grafo com comunidades", func() {
		comunidades := girvanNewman(grafoComComunidades, true) // Detecta comunidades
		imprimirComunidades(comunidades)                      // Imprime as comunidades
	})

	fmt.Printf("\n\n  Diferença de Tempo (grafo com comunidades): Sequencial levou %s, Concorrente levou %s\n", seqTempoComComunidades, conTempoComComunidades)


    // Caminho para a pasta com os arquivos .edges
	diretorio := "./facebook"
	// Lê todos os arquivos no diretório
	arquivos, err := ioutil.ReadDir(diretorio)
	if err != nil {
		panic(err)
	}
	// Processa todos os arquivos .edges
	for _, arquivo := range arquivos {
		if filepath.Ext(arquivo.Name()) == ".edges" {
			caminhoArquivo := filepath.Join(diretorio, arquivo.Name())
			fmt.Printf("\nCarregando arquivo: %s\n\n", caminhoArquivo)
			// Carrega o grafo
			grafo := CarregarGrafoEdges(caminhoArquivo)
         // Rodando o algoritmo Girvan-Newman no modo sequencial
			fmt.Println("Executando Girvan-Newman (Sequencial)...")
			inicioSequencial := time.Now()
			comunidadesSequencial := girvanNewman(grafo, false)
			tempoSequencial := time.Since(inicioSequencial)
			fmt.Printf("Comunidades detectadas (Sequencial):\n")
			imprimirComunidades(comunidadesSequencial)
         fmt.Printf("Tempo de execução (Sequencial): %v\n", tempoSequencial)
         fmt.Printf("\n")
         // Rodando o algoritmo Girvan-Newman no modo concorrente
			fmt.Println("Executando Girvan-Newman (Concorrente)...")
			inicioConcorrente := time.Now()
			comunidadesConcorrente := girvanNewman(grafo, true)
			tempoConcorrente := time.Since(inicioConcorrente)
			fmt.Printf("Comunidades detectadas (Concorrente):\n")
			imprimirComunidades(comunidadesConcorrente)
			fmt.Printf("Tempo de execução (Concorrente): %v\n", tempoConcorrente)
			
         fmt.Printf("\n\n\n")
			
		}
	}
    caminhoArquivo := "./email-Eu-core.txt"
    fmt.Printf("Carregando grafo do arquivo: %s\n", caminhoArquivo)
    grafo := CarregarGrafoTxt(caminhoArquivo)
    // Processando o grafo com o algoritmo Girvan-Newman
    fmt.Println("Executando Girvan-Newman (Sequencial)...")
    inicioSequencial := time.Now()
    comunidadesSequencial := girvanNewman(grafo, false)
    tempoSequencial := time.Since(inicioSequencial)
    fmt.Println("Comunidades detectadas (Sequencial):\n")
    imprimirComunidades(comunidadesSequencial)
    fmt.Printf("Tempo de execução (Sequencial): %v\n", tempoSequencial)
    fmt.Printf("\n")
    fmt.Println("Executando Girvan-Newman (Concorrente)...")
    inicioConcorrente := time.Now()
    comunidadesConcorrente := girvanNewman(grafo, true)
    tempoConcorrente := time.Since(inicioConcorrente)
    fmt.Println("Comunidades detectadas (Concorrente):")
    imprimirComunidades(comunidadesConcorrente)
    fmt.Printf("Tempo de execução (Concorrente): %v\n", tempoConcorrente)
    fmt.Printf("\n\n\n")
}