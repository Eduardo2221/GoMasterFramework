package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// Agora recebemos as variáveis como parâmetros
func PortScan(alvo string, portas string, threads int) {
	fmt.Printf("[*] Escaneando portas em: %s\n", alvo)
	fmt.Printf("[*] Utilizando %d workers...\n", threads)
	fmt.Println("--------------------------------------------------")

	// Higiene: remove espaços caso o usuário tenha digitado "80, 443, 8080"
	portas = strings.ReplaceAll(portas, " ", "")
	portList := strings.Split(portas, ",")

	// Buffer inteligente (mantém o consumo de memória baixo, mesmo para 65 mil portas)
	tarefas := make(chan string, threads*2)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var count int

	// 1. INICIA OS WORKERS
	for i := 1; i <= threads; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for p := range tarefas {
				addr := net.JoinHostPort(alvo, p)
				conn, err := net.DialTimeout("tcp", addr, 1*time.Second)

				// Contador de progresso protegido contra concorrência
				mu.Lock()
				count++
				if count%500 == 0 { // Avisa a cada 500 portas testadas
					fmt.Printf("[*] Progresso: %d portas testadas...\n", count)
				}
				mu.Unlock()

				// Se conectou, a porta está aberta!
				if err == nil {
					fmt.Printf("[Worker %d] [+] Porta %s: ABERTA\n", workerID, p)
					conn.Close() // Importante fechar a conexão para não esgotar as portas locais
				}
			}
		}(i)
	}

	// 2. ALIMENTA OS WORKERS
	for _, p := range portList {
		if p != "" { // Evita enviar strings vazias caso a string termine em vírgula
			tarefas <- p
		}
	}

	// 3. ENCERRAMENTO
	close(tarefas)
	wg.Wait()

	fmt.Println("--------------------------------------------------")
	fmt.Printf("[*] Scan de portas finalizado! Total testado: %d\n", count)
}
