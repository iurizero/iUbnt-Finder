package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"iubnt-finder/internal/scanner"
)

func main() {
	timeout := flag.Duration("timeout", 5*time.Second, "tempo de espera para respostas de descoberta")
	targetsFlag := flag.String("targets", "", "lista de destinos UDP separados por vírgula; vazio usa broadcast")
	bindFlag := flag.String("bind", "", "IP local de origem para selecionar a interface de saída")
	flag.Parse()

	fmt.Println("\nIniciando análise de rede...")

	sc := scanner.NewScanner(*timeout)
	sc.Config.LocalAddr = strings.TrimSpace(*bindFlag)
	if trimmed := strings.TrimSpace(*targetsFlag); trimmed != "" {
		for _, target := range strings.Split(trimmed, ",") {
			target = strings.TrimSpace(target)
			if target != "" {
				sc.Config.Targets = append(sc.Config.Targets, target)
			}
		}
	}

	devices, err := sc.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "erro: %v\n", err)
		os.Exit(1)
	}

	if len(devices) == 0 {
		fmt.Println("\nNenhum dispositivo localizado")
		fmt.Println()
		return
	}

	fmt.Printf("\nDispositivos encontrados: %d\n", len(devices))
	for _, device := range devices {
		fmt.Print(scanner.FormatDeviceInfo(device))
	}
}
