# iUbnt-Finder

Uma ferramenta de linha de comando para descoberta e análise de dispositivos Ubiquiti em rede local.

## 📋 Descrição

O iUbnt-Finder é uma ferramenta em Go que permite descobrir e analisar dispositivos Ubiquiti em sua rede local. Ele utiliza o protocolo de descoberta proprietário da Ubiquiti para identificar dispositivos, coletando informações como:

- Endereço IP
- Endereço MAC
- Nome do dispositivo
- Modelo do hardware
- Versão do firmware
- Nome da rede (ESSID)

## ⚠️ Importante: Permissões

**Esta ferramenta pode exigir privilégios elevados em alguns ambientes de rede.** Isso depende de políticas locais para broadcast/UDP e leitura da tabela ARP.

### Linux
```bash
go build -o iubnt-finder ./cmd/iubnt-finder
sudo ./iubnt-finder
```

## 🚀 Instalação

### Pré-requisitos

- Go 1.26 ou superior
- Privilégios de root/administrador podem ser necessários em alguns ambientes de rede

### Build

```bash
go build -o iubnt-finder ./cmd/iubnt-finder
```

### Execução

```bash
./iubnt-finder
```

Se quiser testar sem compilar:

```bash
go run ./cmd/iubnt-finder
```

## Modos De Uso

- `scan local`: use sem flags. Faz broadcast na rede local.
- `scan by targets`: use `-targets` com um ou mais IPs.
- `scan by interface`: use `-bind` para escolher o IP de origem da interface.

## Uso

### Sintaxe

```bash
./iubnt-finder [flags]
```

### Flags disponíveis

| Flag | Padrão | Descrição |
| --- | --- | --- |
| `-timeout` | `5s` | Tempo máximo de espera por respostas de descoberta. |
| `-targets` | vazio | Lista de destinos UDP separada por vírgulas. Se vazio, usa broadcast. |
| `-bind` | vazio | IP local de origem para selecionar a interface de saída. |
| `-h`, `-help` | n/a | Exibe a ajuda padrão do programa. |

### Comportamento padrão

Sem flags, o programa:

1. Usa broadcast UDP para `255.255.255.255:10001`
2. Espera até `5s` por respostas
3. Resolve o MAC via `/proc/net/arp` quando disponível

### Modos de uso recomendados

- Descoberta na rede local:

```bash
./iubnt-finder
```

- Descoberta com timeout maior:

```bash
./iubnt-finder -timeout 10s
```

- Varredura por IPs conhecidos:

```bash
./iubnt-finder -targets 192.168.1.10,192.168.1.11
```

- Varredura usando porta explícita:

```bash
./iubnt-finder -targets 192.168.1.10:10001,192.168.1.11:2000
```

- Varredura forçando a interface/IP de origem:

```bash
./iubnt-finder -bind 10.0.0.10
```

- Combinação de interface e targets:

```bash
./iubnt-finder -bind 10.0.0.10 -targets 10.0.0.50,10.0.0.51
```

### Quando usar cada modo

- `scan local`: quando você está na mesma rede ou VLAN do equipamento
- `scan by targets`: quando você já conhece o IP do rádio ou quer testar alvos específicos
- `scan by interface`: quando sua máquina tem mais de uma interface ou você quer sair por uma rede específica

### Solução de Problemas Comuns

Se o binário não conseguir enviar o broadcast ou ler respostas, verifique:

1. Se a rede local permite UDP broadcast na porta `10001`
2. Se firewall local ou de rede está bloqueando as respostas
3. Se o dispositivo alvo está na mesma sub-rede

### Varredura Unicast

Para apontar IPs específicos, use `-targets` com uma lista separada por vírgulas:

```bash
./iubnt-finder -targets 192.168.1.10,192.168.1.11
```

Se um alvo não tiver porta explícita, o scanner usa `10001` automaticamente.

### Varredura por Interface

Se você quiser forçar a origem do tráfego em uma interface/IP específico, use `-bind`:

```bash
./iubnt-finder -bind 10.0.0.10
```

Isso é útil quando o PC tem mais de uma interface ativa ou quando você quer testar a partir de uma sub-rede específica.

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

- A ferramenta pode precisar de privilégios de rede dependendo do ambiente
- Recomenda-se executar apenas em redes que você tem permissão para analisar
- O scanner utiliza broadcast UDP na porta 10001
- Verifique o código fonte antes de executar com privilégios elevados

## ⚠️ Limitações

- Funciona apenas com dispositivos Ubiquiti que suportam o protocolo de descoberta
- Requer que os dispositivos estejam na mesma sub-rede
- Pode ser bloqueado por firewalls que bloqueiam tráfego UDP
- A leitura de MAC usa `/proc/net/arp`, então a validação atual é principalmente em Linux

## Exemplo Completo

```bash
./iubnt-finder -bind 10.0.0.10 -targets 10.0.0.50,10.0.0.51 -timeout 8s
```

Esse comando:

1. Usa o IP `10.0.0.10` como origem
2. Envia probes para `10.0.0.50:10001` e `10.0.0.51:10001`
3. Aguarda até `8s` por respostas
4. Exibe os dispositivos encontrados no terminal

## 🤝 Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou enviar pull requests.

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está licenciado sob a licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 📚 Documentação Adicional

Para mais informações sobre o protocolo de descoberta Ubiquiti, consulte a documentação oficial da Ubiquiti.
