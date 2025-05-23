#!/usr/bin/env python3

from random import choice
from scapy.all import *

# Protocol constants for device identification
PROTOCOL_MARKERS = {
    'DEVICE_ID': '01',
    'DEVICE_ID_EXT': '02',
    'SOFTWARE_VERSION': '03',
    'DEVICE_ALIAS': '0b',
    'HARDWARE_CODE': '0c',
    'WIRELESS_ID': '0d',
    'HARDWARE_FULL': '14'
}

# Network protocol parameters
PROBE_MESSAGE = bytes.fromhex('01000000')
RESPONSE_HEADER = '010000'
SCAN_DURATION = 5

def scan_network_devices():
    # Disable IP address validation for broadcast
    conf.checkIPaddr = False
    
    # Construct network probe packet
    probe = (
        Ether(dst="ff:ff:ff:ff:ff:ff")/
        IP(dst="255.255.255.255")/
        UDP(sport=choice(range(1024, 65536)), dport=10001)/
        Raw(PROBE_MESSAGE)
    )
    
    # Send probe and collect responses
    responses, _ = srp(
        probe,
        multi=True,
        verbose=0,
        timeout=SCAN_DURATION
    )

    found_devices = []
    for _, response in responses:
        data = response[IP].load
        
        # Validate response format
        if data[0:3].hex() != RESPONSE_HEADER:
            continue

        device_info = {
            'network_address': response[IP].src,
            'hardware_id': response[Ether].src.upper()
        }

        # Parse response data
        data_index = 3
        data_length = int(data[data_index].hex(), 16)
        data_index += 1
        data_length -= 1

        while data_length > 0:
            marker_type = data[data_index].hex()
            data_index += 1
            data_length -= 1
            
            segment_size = int(data[data_index:data_index+2].hex(), 16)
            data_index += 2
            data_length -= 2
            
            segment_data = data[data_index:data_index+segment_size].decode('utf-8', errors='ignore')
            
            # Map response data to device properties
            if marker_type == PROTOCOL_MARKERS['DEVICE_ALIAS']:
                device_info['alias'] = segment_data
            elif marker_type == PROTOCOL_MARKERS['HARDWARE_FULL']:
                device_info['hardware_model'] = segment_data
            elif marker_type == PROTOCOL_MARKERS['HARDWARE_CODE']:
                device_info['model_code'] = segment_data
            elif marker_type == PROTOCOL_MARKERS['SOFTWARE_VERSION']:
                device_info['version'] = segment_data
            elif marker_type == PROTOCOL_MARKERS['WIRELESS_ID']:
                device_info['wireless_name'] = segment_data
                
            data_index += segment_size
            data_length -= segment_size

        found_devices.append(device_info)

    return found_devices

def display_results():
    print("\nIniciando varredura de rede...")
    network_devices = scan_network_devices()
    
    if network_devices:
        print(f"\nDetectados {len(network_devices)} dispositivo(s):")
        for device in network_devices:
            print(f"\n--- [{device.get('hardware_model', 'N/A')}] ---")
            print(f"  Endereço IP : {device.get('network_address', 'N/A')}")
            print(f"  Apelido     : {device.get('alias', 'N/A')}")
            print(f"  Código      : {device.get('model_code', 'N/A')}")
            print(f"  Versão      : {device.get('version', 'N/A')}")
            print(f"  Rede        : {device.get('wireless_name', 'N/A')}")
            print(f"  ID Hardware : {device.get('hardware_id', 'N/A')}")
    else:
        print("\nNenhum dispositivo detectado na rede\n")

if __name__ == "__main__":
    display_results()