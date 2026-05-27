package main

import (
	_ "embed"
	"fmt"
	"os"
	"strings"
)

// 1º e 2º: Embutindo as wordlists padrão
//
//go:embed defaults/dns.txt
var defaultDNSWordlist string

//go:embed defaults/dir.txt
var defaultDIRWordlist string

// 4º: Perfis de Portas e Extensões padrão
var (
	defaultPortsLow  = "21,22,23,25,53,80,110,443"
	defaultPortsMid  = "21,22,23,25,53,80,110,135,139,443,445,1433,3306,3389,8080"
	defaultPortsHigh = "1-65535"

	defaultExtsLow  = "html,php,txt"
	defaultExtsMid  = "html,php,txt,json,xml,log,zip"
	defaultExtsHigh = "html,php,txt,json,xml,log,zip,tar.gz,bak,old,conf,config,sql"
)

// CarregarConfiguracoes processa as flags e retorna os dados finais que o scanner vai usar
func CarregarConfiguracoes(wordlistFlag, dnsFlag, dirFlag, portProfile, extProfile string) ([]string, []string, string, string) {
	var finalDNSWords []string
	var finalDIRWords []string
	var finalPorts string
	var finalExts string

	// ==========================================
	// REGRA 3º: Lógica das Wordlists (Prioridades)
	// ==========================================
	if wordlistFlag != "" {
		fmt.Printf("[*] Usando wordlist customizada única para ambos: %s\n", wordlistFlag)
		linhas, err := lerArquivoExterno(wordlistFlag)
		if err != nil {
			fmt.Printf("Erro ao ler wordlist: %v\n", err)
			os.Exit(1)
		}
		finalDNSWords = linhas
		finalDIRWords = linhas
	} else {
		// Tratamento de DNS
		if dnsFlag != "" {
			linhas, _ := lerArquivoExterno(dnsFlag)
			finalDNSWords = linhas
		} else {
			finalDNSWords = strings.Split(defaultDNSWordlist, "\n")
		}

		// Tratamento de DIR
		if dirFlag != "" {
			linhas, _ := lerArquivoExterno(dirFlag)
			finalDIRWords = linhas
		} else {
			finalDIRWords = strings.Split(defaultDIRWordlist, "\n")
		}
	}

	// ==========================================
	// REGRA 4º: Perfis de Portas
	// ==========================================
	switch strings.ToLower(portProfile) {
	case "mid":
		finalPorts = defaultPortsMid
	case "high":
		finalPorts = defaultPortsHigh
	default:
		finalPorts = defaultPortsLow
	}

	// Perfis de Extensões
	switch strings.ToLower(extProfile) {
	case "mid":
		finalExts = defaultExtsMid
	case "high":
		finalExts = defaultExtsHigh
	default:
		finalExts = defaultExtsLow
	}

	// Retorna os slices limpos e as strings de portas/extensões
	return limparLinhas(finalDNSWords), limparLinhas(finalDIRWords), finalPorts, finalExts
}

// Funções auxiliares mantidas aqui para não poluir o main
func lerArquivoExterno(caminho string) ([]string, error) {
	conteudo, err := os.ReadFile(caminho)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(conteudo), "\n"), nil
}

func limparLinhas(linhas []string) []string {
	var limpas []string
	for _, l := range linhas {
		l = strings.TrimSpace(l)
		if l != "" {
			limpas = append(limpas, l)
		}
	}
	return limpas
}
