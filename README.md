<p align="center">
  <h1 align="center">Faraway Challenge</h1>
  <p align="center">
    <a href="https://github.com/alexandermac/faraway-chal/actions/workflows/ci.yml?query=branch%3Amaster"><img src="https://github.com/alexandermac/faraway-chal/actions/workflows/ci.yml/badge.svg" alt="Build Status"></a>
    <a href="https://goreportcard.com/report/github.com/alexandermac/faraway-chal"><img src="https://goreportcard.com/badge/github.com/alexandermac/faraway-chal" alt="Go Report Card"></a>
    <a href="LICENSE"><img src="https://img.shields.io/github/license/alexandermac/faraway-chal.svg" alt="License"></a>
  </p>
</p>

### Description

An implementation of a TCP server protected from DDoS attacks with the [Proof of Work](https://en.wikipedia.org/wiki/Proof_of_work).
When the clients solves the challenge, the server sends one of the quotes from _Word of wisdom_ book.
Docker files are provided for the server and the client sides.

#### Challenge-response protocol:

I haven't found any RFC specification related to the protocol format and data flow on the Internet. So I've created a simple protocol by myself. This protocol works with binary encoded data in the following format:
- first byte: message id
- rest bytes: message structure in gob format.

There are a few structures, a structure type for each message.

#### Data flow

| Client                                                       | Message         | Server                                                       |
| ------------------------------------------------------------ | --------------- | ------------------------------------------------------------ |
| Sends an initial message without data.                       | => challenge    |                                                              |
|                                                              | <= challenge    | Creates a hash by combining the secret key, client's IP address and nonce.<br />Sends the challenge to the client. |
| Solves the challenge.<br />Sends the solution to the server. | => solution     |                                                              |
|                                                              | <= grant        | Validates the client solution.<br />If it's valid, grants the access (returns a poem), otherwise returns an error. |
| Prints the poem to console and exits.                        |                 |                                                              |

#### POW algorithm
In my implementation I've used _Hashcash_ PoW algorithm. It was chosen because of its simplicity and prevalence. The difficulty is hardcoded and equals 10.

### Usage
```sh
# First, run the server
$ scripts/server/run.sh

# Then, run the client
$ scripts/client/run.sh

# Or run applications as Docker Containers
$ sudo docker compose up
```

### License
Licensed under the MIT license.

### Author
Alexander Mac
