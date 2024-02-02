
# Testing Golang Project 

Welcome to the **Testing Golang Project Template!** This project provides a structured and comprehensive template for testing **APIs**, covering basic **CRUD operations** for user management, including features like login and register.

# Table of Contents
- Introduction
- Features
- Getting Started
- Testing

## Introduction

This Golang project template is designed to showcase best practices for testing in Go applications. It focuses on creating a simple API for user management, including registration, login, and basic CRUD operations. The goal is to provide a clean and well-organized foundation for building scalable and maintainable applications.


## Features

- **User Management**: Implement user registration, login, and basic CRUD operations.
- **API Testing**: Comprehensive testing suite covering API endpoints.
- **Environment Configuration**: Utilize environment variables for configuration.
- **Database Interaction**: Interact with a MySQL database for user data storage.
- **Structured Logging**: Employ structured logging for better traceability.
- **Dependency Management**: Use Go modules for efficient dependency management.
- **Consistent Coding Style**: Follow a consistent coding style for better code readability.
- **Documentation**: Well-documented code and a README for easy understanding.

## Getting Started

To get started with this project, follow these steps:

1. Clone the repository:
```git
git clone this repo link
cd your-repo

```
2. Set up your environment variables:

Create a **.env** file based on **.env.example** and fill in the required configuration.

3. Install dependencies:

```go
go mod download
```
4. change your database configuration in **db.go** and setup tabels in package **migrate**

5. run migrate to migrate your tabel to database 
 ```make
 make migrate
 ```
 6. run projects
 ```make
 make or make run
 ```
## Running Tests

To run tests, run the following command

```bash
  make test
```

