if [ -z "${KEY}" ]; then
    echo "KEY is undefined"
    exit
fi
docker stack deploy -c docker-compose.yml eos
