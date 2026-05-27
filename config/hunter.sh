#!/bin/bash
# hunter.sh - Script de automação

ALVO=$1

echo "[*] Coletando subdomínios usando o GMF..."
./gmf -m dns -t $ALVO -w subdomains.txt > dns_results.txt

echo "[*] Coletando diretórios usando o GMF..."
./gmf -m dir -t $ALVO -w dirs.txt > dir_results.txt

echo "[+] Automação concluída!"