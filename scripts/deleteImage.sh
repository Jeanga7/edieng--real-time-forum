# Stop all containers
docker rm -vf $(docker ps -aq)
# Delete all images
docker rmi -f $(docker images -aq)
