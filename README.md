# Golang, GORM, & Fiber: Ramen

In this comprehensive guide, you'll learn how to implement JWT (JSON Web Token) authentication in a Golang application using GORM and the Fiber web framework. The REST API will be powered by a high-performance Fiber HTTP server, offering endpoints dedicated to secure user authentication, and persist data in a PostgreSQL database.

![Golang, GORM, & Fiber: JWT Authentication](https://ramen.mn/wp-content/themes/aunm/assets/img/logo-light.svg)

## Topics Covered

- Run the Golang & Fiber JWT Auth Projec
- Setup the Golang Project
- Setup PostgreSQL and pgAdmin with Docker
- Create the GORM Model
- Database Migration with GORM
  - Load the Environment Variables with Viper
  - Create the Database Pool with GORM
  - Migrate the GORM Model to the Database
- Create the JWT Authentication Controllers
  - SignUp User Fiber Context Handler
  - SignIn User Fiber Context Handler
  - Logout User Fiber Context Handler
- Get the Authenticated User
- Create the JWT Middleware Guard
- Register the Routes and Add CORS
- Testing the JWT Authentication Flow
  - Register a New Account
  - Log into the Account
  - Access Protected Routes
  - Logout from the API

Read the entire article here: [https://codevoweb.com/golang-gorm-fiber-jwt-authentication/](https://codevoweb.com/golang-gorm-fiber-jwt-authentication/)
