
# Without CADT MySQL Mirror Database:

# CONFIG_FILE=/root/.chia/testnet/cadt/v1/config.yaml

# if [ -f "$CONFIG_FILE" ]; then
#     sed -i 's/CHIA_NETWORK:.*/CHIA_NETWORK: testnet/' $CONFIG_FILE
#     sed -i 's/USE_DEVELOPMENT_MODE:.*/USE_DEVELOPMENT_MODE: true/' $CONFIG_FILE
# fi

#############################################################################################

#With CADT MySQL Mirror Database:

CONFIG_FILE=/root/.chia/testnet/cadt/v1/config.yaml

# External MySQL Database

# DB_USERNAME=${DB_USERNAME: <MySQL Database Username>}
# DB_PASSWORD=${DB_PASSWORD: <MySQL Database Password>}
# DB_NAME=${DB_NAME: <MySQL Database Name>}
# DB_HOST=${DB_HOST: <MySQL Database Host>}

# The Docker Compose MySQL Database(s)

DB_USERNAME=${DB_USERNAME}
DB_PASSWORD=${DB_PASSWORD}
DB_NAME=${DB_NAME}
DB_HOST=${DB_HOST}

if [ -f "$CONFIG_FILE" ]; then
    sed -i 's/CHIA_NETWORK:.*/CHIA_NETWORK: testnet/' $CONFIG_FILE
    sed -i 's/USE_DEVELOPMENT_MODE:.*/USE_DEVELOPMENT_MODE: true/' $CONFIG_FILE

    if ! grep -q "MIRROR_DB" $CONFIG_FILE; then
        cat <<EOL >> $CONFIG_FILE
MIRROR_DB:
  DB_USERNAME: $DB_USERNAME
  DB_PASSWORD: $DB_PASSWORD
  DB_NAME: $DB_NAME
  DB_HOST: $DB_HOST
EOL
    else
        sed -i "s|DB_USERNAME:.*|DB_USERNAME: $DB_USERNAME|" $CONFIG_FILE
        sed -i "s|DB_PASSWORD:.*|DB_PASSWORD: $DB_PASSWORD|" $CONFIG_FILE
        sed -i "s|DB_NAME:.*|DB_NAME: $DB_NAME|" $CONFIG_FILE
        sed -i "s|DB_HOST:.*|DB_HOST: $DB_HOST|" $CONFIG_FILE
    fi
fi
