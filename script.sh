docker image build -f Dockerfile -t image_test .
docker container run -p 8080:8080 --detach --name testing image_test
docker exec -it testing /bin/bash
docker image ls
docker ps
docker container run image_test