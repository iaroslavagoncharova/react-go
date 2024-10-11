# Language learning API with Go
## Technologies Used
- Go
- MongoDB
- Docker
- Swagger

## Features
- **CRUD Operations**: Create, Read, Update, and Delete functionalities for users, collections, and words.
- **Authentication**: Token-based authentication for user access.
- **Authorization**: Role-based access control to restrict certain actions.
- **Validation**: Input validation to ensure data integrity.

## API Endpoints
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

