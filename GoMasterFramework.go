package main

import (
	"bufio"
	"embed"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var defaultWordlist embed.FS

var (
	target   string
	wordlist string = "EMBEDDED"
	ports           = "80,443,22,3306"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("--- Go Master Framework ---")
	fmt.Println("Digite 'help' para comandos.")

	for {
		fmt.Print("\nGMF > ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		parts := strings.Split(input, " ")
		command := parts[0]

		switch command {
		case "set":
			handleSet(parts)
		case "show":
			handleShow()
		case "scan":
			if target == "" {
				fmt.Println("[-] Erro: Defina 'target' (IP ou Domínio) antes de rodar o scan de portas.")
			} else {
				executeScanPort()
			}
		case "run":
			if target == "" || wordlist == "" {
				fmt.Println("[-] Erro: Defina 'target' e 'wordlist' antes de rodar o brute force.")
			} else {
				executeScanSubdomain()
			}
		case "help":
			fmt.Println("\nComandos:\n  set <option> <value> : target, wordlist, ports\n  show                 : Ver configs\n  scan                 : Scan de portas no target\n  run                  : Brute force de subdomínios\n  exit                 : Sair")
		case "exit":
			return
		default:
			if command != "" {
				fmt.Printf("Comando desconhecido: %s\n", command)
			}
		}
	}
}

func executeScanPort() {
	fmt.Printf("[*] Escaneando portas em: %s\n", target)
	portList := strings.Split(ports, ",")

	for _, p := range portList {
		addr := net.JoinHostPort(target, p)
		conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
		if err == nil {
			fmt.Printf("[+] Porta %s: ABERTA\n", p)
			conn.Close()
		}
	}
	fmt.Println("[*] Scan de portas finalizado.")
}

func executeScanSubdomain() {
	fmt.Printf("[*] Buscando subdomínios para: %s\n", target)

	var scanner *bufio.Scanner

	if wordlist == "EMBEDDED" {
		f, err := defaultWordlist.Open("subdomains-top1million-110000.txt")
		if err != nil {
			fmt.Printf("[-] Erro ao carregar wordlist embutida: %v\n", err)
			return
		}
		defer f.Close()
		scanner = bufio.NewScanner(f)
		fmt.Println("[+] Usando wordlist padrão (embutida)")
	} else {
		arquivo, err := os.Open(wordlist)
		if err != nil {
			fmt.Printf("[-] Erro ao abrir wordlist externa: %v\n", err)
			return
		}
		defer arquivo.Close()
		scanner = bufio.NewScanner(arquivo)
		fmt.Printf("[+] Usando wordlist externa: %s\n", wordlist)
	}

	foundCount := 0
	for scanner.Scan() {
		sub := strings.TrimSpace(scanner.Text())
		if sub == "" || strings.HasPrefix(sub, "#") {
			continue
		}
		subAlvo := fmt.Sprintf("%s.%s", sub, target)

		ips, err := net.LookupHost(subAlvo)
		if err == nil {
			fmt.Printf("[FOUND] %s -> %v\n", subAlvo, ips)
			foundCount++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("[-] Erro durante a leitura: %v\n", err)
	}
	fmt.Printf("[*] Brute force finalizado. %d subdomínios encontrados.\n", foundCount)
}

func handleSet(parts []string) {
	if len(parts) < 3 {
		fmt.Println("Uso: set <campo> <valor>")
		return
	}
	campo := strings.ToLower(parts[1])
	valor := parts[2]

	switch campo {
	case "target":
		target = valor
		fmt.Printf("[+] Target: %s\n", target)
	case "wordlist":
		wordlist = valor
		fmt.Printf("[+] Wordlist: %s\n", wordlist)
	case "ports":
		ports = valor
		fmt.Printf("[+] Ports: %s\n", ports)
	}
}

func handleShow() {
	fmt.Printf("\nSESSÃO ATUAL:\n  TARGET:   %s\n  WORDLIST: %s\n  PORTS:    %s\n", target, wordlist, ports)
}
