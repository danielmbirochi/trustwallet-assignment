# Tx Parser
Implement Ethereum blockchain parser that will allow to query transactions for subscribed addresses.

## Project Structure

The project is organized as follows:

- **internal**: This directory contains three packages.
  - **state**: This package is used to define the application state. It includes the data structures and methods necessary for maintaining and manipulating the application state during its lifecycle.
  - **txparser**: This package is responsible for implementing the core functionalities of the application. It interacts with blockchain for extracting and parsing on-chain data.
  - **service**: The file `service.go` defines the `Parser` interface, which can be used by API implementations and consumed by client applications. It also defines the `Transaction` domain type, which represents the transaction object in the application context.
  
In the root directory, you'll find a `Makefile`, which includes commands for building, running, and testing the application.

## Testing

The project includes automated tests that can be run with the `make test` command. These tests ensure that all functionalities are working as expected.

## Running the Application

To run the application, use the command: `./txparser -block=<block number>`. The `<block number>` should be replaced with the actual number of the initial block to be scanned. After the initial block, the application will continue to scan subsequent blocks.

## Future Improvements

While the current application serves its primary purpose, the following improvements could enrich the application:

- **Error handling**: By far the most important improvement is error handling. The current state of this app lacks proper error handling for targeting production environment.
- **EIP-55 support**: The current application is not compliant with EIP-55. Cloudflare API return addresses in lower case mode. Subscribed addresses will be persisted in lower case mode.
- **Logging System**: Implement a robust logging system to trace the application's operations and potential issues.
- **Application Metrics**: Incorporate a monitoring tool to gather important metrics like memory footprint, block scanning time, request count, error rates, etc., which could provide useful insights into the application's behavior.
- **KeyValueStorer Database Engine**: Currently, a mock key-value store is used. However, implementing a proper database engine will provide reliable, persistent data storage, improving data handling capabilities.


## License

This project is licensed under the MIT License. Please refer to the `LICENSE` file for more details.


