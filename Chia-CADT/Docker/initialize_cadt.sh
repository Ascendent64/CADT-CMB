export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

cd /cadt
echo 'In CADT directory'
echo 'PATH is: $PATH'
. $NVM_DIR/nvm.sh
echo 'NVM sourced'
nvm use $NODE_VERSION
echo 'Using Node version $NODE_VERSION'
npm run start &

CADT_PID=$!
echo 'CADT started with PID' $CADT_PID

sleep 30

echo "Terminating CADT process with PID $CADT_PID"
kill $CADT_PID
sleep 5

if ps -p $CADT_PID > /dev/null; then
  echo 'CADT process did not terminate. Forcing termination...'
  pkill -f 'node'
  sleep 5
fi

echo 'Modifying CADT configuration'
/modify_config.sh

echo 'Checking if port 31310 is still in use'
while lsof -i:31310; do
  echo 'Waiting for port 31310 to be free...'
  fuser -k 31310/tcp
  sleep 1
done

echo 'Port 31310 is free. Restarting CADT...'

. $NVM_DIR/nvm.sh
echo 'NVM sourced'
nvm use $NODE_VERSION
echo 'Using Node version $NODE_VERSION'
npm run start &
echo 'CADT restarted with new PID' $!
tail -f /dev/null
