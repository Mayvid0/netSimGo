# NetSimGo

NetSimGo is a networking simulation project developed in Go. It aims to simulate various networking concepts,encompassing the entire **TCP/IP** model, including layers beyond the physical and data link layers. For now it includes implementations of various network topologies, data link layer protocols, and data transmission simulations.

## Features

- **Device Structs**: Provides basic structures for devices like hubs, bridges, switches, and end devices.
- **Topologies**: Implements different network topologies such as point-to-point and star topologies.
- **Data Link Layer**: Includes functionalities like connecting devices, sending and receiving data, _**error correction and detection**_ using **_Hamming_** code, **token passing** _access_ control protocol and _**selective repeat**_ sliding window **flow control** protocol.


## Programming Language

The entire project is written in Go (Golang), a statically typed, compiled programming language designed for building efficient and reliable software.

## Structs and Interfaces

### Physical Layer
- **Device**: Represents a network device with attributes like Name, MACAddress, LinkStatus, PortNumber, and HasToken.
- **Hub**: Extends Device and includes additional fields such as NumberPorts and EndDevices (connected devices).
- **Connection**: Defines a connection between two devices.
- **Bridge**: Extends Device and includes fields like NumberPorts, EndDevices (hubs), and a ForwardingTable.
- **Switch**: Extends Hub and includes a SwitchingTable for packet routing.

### Data Link Layer
- **Frame**: Represents a data frame with fields like SequenceNumber, Acknowledgment, Data, Checksum, and Retransmit.

### Topologies
- **PointToPoint**: Implements a point-to-point network topology with methods for adding devices and data transmission.
- **Star**: Represents a star network topology using a hub, with methods for connecting devices and data transmission via the hub.

### Token Passing and Data Link Layer Protocols
- **Token**: Represents a token used in token passing protocols.
- **StarTopologyWithSwitch**: Extends Switch and implements data transmission using the Selective Repeat protocol and token passing.

### Data Link Layer Operations
- **hammingEncoding**: Implements Hamming encoding for error detection and correction.
- **addNoise**: Simulates noise in data transmission by flipping random bits.
- **hammingDecoding**: Implements Hamming decoding for error detection and correction.


## Input/Output Representation
- Input: Network devices, messages for data transmission, topology configurations.
- Output: Log messages indicating connections, data transmission, acknowledgments, errors, and received data.

## Formats
- MACAddress: Represented as a string in the format "XX:XX:XX:XX:XX:XX".
- Messages: Text strings for data transmission, encoded and decoded using Hamming codes.
- Frame: Data frames with sequence numbers, checksums, acknowledgment flags, and payload data.

## Implementation Highlights
- Network topologies: Point-to-point, star with a hub, and switch-based networks.
- Data Link Layer Protocols: Selective Repeat algorithm, Hamming encoding and decoding.
- Token Passing: Simulated token passing for access control in networks.


## Installation

1. Clone the repository:

```bash
git clone https://github.com/Mayvid0/netSimGo.git
```
2. Navigate to the project directory:
```bash
cd netSimGo
```
3. Run the simulation:
```bash
go run ./cmd/main/main.go
```


## Running Tests
- Use `go test ./...` command to run all tests in the project.
- Individual test cases can be run using `go test -run TestFunctionName`.

## Dependencies
- No external dependencies are required beyond the Go standard library.

## Conclusion
This project provides a comprehensive simulation of network protocols and topologies using Go, demonstrating key concepts in networking and data transmission.

## Future Work

In progress and planned future work includes:

- **TCP/IP Model**: Encompassing the entire TCP/IP model, including layers beyond the physical and data link layers.


## Contributing

Contributions to this project are welcome. Feel free to fork the repository and submit pull requests with your enhancements or bug fixes.

