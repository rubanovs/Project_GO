**Building the Code**:
1. Open your terminal or command prompt.
2. Navigate to the directory where your Go application is located.
3. Run the following command to build your application:
   ```
   go build
   ```

**Running Tests**:
1. In the same directory as your application, execute the following command:
   ```
   go test
   ```

This will run all the tests in your application. If you have separate test packages, you can specify them instead of `go test`.

Please ensure that you have Go installed and that your application is set up for building and testing. If you have any questions, feel free to ask. Good luck with your development! 🚀

## Context

We are developing a REST API in **Golang** that allows users to create, retrieve, and stop commands. The commands are bash scripts, and the API should support long-running commands by saving their output to a database and displaying the output when retrieving a specific command.

## Decision

We will implement the API using the following architecture and design decisions:

1. **Language and Libraries**:
   - **Language**: Golang
   - **Database**: PostgreSQL
   - **Router**: Gorilla Mux

2. **Database Schema**:
   - We will create a table named `commands` with the following columns:
     - `id` (serial primary key)
     - `content` (text): The bash script content
     - `output` (text): The output of the executed command

3. **API Endpoints**:
   - `POST /commands`: Create a new command. The request body should contain the bash script content.
   - `GET /commands`: Retrieve a list of all commands.
   - `GET /commands/{id}`: Retrieve a specific command by its ID.
   - `POST /commands/{id}/stop`: Stop the execution of a specific command (not implemented in this ADR).

4. **Command Execution**:
   - When creating a new command, we will execute it using `exec.Command("bash", "-c", content)`.
   - The output of the command will be saved to the `output` column in the database.

5. **Testing**:
   - We will write unit tests for each API endpoint and the command execution logic.
   - Test data will be stored in a separate test database (`commands_test`).

6. **Deployment and CI/CD** (optional):
   - We can set up a GitLab CI/CD pipeline to build Docker images, create deb/rpm packages, and deploy the application.

## Consequences

- **Scalability**: The API can handle multiple concurrent requests for command execution and retrieval.
- **Data Consistency**: The output of long-running commands will be saved to the database, ensuring data consistency.
- **Testability**: The API endpoints and command execution logic will be thoroughly tested.
- **Complexity**: The additional functionality for stopping commands and handling long-running commands increases complexity.

---

This ADR outlines the architectural decisions for our Command Execution and Storage API. It provides a clear direction for implementation and testing. If any further decisions or changes arise during development, we will document them in additional ADRs. 🚀