# empty
Empty repo for github classroom starter

# Login to get JWT token
curl -X POST -d 'username=admin' -d 'password=password' localhost:3000/login

# endpoints without JWT
- get book list : GET /books
- get book detail : GET /books/{book_id}

# endpoints with JWT as bearer token
- create book : POST /admin/books 
body
JSON : {
    "isbn" : "111111",
    "title" : "malin kundang",
    "author" : "wiro",
    "price" : 2000
}
- update book : PUT /admin/books/{book_id}
body
JSON : {
    "isbn" : "111111",
    "title" : "malin kundang edited",
    "author" : "wiro e",
    "price" : 2200
}   
- delete book : DELETE /admin/books/{book_id}