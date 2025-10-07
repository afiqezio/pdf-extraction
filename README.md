# PDF Extraction Platform (Mac user)

A full-stack application for PDF extraction and data management with a Go backend and Nuxt.js frontend.

## 🏗️ Architecture

- **Backend**: Go with Echo framework, GORM ORM, PostgreSQL database
- **Frontend**: Nuxt.js 3 with Vue 3, Tailwind CSS
- **Database**: PostgreSQL with UUID support

## 📋 Prerequisites

Before you begin, ensure you have the following installed on your machine:

- **Go** (version 1.23.2 or higher)
- **Node.js** (version 18 or higher)
- **PostgreSQL** (version 12 or higher)
- **Git**

### Installing Prerequisites

#### macOS (using Homebrew)
```bash
# Install Go
brew install go

# Install Node.js
brew install node

# Install PostgreSQL
brew install postgresql
brew services start postgresql
```

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/afiqezio/pdf-extraction
cd pdf-extraction
```

### 2. Set Up Environment Variables
```bash
cp .env.example .env
# Edit .env with your configuration
```

### 3. Database Setup

Connect to PostgreSQL using your system user (no password required for local setup):

```bash
# Connect to PostgreSQL as your system user
psql postgres
```

Run the following SQL commands in the PostgreSQL prompt:

```sql
-- Create database
CREATE DATABASE pdf_extraction;

-- Use default postgres user, grant privileges:
GRANT ALL PRIVILEGES ON DATABASE pdf_extraction TO postgres;

-- Exit PostgreSQL
\q
```

### 4. Start Development Environment
```bash
# Terminal 1: Start backend
cd backend
go mod download
go run cmd/server/main.go

# Terminal 2: Start frontend
cd frontend
npm install
npm run dev
```
