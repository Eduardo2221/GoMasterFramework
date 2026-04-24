package main

import (
	"fmt"
	"strings"
)

func HandleSet(parts []string) {
	if len(parts) < 3 {
		fmt.Println("Uso: set <campo> <valor>")
		return
	}
	campo := strings.ToLower(parts[1])
	valor := parts[2]

	switch campo {
	case "target":
		target = valor
		fmt.Printf("[+] Target definido: %s\n", target)
	case "wordlist":
		wordlist = valor
		fmt.Printf("[+] Wordlist definida: %s\n", wordlist)
	case "ports":
		ports = valor
		fmt.Printf("[+] Portas definidas: %s\n", ports)
	case "threads":
		fmt.Sscanf(valor, "%d", &threads)
		fmt.Printf("[+] Threads definidas: %d\n", threads)
	}
}

func HandleShow() {
	fmt.Printf("\n--- CONFIGURAÇÕES ATUAIS ---\n")
	fmt.Printf(" Target:   %s\n Wordlist: %s\n Portas:   %s\n Threads:  %d\n", target, wordlist, ports, threads)
}
