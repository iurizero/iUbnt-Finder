# iUbnt-Finder

Uma ferramenta de linha de comando para descoberta e análise de dispositivos Ubiquiti em rede local.

## 📋 Descrição

O iUbnt-Finder é uma ferramenta Python que permite descobrir e analisar dispositivos Ubiquiti em sua rede local. Ele utiliza o protocolo de descoberta proprietário da Ubiquiti para identificar dispositivos, coletando informações como:

- Endereço IP
- Endereço MAC
- Nome do dispositivo
- Modelo do hardware
- Versão do firmware
- Nome da rede (ESSID)

## ⚠️ Importante: Permissões

**Esta ferramenta requer privilégios de root/administrador para funcionar corretamente.** Isso é necessário porque o programa precisa acessar diretamente as interfaces de rede para enviar e receber pacotes UDP.

### Linux
```bash
# Execute com sudo
sudo python3 iUbnt.py

# Ou, se estiver usando um ambiente virtual
sudo .venv/bin/python iUbnt.py
```

### Windows
Execute o prompt de comando como administrador e então execute o script.

## 🚀 Instalação

### Pré-requisitos

- Python 3.7 ou superior
- pip (gerenciador de pacotes Python)
- Privilégios de root/administrador

### Instalação das Dependências

```bash
# Clone o repositório
git clone https://github.com/seu-usuario/iUbnt-Finder.git
cd iUbnt-Finder

# Instale as dependências
pip install -r requirements.txt
```

## 📦 Dependências

O projeto requer as seguintes dependências:

- scapy>=2.5.0
- typing-extensions>=4.0.0

## 🎯 Uso

Para executar o scanner, você **deve** usar privilégios de administrador:

```bash
# Linux
sudo python3 iUbnt.py

# Se estiver usando ambiente virtual
sudo .venv/bin/python iUbnt.py

# Windows (como administrador)
python iUbnt.py
```

### Solução de Problemas Comuns

Se você encontrar o erro "PermissionError: [Errno 1] Operation not permitted", isso significa que o script não está sendo executado com privilégios de administrador. Certifique-se de:

1. Usar `sudo` no Linux
2. Executar como administrador no Windows
3. Verificar se seu usuário tem permissões de administrador

### Exemplo de Saída

```
Iniciando análise de rede...

Dispositivos encontrados: 2

=== UniFi AP AC Pro ===
IP: 192.168.1.100
Nome: AP-Sala
Modelo: UAP-AC-PRO
Versão: 6.5.28
Rede: Rede-Corporativa
MAC: 00:11:22:33:44:55

=== UniFi Switch 24 ===
IP: 192.168.1.101
Nome: Switch-Patrimonio
Modelo: US-24-250W
Versão: 6.5.28
Rede: Rede-Corporativa
MAC: 00:11:22:33:44:56
```

## 🔒 Observações de Segurança

- A ferramenta **requer** privilégios de root/administrador para executar operações de rede
- Recomenda-se executar apenas em redes que você tem permissão para analisar
- O scanner utiliza broadcast UDP na porta 10001
- Nunca execute scripts Python com privilégios de administrador sem verificar o código fonte

## ⚠️ Limitações

- Funciona apenas com dispositivos Ubiquiti que suportam o protocolo de descoberta
- Requer que os dispositivos estejam na mesma sub-rede
- Pode ser bloqueado por firewalls que bloqueiam tráfego UDP

## 🤝 Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou enviar pull requests.

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## ⚡️ Desenvolvimento Rápido

Para desenvolvimento rápido, você pode usar o ambiente virtual Python:

```bash
# Criar ambiente virtual
python3 -m venv venv

# Ativar ambiente virtual
source venv/bin/activate  # Linux/Mac
# ou
.\venv\Scripts\activate  # Windows

# Instalar dependências
pip install -r requirements.txt
```

## 📚 Documentação Adicional

Para mais informações sobre o protocolo de descoberta Ubiquiti, consulte a documentação oficial da Ubiquiti. 