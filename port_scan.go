package main

import (
	"net"
	"strings"
	"sync"
	"time"
)

// Representa uma porta encontrada
type PortResult struct {
	Port  string
	State string // Ex: "ABERTA"
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

	// Higiene: remove espaços caso o usuário tenha digitado "80, 443, 8080"
	portas = strings.ReplaceAll(portas, " ", "")
	portList := strings.Split(portas, ",")

	// Buffer inteligente
	tarefas := make(chan string, threads*2)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var count int
	var encontrados []PortResult // Lista onde vamos guardar as portas abertas

	// 1. INICIA OS WORKERS
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

	// 2. ALIMENTA OS WORKERS
	for _, p := range portList {
		if p != "" { // Evita enviar strings vazias
			tarefas <- p
		}
	}

	// 3. ENCERRAMENTO
	close(tarefas)
	wg.Wait()

	// Preenche o relatório final
	report.TotalScanned = count
	report.OpenPorts = encontrados

	return report, nil
}
