docker network inspect infrastructure_infra-network

docker network connect infrastructure_infra-network infrastructure-application-server-1
docker network connect infrastructure_infra-network infrastructure-reverse-proxy-1

docker run -p 8080:8080 --name reverse-proxy --network infrastructure_infra-network reverse-proxy
