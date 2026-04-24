package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	target   string = "example.com"
	wordlist string = "wordlist.txt"
	ports    string = "21,22,80,443,3306,8080"
	threads  int    = 10
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("--- Go Master Framework (GMF) ---")
	fmt.Println("Digite 'help' para comandos.")

	for {
		fmt.Print("\nGMF > ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		parts := strings.Split(input, " ")
		command := parts[0]

		switch command {

		case "set":
			// Supondo que você já tem a HandleSet definida
			HandleSet(parts)

		case "show":
			// Supondo que você já tem a HandleShow definida
			HandleShow()

		case "dns":
			if target == "" {
				fmt.Println("[-] Erro: Defina 'target' antes.")
			} else {
				EnumerateDNS(target, wordlist, threads)
			}

		case "dir":
			if target == "" {
				fmt.Println("[-] Erro: Defina 'target' antes.")
			} else {
				EnumerateDIR(target, wordlist, threads)
			}

		case "port":
			if target == "" {
				fmt.Println("[-] Erro: Defina 'target' antes.")
			} else {
				PortScan(target, ports, threads)
			}

		case "help":
			fmt.Println("\nComandos Disponíveis:")
			fmt.Println("  set <campo> <valor> : Configura variáveis (target, wordlist, ports, threads)")
			fmt.Println("  show                : Exibe as configurações atuais")
			fmt.Println("  dns                 : Inicia a enumeração de subdomínios (DNS)")
			fmt.Println("  dir                 : Inicia a enumeração de diretórios (Web)")
			fmt.Println("  port                : Inicia o scan de portas TCP")
			fmt.Println("  help                : Exibe este menu de ajuda")
			fmt.Println("  exit                : Fecha o framework")

		case "exit":
			fmt.Println("[*] Saindo do GMF. Até a próxima!")
			return

		default:
			fmt.Printf("[-] Comando desconhecido: %s\n", command)
		}
	}
}
