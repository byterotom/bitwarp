# Bitwarp: A P2P file sharing system
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

Bitwarp is a peer-to-peer file sharing system written in Go.

It uses **gRPC** for file transfers. A **tracker** service (with multiple instances) uses **Redis** to maintain metadata of chunk holders. These trackers are load balanced using **HAProxy** and kept in sync via **RabbitMQ**.

### How It Works

* A peer (the seeder) seeds a file into the network.
* Other peers contact the tracker to get the list of chunk holders.
* Instead of downloading all chunks from the seeder, peers fetch specific chunks from other peers.
* Peers then upload their chunks to others, enabling distributed sharing.

This leads to an efficient swarm-like mechanism, where chunks of a file are exchanged among nodes.


## Requirements

* Go 1.24.2
* Latest Docker and Docker Compose
* `make`

Install `make` if not available:

```bash
sudo apt install make
```


## Getting Started

### 1. Build the System

```bash
make build
```

### 2. Prepare the File

Place the actual file to be shared in the `storage/downloads` directory.

Then generate a warp file:

```bash
go run cmd/warpgen/main.go storage/downloads/<your_file>
```

This will create a corresponding warp JSON file in `storage/warp`.

### 3. Configure the File

Update `docker-compose.yaml` with the generated file name.

### 4. Run the System

Spin up trackers and nodes:

```bash
docker compose up --scale tracker=2 --scale node=5
```

You’ll now observe chunks being transferred and shared among the nodes.


## Contribution

Open-source contributions are welcome. Feel free to fork, suggest changes, or file issues.

## License

This project is licensed under the MIT License.