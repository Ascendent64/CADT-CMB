export PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

check_wallet_sync() {
  cd /chia-blockchain
  source ./activate
  sync_status=$(chia wallet show | grep 'Sync status:' | awk '{print $3}')
  deactivate
  if [ "$sync_status" = "Synced" ]; then
    return 0
  else
    return 1
  fi
}

KEYS_DIR="/keys"

echo "Checking if KEYS_DIR exists..."
if [ ! -d "$KEYS_DIR" ]; then
  echo "Creating KEYS_DIR..."
  mkdir -p $KEYS_DIR
else
  echo "KEYS_DIR already exists."
fi

add_keys() {
  for key in "$KEYS_DIR"/*; do
    if [ -f "$key" ]; then
      echo "Adding key from file $key..."
      mnemonic=$(cat "$key")
      chia keys add <<EOF
$mnemonic
Key1
EOF
    else
      echo "$key is not a regular file."
    fi
  done
}

generate_and_store_key() {
  new_key_file="$KEYS_DIR/mnemonic.key"
  echo "Generating a new key..."
  chia keys generate -l "Key1"
  echo "Storing the mnemonic seed in $new_key_file..."
  chia keys show --show-mnemonic-seed | awk '/Mnemonic seed/{getline; print}' > "$new_key_file"
}

cd /chia-blockchain
source ./activate
chia init

echo "Checking for existing keys in KEYS_DIR..."
if [ "$(ls -A $KEYS_DIR)" ]; then
  echo "Existing keys found. Adding keys..."
  add_keys
else
  echo "No keys found. Generating a new key..."
  generate_and_store_key
  add_keys
fi

chia init
chia init --fix-ssl-permissions

chia configure --testnet true && \
export CHIA_ROOT=/root/.chia/testnet && \
chia init && \
chia start wallet -r && \
echo 'Chia wallet started' && \
chia start data && \
echo 'Chia data started' && \
chia start data_layer_http && \
echo 'Chia data_layer_http started'

deactivate

while ! check_wallet_sync; do
  echo "Wallet is not synced yet. Waiting for 60 seconds..."
  sleep 60
done

echo "Wallet is synced. Starting CADT..."
