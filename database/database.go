/*Diagramme Merise :
Boite carré : table
Boite arrondie : relation (avec des verbes à l'infinitif)
Attribut souligné : clé primaire
Voir : looping-mcd
*/

package database

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func Initialize() {
	// Déclaration de toutes les tables de la base de données :
	dbTables := []string{
		`CREATE TABLE IF NOT EXISTS "user" (
			"id_user"			INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			"username"			TEXT UNIQUE NOT NULL,
			"password"			TEXT NOT NULL,
			"email"				TEXT UNIQUE NOT NULL,
			"role"				INTEGER DEFAULT 0,
			"avatar"			TEXT DEFAULT "/images/avatars/defaultAvatar.jpg",
			"date_user"			DATETIME DEFAULT CURRENT_TIMESTAMP,
			"state_user"		INTEGER DEFAULT 0,		
			"secretQuestion"	TEXT DEFAULT '',
			"secretAnswer"		TEXT DEFAULT '',
			"house_id"			INTEGER DEFAULT 0
		)`,

		`CREATE TABLE IF NOT EXISTS "post" (
			"id_post"				INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			"title_post"			TEXT NOT NULL,
			"author_id"				INTEGER NOT NULL,
			"content_post"			TEXT NOT NULL,
			"category_id"			INTEGER NOT NULL,
			"date_post"				DATETIME DEFAULT CURRENT_TIMESTAMP,
			"state_post"			INTEGER DEFAULT 0,
			"promoted"				INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY(author_id) REFERENCES "user"(id_user), 
			FOREIGN KEY(category_id) REFERENCES "category"(id_category) 
		)`,

		`CREATE TABLE IF NOT EXISTS "comment" (
			"id_comment"			INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			"author_id"				INTEGER NOT NULL,
			"post_id"				INTEGER NOT NULL,
			"content_comment"		TEXT NOT NULL,
			"gif"					TEXT NOT NULL DEFAULT '',
			"date_comment"			DATETIME DEFAULT CURRENT_TIMESTAMP,
			"state_comment"			INTEGER DEFAULT 0,
			"reason_comment"		TEXT DEFAULT "Supprimer par l'utilisateur lui même",
			FOREIGN KEY(author_id) REFERENCES "user"(id_user),
			FOREIGN KEY(post_id) REFERENCES "post"(id_post)
		)`,

		`CREATE TABLE IF NOT EXISTS "session" (
			"id_session"		INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			"user_id"			INTEGER NOT NULL,
			"uuid"				TEXT NOT NULL,
			"date_session"		DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES "user"(id_user)
		)`,

		`CREATE TABLE IF NOT EXISTS "category" (
			"id_category"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			"name_category"	TEXT NOT NULL UNIQUE,
			"theme" TEXT NOT NULL DEFAULT ' ',
			"description"	TEXT NOT NULL DEFAULT ' '
		)`,

		`CREATE TABLE IF NOT EXISTS "post_like" (
			"post_id"			INTEGER NOT NULL,
			"user_id"			INTEGER NOT NULL,
			"type"				TEXT NOT NULL,	
			"date_post_like"	DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(post_id) REFERENCES "post"(id_post),
			FOREIGN KEY(user_id) REFERENCES "user"(id_user),
			PRIMARY KEY(post_id, user_id)
		)`,

		`CREATE TABLE IF NOT EXISTS "comment_like" (
			"comment_id"	INTEGER NOT NULL,
			"user_id"		INTEGER NOT NULL,
			"type"			TEXT NOT NULL,	
			"date"			DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(comment_id) REFERENCES "comment"(id_user),
			FOREIGN KEY(user_id) REFERENCES "user"(id_user),
			PRIMARY KEY(comment_id, user_id)
		)`,

		`CREATE TABLE IF NOT EXISTS "badge" (
			"id_badge"			INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			"type"				TEXT NOT NULL,
			"image_badge"				TEXT NOT NULL
		)`,

		`CREATE TABLE IF NOT EXISTS "user_badge" (
			"user_id"		INTEGER NOT NULL,
			"badge_id"		INTEGER NOT NULL,
			FOREIGN KEY(user_id) REFERENCES "user"(id_user),
			FOREIGN KEY(badge_id) REFERENCES "badge"(id_badge),
			PRIMARY KEY(user_id, badge_id)
		)`,

		`CREATE TABLE IF NOT EXISTS "ticket" (
			"id_ticket"			INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			"author_id"		INTEGER NOT NULL,
			"actual_admin" 	INTERGER NOT NULL,
			"title_ticket"			TEXT NOT NULL,
			"content_ticket"		TEXT NOT NULL,
			"date_ticket"			DATETIME DEFAULT CURRENT_TIMESTAMP,
			"state_ticket"			INTEGER DEFAULT 0,
			FOREIGN KEY(author_id) REFERENCES "user"(id_user),
			FOREIGN KEY(actual_admin) REFERENCES "user"(id_user)
		)`,

		`CREATE TABLE IF NOT EXISTS "house" (
			"id_house"			INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			"name"			TEXT NOT NULL,
			"image_house"			TEXT NOT NULL
		)`,

		`CREATE TABLE IF NOT EXISTS "chat" (
			"id_chat"				INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			"date_creation"			TEXT NOT NULL,
			"state"					TEXT NOT NULL
		)`,

		`CREATE TABLE IF NOT EXISTS "message" (
			"id_message"			INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			"id_chat"				INTEGER NOT NULL,
			"date_creation"			TEXT NOT NULL,
			"author"				INTEGER NOT NULL,
			"content_message"				TEXT NOT NULL,
			"state_message"					TEXT NOT NULL,
			FOREIGN KEY(id_chat) REFERENCES "chat"(id_chat),
			FOREIGN KEY(author) REFERENCES "user"(user_id)  
		)`,

		`CREATE TABLE IF NOT EXISTS "assoc_chat" (
			"id_assoc"		INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
			"id_chat"		INTEGER NOT NULL,
			"user_id"		INTEGER NOT NULL,
			"author_ID"		INTEGER NOT NULL,
			"content_assoc"		TEXT NOT NULL,
			"date_assoc"			DATETIME DEFAULT CURRENT_TIMESTAMP,
			"image_assoc"			TEXT DEFAULT '',
			"state_assoc"			INTEGER DEFAULT 0,	
			"reason_assoc"		TEXT DEFAULT "Supprimer par l'utilisateur lui même",
			FOREIGN KEY(author_ID) REFERENCES "user"(id_user), 
			FOREIGN KEY(id_chat) REFERENCES "chat"(id_chat) 
		)`,
	}
	var err error
	Db, err = sql.Open("sqlite3", "./database/database.db")
	if err != nil {
		log.Println("❌ ERREUR | Impossible de créer le fichier database.db")
		panic(err)
	}

	// Création de chaque table de la base de données :
	for _, table := range dbTables {
		err := createDatabase(table)
		if err != nil {
			panic(err)
		}
	}
	log.Println("✔️ DATABASE | Database created and initialized successfully.")
}

