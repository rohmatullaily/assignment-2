## Step 1: Install Dependencies
Before you can run the code, make sure you have the necessary Go packages installed:

go get github.com/antonlindstrom/pgstore
go get golang.org/x/crypto/scrypt
go get gorm.io/gorm

## Step 2: Set Up PostgreSQL Database
Create a PostgreSQL database with the name "postgres" and a schema named "training." Make sure you have PostgreSQL running locally on the default port (5432). You can adjust the database connection parameters in the code if needed.

## Step 3: Create Required HTML Templates
Create the following HTML templates in a "layout" folder in the same directory as your Go code:

layout/login.html for the login page.
layout/register.html for the registration page.
layout/home.html for the home page.
These templates should match the placeholders in your code.

## Step 4: Run the Application
You can run the Go application by executing the main Go file in your terminal:
go run main.go

## Step 5: Access the Application
Open a web browser and access the following URLs:

http://localhost:5555/ - This should take you to the login page.
<img width="1680" alt="Screen Shot 2023-09-09 at 02 47 20" src="https://github.com/rohmatullaily/assignment-2/assets/31227217/e353e03e-047d-4d7f-b6c3-6a23a3f34e69">

http://localhost:5555/register - This is the registration page.
<img width="1679" alt="Screen Shot 2023-09-09 at 02 26 23" src="https://github.com/rohmatullaily/assignment-2/assets/31227217/fd27f026-0e3d-41f0-b475-b0a418f99f79">

http://localhost:5555/logout - This logs out the user and redirects to the login page.
<img width="1680" alt="Screen Shot 2023-09-09 at 02 27 23" src="https://github.com/rohmatullaily/assignment-2/assets/31227217/162d2a45-8aa5-44c6-9efd-33a92f03c7a5">

You can then interact with the application by registering new users and logging in.
