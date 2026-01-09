# 📚 Rent and Read API

A robust Go-based backend for a physical book rental system. This project features a virtual wallet system, real-time stock management, and Xendit payment gateway integration.

## 🚀 Features

    Book Management: Browse books with category/genre filtering.

    Virtual Wallet: Top up your balance via Xendit and use it for one-click rentals.

    Rental System: Automated stock checking and atomic transactions.

    Payment Integration: Secure payment handling using Xendit Invoices and Webhooks.

    Layered Architecture: Organized into Handler, Service, and Repository layers for maintainability.

## 🏗️ System Architecture

The project follows a Clean Architecture approach to separate business logic from technical implementation.

### 🛠️ Tech Stack

    Language: Go (Golang)

    Web Framework: Echo

    ORM: GORM

    Database: PostgreSQL

    Payment Gateway: Xendit

### 📋 Database Schema

The system uses a relational PostgreSQL schema to ensure data integrity, especially for financial transactions.

![ERD Image](/docs/erd/mermaid-diagram-2026-01-08-101502.png)

🚦 API Endpoints (Quick Reference)
Books

    GET /books - List all books with their genre names.

    GET /books/:id - Get specific book details.

Wallet & Payments

    POST /wallet/topup - Request a balance top-up (returns Xendit URL).

    POST /webhooks/xendit - Callback for Xendit to confirm payments.

Rentals

    POST /rentals - Rent a book using wallet balance.
    