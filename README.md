# NetSimGo

NetSimGo is a networking simulation project developed in Go. It aims to simulate various networking concepts,encompassing the entire **TCP/IP** model, including layers beyond the physical and data link layers.

## Features

- **Device Structs**: Provides basic structures for devices like hubs, bridges, switches, and end devices.
- **Topologies**: Implements different network topologies such as point-to-point and star topologies.
- **Data Link Layer**: Includes functionalities like connecting devices, sending and receiving data.

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


## Usage

- **Device Structs**: Use the provided structs to create and manage networking devices.
- **Topologies**: Explore different network topologies by instantiating and connecting devices.
- **Data Link Layer**: Simulate data transmission and reception between devices.

## Unit Testing

The project includes unit tests to ensure the correctness of data transmission and topology functionalities. You can run the tests using the following command:

```bash
go test
```

## Future Work

In progress and planned future work includes:

- **TCP/IP Model**: Encompassing the entire TCP/IP model, including layers beyond the physical and data link layers.
- **Error Handling**: Implementing error handling mechanisms for robust network communication.
- **Access Control**: Developing access control protocols to manage device permissions and security.
- **Flow Control**: Implementing flow control mechanisms to regulate data transmission rates and prevent congestion.

## Contributing

Contributions to this project are welcome. Feel free to fork the repository and submit pull requests with your enhancements or bug fixes.

