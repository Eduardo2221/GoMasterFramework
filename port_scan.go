package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Representa uma porta encontrada
type PortResult struct {
	Port  string
	State string
}

// Representa o resumo final do scan de portas
type PortScanReport struct {
	Target       string
	TotalScanned int
	OpenPorts    []PortResult
}

// Agora recebemos as variáveis como parâmetros
func PortScan(alvo string, portas string, threads int) (PortScanReport, error) {
	var report PortScanReport
	report.Target = alvo

	portList, err := ParsePorts(portas)
	if err != nil {
		// Se o usuário digitou um range inválido, o scan para e retornamos o erro
		return report, fmt.Errorf("erro nas portas: %v", err)
	}

	// Buffer inteligente
	tarefas := make(chan string, threads*2)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var count int
	var encontrados []PortResult // Lista onde vamos guardar as portas abertas

	// INICIA OS WORKERS
	for i := 1; i <= threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for p := range tarefas {
				addr := net.JoinHostPort(alvo, p)

				// A tentativa de conexão fica FORA do Mutex (para não travar os outros workers)
				conn, err := net.DialTimeout("tcp", addr, 1*time.Second)

				// Bloqueia apenas na hora de somar e salvar o resultado
				mu.Lock()
				count++

				// Se conectou (err == nil), a porta está aberta!
				if err == nil {
					encontrados = append(encontrados, PortResult{
						Port:  p,
						State: "ABERTA",
					})
					conn.Close() // Fecha a conexão
				}
				mu.Unlock()
			}
		}()
	}

	// ALIMENTA OS WORKERS
	for _, p := range portList {
		if p != "" { // Evita enviar strings vazias
			tarefas <- p
		}
	}

	// ENCERRAMENTO
	close(tarefas)
	wg.Wait()

	// Preenche o relatório final
	report.TotalScanned = count
	report.OpenPorts = encontrados

	return report, nil
}
