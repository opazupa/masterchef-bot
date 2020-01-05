// User collection
db.users.createIndex( { "UserName": 1 }, { unique: true } )

// Recipe collection
db.recipes.createIndex( { Name: 1, Url: 1, UserId: 1 }, { unique: true } )