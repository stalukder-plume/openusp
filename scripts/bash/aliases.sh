alias dc="docker compose -f $(pwd)/deployments/docker-compose.yaml" 
alias dclocal="docker compose -f $(pwd)/deployments/docker-compose_local.yaml" 
alias cli="docker run --env API_SERVER_ADDR=http://localhost:8081 --env API_SERVER_AUTH_NAME=admin --env API_SERVER_AUTH_PASSWD=admin --network=openusp -it --rm n4networks/openusp-cli"
