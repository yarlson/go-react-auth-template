# Go-React Auth Template

A template for building web applications with a Go backend and a React frontend. It includes cookie-based Google authentication with refresh token functionality and uses Tailwind CSS and shadcn/ui for styling. The project integrates Traefik for reverse proxy management and automatic SSL certificates with Let's Encrypt for production deployments. Docker Compose is used for containerization in production. The monorepo structure uses Turborepo for efficient project management and builds.

## 📋 Table of Contents

- [Features](#-features)
- [Getting Started](#-getting-started)
- [Prerequisites](#-prerequisites)
- [Configuration](#-configuration)
- [Usage](#-usage)
- [Production Deployment with Docker Compose](#-production-deployment-with-docker-compose)
- [API Documentation](#-api-documentation)
- [Project Structure](#-project-structure)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

- Go backend using Gin framework
- React frontend with TypeScript
- Tailwind CSS and shadcn/ui for styling
- Cookie-based Google OAuth authentication with refresh token support
- Traefik reverse proxy with automatic Let's Encrypt SSL certificates (production)
- Docker Compose setup for containerized production deployment
- Monorepo structure managed with Turborepo

## 🚀 Getting Started

To use this template:

1. **Use the Template:**

   - Click the "Use this template" button at the top of this repository.
   - Name your new repository and set its visibility.
   - Click "Create repository from template".

2. **Clone the Repository:**

   ```bash
   git clone https://github.com/yourusername/your-repo-name.git
   cd your-repo-name
   ```

## 📋 Prerequisites

- **Node.js** (v20 or later)
- **Go** (v1.23 or later)
- **Google Cloud Platform Account** (for OAuth credentials)
- **Docker & Docker Compose** (for production deployment)
- **Domain Name** (pointing to your server's IP address for production)
- **Valid Email Address** (for Let's Encrypt SSL certificate registration in production)

## 🔧 Configuration

### Backend Configuration

Create a `.env` file in the `services/backend/` directory with the following variables:

```env
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_HOSTED_DOMAIN=
GOOGLE_REDIRECT_URL=http://localhost:5173/callback
JWT_SECRET=your_jwt_secret
HASH_KEY=your_hash_key
BLOCK_KEY=your_block_key
```

- For local development, set `GOOGLE_REDIRECT_URL` to `http://localhost:5173/callback`.
- Replace the placeholder values with your actual configuration.

### Root Configuration (Production Only)

For production deployment, create a `.env` file in the **project root directory** with:

```env
ACME_EMAIL=your_email@example.com
DOMAIN_NAME=yourdomain.com
```

- Replace `your_email@example.com` with your email.
- Replace `yourdomain.com` with your domain.

### DNS Configuration (Production Only)

Ensure your domain name points to your server's IP address for SSL certificate validation.

## 🚀 Usage

### Local Development

1. **Install Dependencies:**

   ```bash
   npm install
   ```

2. **Start Backend Server:**

   Navigate to the backend service directory and run the server:

   ```bash
   cd services/backend
   go run main.go
   ```

   The backend server will start on port `8080` by default.

3. **Start Frontend Development Server:**

   Open a new terminal window, navigate to the frontend service directory, and start the development server:

   ```bash
   cd services/frontend
   npm install
   npm run dev
   ```

   The frontend development server will start on `http://localhost:5173`.

4. **Run Tests:**

   ```bash
   npm run test
   ```

5. **Build the Project:**

   ```bash
   npm run build
   ```

6. **Lint Code:**

   ```bash
   npm run lint
   ```

7. **Format Code:**

   ```bash
   npm run format
   ```

### Notes

- The backend server expects the frontend to be running on `http://localhost:5173`.
- Ensure that the `GOOGLE_REDIRECT_URL` in your backend `.env` file matches the frontend URL.

## 🐳 Production Deployment with Docker Compose

For production deployment, use Docker Compose to build and run the application with Traefik for reverse proxy and automatic SSL.

### Prerequisites

- Ensure you have a domain name pointing to your server's IP address.
- Ports `80` and `443` should be open and accessible.

### Configuration

- Set up the `.env` file in the project root with your production domain and email.
- Update the `GOOGLE_REDIRECT_URL` in `services/backend/.env` to `https://yourdomain.com/callback`.

### Building and Running

1. **Build and Start Containers:**

   ```bash
   docker-compose up --build -d
   ```

2. **Access the Application:**

   - Frontend: `https://yourdomain.com`
   - Backend API: `https://yourdomain.com/api`

### Services Overview

- **Traefik Service:**
  - Reverse proxy and SSL management with Let's Encrypt.
- **Frontend Service:**
  - Serves the React application.
- **Backend Service:**
  - Provides the API endpoints and authentication.

### Environment Variables

Ensure all required environment variables are set in your `.env` files as per the configuration steps.

### Volumes

- `sqlite_data`: Persists the SQLite database.
- `letsencrypt`: Stores SSL certificates.

### Common Commands

- **Stop Containers:**

  ```bash
  docker-compose down
  ```

- **Stop and Remove Volumes:**

  ```bash
  docker-compose down -v
  ```

- **View Logs:**

  ```bash
  docker-compose logs
  ```

- **Rebuild Images:**

  ```bash
  docker-compose build
  ```

### Troubleshooting

- **Check Logs for Errors:**

  ```bash
  docker-compose logs
  ```

- **Verify Environment Variables:**

  Ensure all variables are correctly set.

- **Confirm Domain DNS:**

  Your domain should point to your server's IP.

- **Check Open Ports:**

  Ports 80 and 443 should be accessible.

## 📖 API Documentation

### Authentication Endpoints

- `GET /auth/google`: Initiate Google OAuth login.
- `GET /auth/google/callback`: Handle OAuth callback.
- `POST /auth/refresh`: Refresh user session.
- `GET /auth/logout`: Log out user.

### Protected Endpoints

- `GET /api/user/profile`: Get authenticated user profile.

_Note: Protected endpoints require a valid session cookie._

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
├── .env           # (Production environment variables)
└── README.md
```

## 🤝 Contributing

Contributions are welcome.

1. **Fork the Repository**
2. **Create a Branch**

   ```bash
   git checkout -b feature/your-feature
   ```

3. **Make Changes**
4. **Commit Changes**

   ```bash
   git commit -m 'Add your feature'
   ```

5. **Push to Branch**

   ```bash
   git push origin feature/your-feature
   ```

6. **Open a Pull Request**

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

For support or more information, please open an issue on this repository.
