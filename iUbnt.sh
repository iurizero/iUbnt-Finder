#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Verifica se está rodando como root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Este script precisa ser executado como root (sudo)${NC}"
    exit 1
fi

# Verifica se o Python3 está instalado
if ! command -v python3 &> /dev/null; then
    echo -e "${RED}Python3 não está instalado. Por favor, instale-o primeiro:${NC}"
    echo "sudo apt-get update && sudo apt-get install python3 python3-venv python3-pip"
    exit 1
fi

# Verifica se python3-venv está instalado
if ! dpkg -l | grep -q python3-venv; then
    echo -e "${YELLOW}Instalando python3-venv...${NC}"
    apt-get install -y python3-venv
fi

# Cria e ativa o ambiente virtual
VENV_DIR=".venv"
if [ ! -d "$VENV_DIR" ]; then
    echo -e "${YELLOW}Criando ambiente virtual Python...${NC}"
    python3 -m venv "$VENV_DIR"
fi

# Ativa o ambiente virtual
source "$VENV_DIR/bin/activate"

# Instala o scapy no ambiente virtual
if ! python3 -c "import scapy" &> /dev/null; then
    echo -e "${YELLOW}Instalando scapy no ambiente virtual...${NC}"
    pip install scapy
fi

# Adiciona o shebang ao script Python se não existir
if ! head -n 1 iUbnt.py | grep -q "^#!.*python"; then
    echo -e "${YELLOW}Adicionando shebang ao script Python...${NC}"
    sed -i '1i#!/usr/bin/env python3' iUbnt.py
fi

# Torna o script Python executável
chmod +x iUbnt.py

# Executa o script usando o Python do ambiente virtual
echo -e "${GREEN}Executando iUbnt Finder...${NC}"
"$VENV_DIR/bin/python" iUbnt.py

# Desativa o ambiente virtual
deactivate 