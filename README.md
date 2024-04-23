# innovite

Instructions:
1. Have to change the "source: /home/think/aq/temp" in compose.yaml to map the folder in your local to be monitored for file changes

Run:
    docker compose up -d
    docker compose down

Update the <source> folder files and also add new files and observe the docker logs.
Once done open the sqlite3 db file (gorm.db) in container (/go/src/gorm.db) and observe the pathFileName and size contents