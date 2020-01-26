// User collection
db.users.createIndex( { "UserName": 1,  }, { unique: true } )
db.users.createIndex( { "TelegramID": 1,  }, { unique: true } )

// Recipe collection
db.recipes.createIndex( { "Name": 1, "URL": 1, "UserID": 1 }, { unique: true } )

// SelectedRecipe collection
db.selectedrecipes.createIndex( {"ChatID" : 1, "UserID": 1 }, { unique: true } )