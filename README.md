# Evermos Store - REST API

Evermos Service is a backend API built using **Golang, Fiber, Cloudinary, and MySQL**. This project provides essential services for store management, product management, transactions, and user authentication. The API follows RESTful principles and includes authentication using JWT tokens.

## Features

- User Authentication (JWT)
- Store Management
- Product Management (with image upload to Cloudinary)
- Category Management
- Transactions & Orders
- Address Management
- MySQL Database (hosted on **Railway**)
- Cloudinary for Image Uploads
- JWT-based authentication
- Swagger API Documentation
- Deployment using **Google Cloud Run**

---

## Getting Started

### Prerequisites

Ensure you have the following installed:

- [Go](https://go.dev/dl/)
- MySQL (or use Railway for cloud hosting)
- Cloudinary account (for image uploads)
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install) (for deployment)

### Installation

1. **Clone the repository:**

   ```sh
   git clone https://github.com/habbazettt/evermos-service-go.git
   cd evermos-service-go
   ```

2. **Install dependencies:**

   ```sh
   go mod tidy
   ```

3. **Setup Environment Variables**

   Create a `.env` file in the root directory and add the following:

   ```env
   # Database URL (For Railway MySQL)
   DB_URL="root:PASSWORD@tcp(HOST:PORT)/DATABASE?charset=utf8mb4&parseTime=True&loc=Local"
   
   # Cloudinary Configuration
   CLOUDINARY_URL="cloudinary://API_KEY:API_SECRET@CLOUD_NAME"
   CLOUDINARY_CLOUD_NAME="your_cloud_name"
   
   # JWT Secret Key
   JWT_SECRET="your_jwt_secret"
   ```

### Running the Application Locally

To start the API locally:

```sh
   go run main.go
```

The API will be available at: `http://localhost:8080`

---

## Deploying MySQL Database using Railway

Instead of setting up MySQL locally, you can deploy it for free using **Railway**:

1. **Create an account & project:**
   - Go to [Railway](https://railway.app/) and sign up.
   - Create a new project.

2. **Deploy MySQL:**
   - Click on "New Service" → Select "MySQL"
   - Wait for the deployment to complete.

3. **Get database credentials:**
   - Go to the **MySQL service** in Railway.
   - Copy the **connection URL** and replace it in your `.env` file under `DB_URL`.

4. **Apply database migrations:**

   ```sh
   go run main.go
   ```

Now your application is connected to **Railway MySQL**.

---

## API Documentation

The API is documented using Swagger. Once the server is running, access the documentation at:

```
http://localhost:8080/swagger/index.html
```

Or if deployed on Cloud Run, replace `localhost` with your Cloud Run URL.

---

## Contributor

Developed by **Hubbal Kholiq Habbaza** as part of **Project-Based Virtual Intern: Backend Developer - Evermos x Rakamin Academy**.

---

## License

MIT License. Feel free to use and contribute!
