# Construit une image Docker en utilisant le Dockerfile spécifié (-f) et lui attribue le tag (-t) "image_test"
docker image build -f Dockerfile -t image_test .
echo "                              --                                    "
echo "                              --                                    "
# Exécute un conteneur Docker à partir de l'image "image_test", en mappant le port 8080 de l'hôte au port 8080 du conteneur,
# en mode détaché (--detach) et en donnant un nom (--name) "testing" au conteneur
docker container run -p 8080:8080 --detach --name testing image_test
echo "                              --                                    "
echo "                              --                                    "
# Exécute une commande interactive dans le conteneur "testing" en démarrant un shell bash (/bin/bash)
docker exec -it testing /bin/bash
echo "                              --                                    "
echo "                              --                                    "
# Liste les images Docker présentes sur le système
docker image ls
echo "                              --                                    "
echo "                              --                                    "
# Liste les conteneurs Docker en cours d'exécution
docker ps
echo "                              --                                    "
echo "                              --                                    "
# Exécute un conteneur Docker à partir de l'image "image_test" (par défaut en mode interactif)
docker container run image_test
