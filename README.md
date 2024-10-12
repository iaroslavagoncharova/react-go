# Language learning API with Go

## Technologies used  
- **Go**: Backend development  
- **MongoDB**: NoSQL database  
- **Docker**: Containerization for easy deployment  
- **Swagger/OpenAPI**: API documentation 

## Features
- **CRUD Operations**: Create, Read, Update, and Delete functionalities for users, collections, and words.
- **Authentication**: Token-based authentication for user access.
- **Authorization**: Role-based access control to restrict certain actions.
- **Validation**: Input validation to ensure data integrity.

## API documentation  
You can access the API documentation through [Swagger](http://10.120.33.51:8080/swagger/index.html).

## API endpoints
| Method | URL                                   | Description                    |
|--------|---------------------------------------|--------------------------------|
| GET    | `/api/users`                          | Retrieve all users             |
| GET    | `/api/users/:id`                      | Retrieve a user by ID          |
| POST   | `/api/users`                          | Create a new user              |
| PATCH  | `/api/users/:id`                      | Update an existing user        |
| DELETE | `/api/users/:id`                      | Delete a user                  |
| GET    | `/api/collections`                    | Retrieve all collections       |
| GET    | `/api/collections/:id`                | Retrieve a collection by ID    |
| POST   | `/api/collections`                    | Create a new collection        |
| PATCH  | `/api/collections/:id`                | Update an existing collection  |
| DELETE | `/api/collections/:id`                | Delete a collection            |
| GET    | `/api/collections/:collectionId/words`| Retrieve words in a collection |
| POST   | `/api/collections/:collectionId/words`| Create a word                  |
| GET    | `/api/words/{id}`                     | Retrieve a word by ID          |
| PATCH  | `/api/words/{id}`                     | Update an existing word        |
| DELETE | `/api/words/{id}`                     | Delete a word                  |

## Deployment
The API has been deployed to a RockyLinux virtual machine using Docker. Below is a summary of the deployment process:

1. Built a Docker image:

``` bash
docker build -t my_app_image .
```
2. Ran the Docker container:

``` bash
sudo docker run -d --restart unless-stopped -p 8080:8080 --name my_app_container my_app_image
```

3. Access the deployed API:
Open the deployed [link](http://10.120.33.51:8080/) in a browser or use Postman.

## Authentication
The API uses JWT for user authentication.
After logging in, you’ll receive a token, which must be passed in the Authorization header for protected routes.
```bash
Authorization: Bearer <your-token>
```

## List of Go Packages Used
Here’s a list of key Go packages used in this project:

- github.com/gofiber/fiber/v2: Web framework for building APIs.
- go.mongodb.org/mongo-driver: MongoDB driver for Go.
- github.com/golang-jwt/jwt/v5: JWT for authentication.
- golang.org/x/crypto/bcrypt: Hashing for user passwords.
- github.com/swaggo/swag: Swagger documentation generator.

## Testing the API with Postman
A Postman collection with sample API requests has been provided in OMA. You can import it into Postman to quickly test the API.
