# Myreads

Myreads is a web-based application that lets you manage your bookshelf digitally.
Here you can maintain the status of the books in three categories.

```
1. Wish List: Books which you want to read.
2. Reading List: Books you're currently reading.
3. Finished List: Books you finished reading.
```
# myreads-go-backend

This project is a backend server for the Myreads project.
This project uses HARPERDB as a database, so most of the things are hardcoded.

This server exposes the following API endpoints.
```
1. /api/register: For user registration.
2. /api/login: For use login.
3. /api/user: To fetch the logged-in user information.
4. /api/logout: To log out the user.
5. /api/books/add: To add book the database.
6. /api/books/all: To fetch all the books.
7. /api/books/reading: To fetch books from the reading list.
8. /api/books/wishlist: To fetch books from wishlist.
9. /api/books/finished: To fetch books from the finished list.
10. /api/books/deletebook: To delete book
11. /api/books/updatestatus: To update the status of the books such as reading->finished.
12. /api/static: To fetch the static content such as the image.
```
## Start the server

To run the server, make sure you've docker installed.

### Clone the repo

```sh
git clone https://github.com/pratikjagrut/myreads-go-backend.git
```
### Build container image

```sh
docker build -t myreads-server --build-arg "DB_HOST=yourdbhost" --build-arg "BASIC_AUTH_TOKEN=yourbasicauthtoken" --build-arg "HDB_ADMIN=dbadminusername" --build-arg "PASSWORD=dbadminpassword" --build-arg "PORT=8000" --build-arg "IMAGES_DIR=images" -f  "Dockerfile" .
```

### Run the container

```sh
docker run -v images:/server/images -p 8000:8000 myreads-server true
```

Here the last arg true is for creating schema, database, and table.
If you're running it the first time or on a new database host, please set it true.

-v will mount the /server/images dir to its volume-specific dir.

Hit the APIs at http://localhost:8000 using POSTMAN or something.