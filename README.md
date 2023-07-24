# Tx Parser
Implement Ethereum blockchain parser that will allow to query transactions for subscribed addresses.

## Project Structure

The project is organized as follows:

- **internal**: This directory contains three packages.
  - **state**: This package is used to define the application state. It includes the data structures and methods necessary for maintaining and manipulating the application state during its lifecycle.
  - **txparser**: This package is responsible for implementing the core functionalities of the application. It interacts with blockchain for extracting and parsing on-chain data.
  - **service**: The file `service.go` defines the `Parser` interface, which can be used by API implementations and consumed by client applications. It also defines the `Transaction` domain type, which represents the transaction object in the application context.
  
In the root directory, you'll find a `Makefile`, which includes commands for building, running, and testing the application.

## License

This project is licensed under the MIT License. Please refer to the `LICENSE` file for more details.
