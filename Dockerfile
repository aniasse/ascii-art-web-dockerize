# Utilise l'image de base golang:alpine
FROM golang:alpine

# Étiquettes pour le mainteneur et la version de l'image
LABEL MAINTAINER="Fatidiop Aboubakdiallo Aniasse" VERSION="1.0"

# Définit le répertoire de travail à /app
WORKDIR /app

# Ajoute tous les fichiers du répertoire courant dans le répertoire /app du conteneur
ADD . /app

# Exécute la commande go build pour construire l'exécutable main
RUN go build -o main .

# Installe les paquets bash et tree dans l'image Alpine
RUN apk update && apk add bash && apk add tree

# Expose le port 8080
EXPOSE 8080

# Définit la commande par défaut à exécuter lorsque le conteneur démarre
CMD ["/app/main"]
