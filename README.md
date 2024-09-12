# Go-React Auth Template

This repository is a GitHub template for building web applications with a Go backend and React frontend. It includes cookie-based Google authentication with refresh token functionality and uses Tailwind CSS and shadcn/ui for styling. The template aims to offer a starting point for developers creating web applications with modern authentication mechanisms and beautiful, accessible UI components.

## 📋 Table of Contents

- [Features](#-features)
- [Getting Started](#-getting-started)
- [Prerequisites](#-prerequisites)
- [Configuration](#-configuration)
- [Usage](#-usage)
- [Docker Compose](#-docker-compose)
- [API Documentation](#-api-documentation)
- [Project Structure](#-project-structure)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

- Go backend using Gin web framework
- React frontend with TypeScript
- Tailwind CSS for utility-first styling
- shadcn/ui for beautiful, accessible UI components
- Google OAuth integration
- Secure cookie-based authentication
- Refresh token mechanism for persistent sessions
- SQLite database for user and token storage
- CORS configuration for local development
- Monorepo structure using Turborepo
- Docker Compose setup for easy deployment

## 🚀 Getting Started

To use this template:

1. Click the "Use this template" button at the top of this repository.
2. Choose a name for your new repository and select its visibility.
3. Click "Create repository from template".

After creating your repository, clone it locally:

```bash
git clone https://github.com/yourusername/your-repo-name.git
cd your-repo-name
```

## 📋 Prerequisites

- Node.js (v14 or later)
- Go (v1.16 or later)
- Google Cloud Platform account (for OAuth credentials)
- Docker and Docker Compose (for containerized deployment)

## 🔧 Configuration

The backend service requires several environment variables to be set. Create a `.env` file in the `services/backend/` directory with the following variables:

```
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_HOSTED_DOMAIN=
GOOGLE_REDIRECT_URL=http://localhost:5173/callback
JWT_SECRET=your_jwt_secret
HASH_KEY=your_hash_key
BLOCK_KEY=your_block_key
```

Replace the placeholder values with your actual configuration. Never commit this file to version control.

For security reasons:

- Use strong, unique values for JWT_SECRET, HASH_KEY, and BLOCK_KEY
- Keep your Google client credentials confidential
- In production, use a secure method to manage environment variables

## 🚀 Usage

1. Install dependencies:

   ```bash
   npm install
   ```

2. Start the development server:

   ```bash
   npm run dev
   ```

3. Run tests:

   ```bash
   npm run test
   ```

4. Build the Docker image:
   ```bash
   npm run docker:build
   ```

## 🐳 Docker Compose

This project includes a `docker-compose.yml` file for easy deployment of the entire stack.

### Building and Running

1. Make sure you're in the project root directory.

2. Build and start the containers:

   ```bash
   docker-compose up --build
   ```

   This command will build the images for both the frontend and backend services and start the containers.

3. Access the application:
   - Frontend: http://localhost
   - Backend: http://localhost:8080

### Configuration

The `docker-compose.yml` file defines two services:

#### Frontend Service

- Built from `services/frontend/Dockerfile`
- Uses Node.js version 20.11.1
- Exposed on port 80
- Depends on the backend service
- Has a healthcheck that curls localhost every 30 seconds
- Connects to the backend using the `BACKEND_URL` environment variable

#### Backend Service

- Built from `services/backend/Dockerfile`
- Uses Go version 1.23
- Exposed on port 8080
- Uses a volume for SQLite data persistence
- Has a healthcheck that curls the `/health` endpoint every 30 seconds
- Uses environment variables from `services/backend/.env` file

### Volumes

- `sqlite_data`: Persists the SQLite database file

### Networks

- `app_network`: A bridge network for communication between services

### Stopping the Stack

To stop the running containers:

```bash
docker-compose down
```

To stop the containers and remove the volumes:

```bash
docker-compose down -v
```

### Viewing Logs

To view logs from all services:

```bash
docker-compose logs
```

To follow logs from a specific service:

```bash
docker-compose logs -f frontend
```

or

```bash
docker-compose logs -f backend
```

### Rebuilding

If you make changes to the code, rebuild the images:

```bash
docker-compose build
```

Then restart the services:

```bash
docker-compose up
```

### Troubleshooting

1. If services fail to start, check the logs for error messages:

   ```bash
   docker-compose logs
   ```

2. Ensure all required environment variables are set in the `services/backend/.env` file.

3. If changes aren't reflecting, try rebuilding the images:

   ```bash
   docker-compose build --no-cache
   ```

4. Check if all services are running:

   ```bash
   docker-compose ps
   ```

## 📖 API Documentation

### Authentication Endpoints

- `GET /auth/google`: Initiates Google OAuth login
- `GET /auth/google/callback`: Handles Google OAuth callback
- `POST /auth/refresh`: Refreshes the user's session
- `GET /auth/logout`: Logs out the user

### Protected Endpoints

- `GET /api/user/profile`: Retrieves the authenticated user's profile

All protected endpoints require a valid session cookie.

## 🗂 Project Structure

```
.
├── services/
│   ├── backend/
│   │   ├── auth/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── main.go
│   │   └── ...
│   └── frontend/
│       ├── src/
│       ├── public/
│       └── ...
├── package.json
├── turbo.json
├── docker-compose.yml
└── README.md
```

## 🤝 Contributing

Contributions to improve this template are welcome. Please follow these steps:

1. Fork the repository
2. Create a new branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

For more information or support, please open an issue on this repository.
