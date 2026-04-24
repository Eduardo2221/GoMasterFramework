package main

import (
	"fmt"
	"strings"
)

func HandleSet(parts []string) error {
	// Verifica se o usuário digitou comando, campo e valor
	if len(parts) < 3 {
		return fmt.Errorf("uso incorreto. Exemplo: set target example.com")
	}

	campo := strings.ToLower(parts[1])
	valor := parts[2]

	switch campo {
	case "target":
		target = valor
	case "wordlist":
		wordlist = valor
	case "ports":
		ports = valor
	case "threads":
		var t int
		// Tenta converter o valor digitado para número inteiro
		_, err := fmt.Sscanf(valor, "%d", &t)
		if err != nil {
			return fmt.Errorf("o valor de threads deve ser um número, ex: set threads 20")
		}
		threads = t
	default:
		return fmt.Errorf("campo desconhecido: %s", campo)
	}

	return nil // Retorna vazio, indicando sucesso!
}

func HandleShow() string {
	var sb strings.Builder

	sb.WriteString("\n--- CONFIGURAÇÕES ATUAIS ---\n")
	sb.WriteString(fmt.Sprintf(" Target:   %s\n", target))
	sb.WriteString(fmt.Sprintf(" Wordlist: %s\n", wordlist))
	sb.WriteString(fmt.Sprintf(" Portas:   %s\n", ports))
	sb.WriteString(fmt.Sprintf(" Threads:  %d\n", threads))

	return sb.String()
}
