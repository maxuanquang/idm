# Internet Download Manager

Internet Download Manager is a Go-based application designed to manage and facilitate the downloading of files from the internet efficiently. This project is designed with a scale-up capability in mind, ensuring that it can handle increased workload and demand as needed.

## Features

- Easy-to-use interface for managing download tasks.
- Support for downloading files from various sources.
- Efficient download management with resume functionality.
- Detailed download progress tracking and status updates.
- Cross-platform compatibility.

## Installation

To install Internet Download Manager, follow these steps:

1. Clone the repository:

```bash
git clone https://github.com/maxuanquang/idm.git
```

2. Navigate to the project directory:

```bash
cd idm
```

3. Build project:

```bash
make build
```

4. Start all necessary services:

```bash
make docker-compose-prod-up
```

5. After all services have started up, we can start the project:

```bash
make run
```

## Usage

Once installed, you can use Internet Download Manager to efficiently manage your download tasks. The application provides a user-friendly interface for adding, resuming downloads.

To start using Internet Download Manager:

1. Access the HTTP server at `localhost:8081`.
2. Use the interface to add download tasks by providing the URL of the file you want to download.
3. Monitor the progress of your downloads and manage them as needed.