func createDatabase(table string) error {
	statement, err := Db.Prepare(table)
	defer statement.Close()
	if err != nil {
		log.Println("❌ ERREUR | Impossible de créer les tables de la base de données.")
		return err
	}
	statement.Exec()
	return nil
}

func AddSessionToDatabase(w http.ResponseWriter, r *http.Request, user User) error {
	// Je supprime la session précédente de l'utilisateur et en créé une nouvelle :
	Db.Exec("DELETE FROM session WHERE user_id = $1", user.ID)
	log.Println("Adding session to database with user's ID : ", user.ID)

	sessionID := uuid.New()
	cookie := &http.Cookie{
		Name:   "session",
		Value:  sessionID.String(),
		Secure: true,
	}
	cookie.MaxAge = 60 * 60 * 24 // 24 heures
	http.SetCookie(w, cookie)

	// Insertion des valeurs de la session dans la table 'session' :
	statement, err := Db.Prepare("INSERT INTO session (user_id, uuid, date_session) VALUES (?, ?, ?)")
	defer statement.Close()
	if err != nil {
		log.Println("❌ ERREUR | Impossible d'insérer la session dans la base de données.")
		log.Println("Hypothèse : Mauvaise syntaxe du statement SQLite “INSERT INTO session (user_id, uuid, date) VALUES (", user.ID, sessionID, time.Now().Add(24*time.Hour), ")”")
		return err
	}

	statement.Exec(user.ID, sessionID, time.Now().Add(24*time.Hour))
	// statement.Exec(user.ID, sessionID, time.Now().Add(60*time.Minute)) // Heure actuelle + 60 minutes
	return nil
}

func CleanExpiredsession() {
	for {
		Db.Exec("DELETE FROM session WHERE date < $1", time.Now())
		time.Sleep(10 * time.Minute)
	}
}
