# Masterchef Telegram Bot :fire:

Hi!

I'm a **`MasterChefBot`** and ready to help you in the kitchen!

## Get it cooking :pizza:
- Search recipes
- Sign up and collect recipes
    - Add your own ones
    - Favourite others
- Get random recipe



## Development

### Start up :rocket:

Configuration
```
API_KEY:                 "<your-bot-api-key>"
    DEBUG_MODE:          true/false
    DATABASE_CONNECTION: "mongo_connection"
	DATABASE_NAME:       "your_db"  
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
