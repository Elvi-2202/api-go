Ce projet consiste a developper une API REST en Go permettant de gerer une ressource Book avec un CRUD complet. L objectif est de pratiquer l acces aux donnees en SQL brut via database/sql, sans aucun ORM.

ETAPE 0 INITIALISATION DU PROJET

Cette phase prepare l environnement de developpement et l architecture logicielle.
1 Initialisation du module : go mod init [github.com/Elvi-2202/book-api]
2 Creation de l arborescence : cmd/api, internal/model, internal/repository, internal/handler, migrations.
3 Installation des dependances :
go get [github.com/go-chi/chi/v5]
go get [github.com/jackc/pgx/v5/stdlib]

ETAPE 1 POSTGRESQL ET MIGRATIONS

Mise en place de la base de donnees et du schema.
Fichier docker-compose.yml pour lancer PostgreSQL localement

Commande de lancement : docker-compose up -d.

Migrations SQL pour la table books

ETAPE 2 CONNEXION ET SHUTDOWN PROPRE

Configuration de la connexion base de donnees via variables d environnement.
1 Ouverture de la connexion avec le driver pgx/stdlib.
2 Verification de la connectivite via PingContext au demarrage.
3 Ajout d un shutdown propre pour fermer la DB et le serveur HTTP lors de l arret.

ETAPE 3 MODELES ET VALIDATION

Definition des structures de donnees dans internal/model/book.go.
L API valide que le titre et l auteur ne sont pas vides.
Format d echange utilise : JSON avec nomenclature snake_case.

ETAPE 4 REPOSITORY SQL SANS ORM

Implementation de l interface BookRepository utilisant des requetes SQL parametrees pour la securite.
Gestion de sql.ErrNoRows pour renvoyer des erreurs 404.

ETAPE 5 EXECUTION DU MAIN ET HANDLERS

Pour lancer l application:
1 Se placer a la racine du projet.
2 Executer la commande : go run cmd/api/main.go.

L API utilise le routeur chi pour gerer les routes /books.
Codes HTTP respectes : 201 pour la creation, 200 pour la lecture, 204 pour la suppression.

ETAPE 6 QUALITE ET TESTS CURL

Format d erreur unique : { "error" "message" }.
Middleware de logs ajoute pour tracer la methode, le chemin et la duree des requetes.

EXEMPLES DE COMMANDES CURL POUR TESTER :

SANTE DE L API :
curl.exe http://localhost:8080/health 

CREER UN LIVRE :
curl.exe -X POST http://localhost:8080/books -H "Content-Type: application/json" -d "{"""title""":"""Dune""","""author""":"""Frank Herbert""","""year""":1965}" 

LISTER LES LIVRES :
curl.exe http://localhost:8080/books 

SUPPRIMER UN LIVRE :
curl.exe -X DELETE http://localhost:8080/books/ID_UUID 
