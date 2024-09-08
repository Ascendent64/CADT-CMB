
MYSQL_CNF="/etc/mysql/mysql.conf.d/mysqld.cnf"

echo "Waiting for MySQL to initialize..."
until mysqladmin ping -h"localhost" --silent; do
    sleep 1
done
echo "MySQL initialized."

if [ -f "$MYSQL_CNF" ]; then
  cp $MYSQL_CNF "${MYSQL_CNF}.backup"
  echo "Backup of MySQL configuration created."
fi

if grep -q "sql_mode" $MYSQL_CNF; then
  sed -i 's/^sql_mode=.*/sql_mode="NO_ENGINE_SUBSTITUTION"/' $MYSQL_CNF
else
  echo 'sql_mode="NO_ENGINE_SUBSTITUTION"' >> $MYSQL_CNF
fi

service mysql restart
echo "MySQL configuration updated and MySQL restarted."
