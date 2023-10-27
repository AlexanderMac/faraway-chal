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

I haven't found any RFC document related to the protocol format and flow on the Internet. So I've created a simple protocol by myself. This protocol operates with the binary data in the format `code:json`, where the `code` is the name of the message and the `json` is the data. There are three JSON schemes, one schema for each message.

#### Data flow

| Client                                                       | Message         | Server                                                       |
| ------------------------------------------------------------ | --------------- | ------------------------------------------------------------ |
| Sends an initial message without data.                       | => init         |                                                              |
|                                                              | <= challenge    | Creates a hash by combining the secret key and the client's IP address.<br />Randomly selects POW algorithm and difficulty.<br />Sends the hash and POW details to the client. |
| Solves the challenge.<br />Sends the generated hash to the server. | => solution     |                                                              |
|                                                              | <= grant-access | Validates the client hash.<br />If it's valid, grants the access (returns a poem). |
| Prints the poem and exits.                                   |                 |                                                              |

#### POW algorithm
I've implemented two POW algorithms. They were chosen because TODO.

### Usage
```sh
# First, run the server
scripts/run-server.sh

# Then, run the client
scripts/run-client.sh

# That's it
```

### License
Licensed under the MIT license.

### Author
Alexander Mac
