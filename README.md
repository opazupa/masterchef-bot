# Masterchef Telegram Bot :fire:

Hi!

I'm a **`MasterChefBot`** and ready to help you in the kitchen!

## Get it cooking :pizza:
- Search recipes
- Sign up and collect recipes
- Shuffle a ßrandom recipe of yours



## Development

### Start up :rocket:

Configuration
```
    API_KEY:            "<your-bot-api-key>"
    DEBUG_MODE:         true/false
    DatabaseConnection: "mongo_connection"
	DatabaseName:       "your_db"  
```

### Running locally
`Hot reload`
```
docker-compose up
docker-compose down
```
`Production`
```
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up
docker-compose -f docker-compose.yml -f docker-compose.prod.yml down
```
