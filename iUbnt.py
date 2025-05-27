#!/usr/bin/env python3

import random as rnd
from scapy.all import *
from typing import Dict, List, Tuple, Any
from dataclasses import dataclass
from functools import reduce

@dataclass
class NetworkConfig:
    broadcast_mac: str = "ff:ff:ff:ff:ff:ff"
    broadcast_ip: str = "255.255.255.255"
    target_port: int = 10001
    min_port: int = 1024
    max_port: int = 65535
    timeout: int = 5

class ProtocolDecoder:
    def __init__(self):
        self._markers = {
            'a': '01',  # Device ID
            'b': '02',  # Extended ID
            'c': '03',  # Version
            'd': '0b',  # Name
            'e': '0c',  # Model
            'f': '0d',  # Network
            'g': '14'   # Full Model
        }
        self._init_packet = bytes.fromhex('01000000')
        self._valid_header = '010000'
    
    def _create_probe(self, config: NetworkConfig) -> Packet:
        sport = rnd.randint(config.min_port, config.max_port)
        return (
            Ether(dst=config.broadcast_mac)/
            IP(dst=config.broadcast_ip)/
            UDP(sport=sport, dport=config.target_port)/
            Raw(self._init_packet)
        )

    def _extract_device_data(self, packet_data: bytes) -> Dict[str, Any]:
        if packet_data[0:3].hex() != self._valid_header:
            return None

        result = {}
        pos = 3
        remaining = int(packet_data[pos].hex(), 16) - 1
        pos += 1

        while remaining > 0:
            marker = packet_data[pos].hex()
            pos += 1
            remaining -= 1

            size = int(packet_data[pos:pos+2].hex(), 16)
            pos += 2
            remaining -= 2

            value = packet_data[pos:pos+size].decode('utf-8', errors='ignore')
            pos += size
            remaining -= size

            # Map markers to properties using a more complex approach
            mapping = {
                self._markers['d']: ('device_name', value),
                self._markers['g']: ('product_model', value),
                self._markers['e']: ('model_id', value),
                self._markers['c']: ('build_version', value),
                self._markers['f']: ('network_id', value)
            }
            
            if marker in mapping:
                key, val = mapping[marker]
                result[key] = val

        return result

class NetworkScanner:
    def __init__(self):
        self.config = NetworkConfig()
        self.decoder = ProtocolDecoder()
        conf.checkIPaddr = False

    def _process_response(self, response: Tuple[Packet, Packet]) -> Dict[str, Any]:
        _, received = response
        device_data = self.decoder._extract_device_data(received[IP].load)
        
        if not device_data:
            return None

        return {
            **device_data,
            'ip_address': received[IP].src,
            'mac_address': received[Ether].src.upper()
        }

    def scan(self) -> List[Dict[str, Any]]:
        probe = self.decoder._create_probe(self.config)
        responses, _ = srp(probe, multi=True, verbose=0, timeout=self.config.timeout)
        
        return list(filter(None, map(self._process_response, responses)))

def format_device_info(device: Dict[str, Any]) -> str:
    template = """
    === {model} ===
    IP: {ip}
    Nome: {name}
    Modelo: {model_id}
    Versão: {version}
    Rede: {network}
    MAC: {mac}
    """
    
    return template.format(
        model=device.get('product_model', 'Desconhecido'),
        ip=device.get('ip_address', 'N/A'),
        name=device.get('device_name', 'N/A'),
        model_id=device.get('model_id', 'N/A'),
        version=device.get('build_version', 'N/A'),
        network=device.get('network_id', 'N/A'),
        mac=device.get('mac_address', 'N/A')
    )

def main():
    print("\nIniciando análise de rede...")
    scanner = NetworkScanner()
    devices = scanner.scan()
    
    if devices:
        print(f"\nDispositivos encontrados: {len(devices)}")
        print(''.join(map(format_device_info, devices)))
    else:
        print("\nNenhum dispositivo localizado\n")

if __name__ == "__main__":
    main()