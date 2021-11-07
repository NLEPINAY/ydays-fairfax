package request

import (
	"errors"
	"forum/database"
	"log"
	"net/http"
	"strings"
)

func Register(w http.ResponseWriter, r *http.Request, user database.User) {
	switch r.Method {

	// 🍔 Méthode 'GET' — Lorsqu'on arrive sur la page register.html pour la 1ère fois :
	case "GET":
		err := MyTemplates.ExecuteTemplate(w, "signin", nil)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}

	// 🍔 Méthode 'POST' — Lorsqu'on sur le bouton 'Create your account' pour s'enregistrer :
	case "POST":
		// Je récupère l'email, le nom d'utilisateur, le mot de passe et la date actuelle :
		email := r.FormValue("signup-email")
		username := r.FormValue("signup-username")
		password := r.FormValue("signup-password")

		// (1) Vérifier si les valeurs envoyées via le formulaire sont valides :
		/**
		*todo: Code to be written here...
		 */

		user.Username = strings.ToLower(username)
		user.Password = password
		user.Email = strings.ToLower(email)

		// (2) Vérifier si l'email ou username existe déjà dans la base de données :
		errorMessage := CheckIfExists(user)
		if errorMessage != "" {
			log.Println("❌ REGISTER | Access denied : ", errorMessage)
			w.Write([]byte(errorMessage))
			return
		}

		// (3) Ajouter l'utilisateur dans la base de données :
		err := user.InsertIntoDatabase()
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}

		// (4) Ajouter la session de l'utilisateur à la base de données :
		err = database.AddSessionToDatabase(w, r, user)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}

		// (5) Après s'être inscrit, on renvoie une réponse Success au JS :
		var message = `{
			"success": true,
			"message": "Account successfully created."}`

		w.Write([]byte(message))
		log.Println("✔️ REGISTER | Account created successfully.")
		log.Println("Successfully registered: ", user)
	}
}

func CheckIfExists(user database.User) string {
	var errorMessage string
	var userToCheck database.User

	// Je vérifie si le username voulu existe déjà dans la DB, et si oui, je l'ajoute dans userToCheck :
	nameInDatabase := database.Db.QueryRow("SELECT username FROM user WHERE username = $1 OR username = $2 OR username = $3", user.Username, strings.ToLower(user.Username), strings.ToUpper(user.Username))
	nameInDatabase.Scan(&userToCheck.Username)

	// Je vérifie si l'email voulu existe déjà dans la DB, et si oui, je l'ajoute dans userToCheck :
	emailInDatabase := database.Db.QueryRow("SELECT email FROM user WHERE email = ?", user.Email)
	emailInDatabase.Scan(&userToCheck.Email)

	// Si un nom ou email a été ajouté à userToCheck, cela veut dire que le nom ou email est déjà pris :
	if userToCheck.Username != "" {
		errorMessage = `{
			"success": false,
			"error": "username",
			"message": "This username is unavailable."}`
		return errorMessage
	}

	if userToCheck.Email != "" {
		errorMessage = `{
			"success": false,
			"error": "email",
			"message": "This email is already registered."}`
		return errorMessage
	}
	return errorMessage
}

func CheckNewAccount(user database.User) database.ErrorData {
	var errorData database.ErrorData
	var userToCheck database.User

	// Je vérifie si le username voulu existe déjà dans la DB, et si oui, je l'ajoute dans userToCheck :
	nameInDatabase := database.Db.QueryRow("SELECT username FROM user WHERE username = $1 OR username = $2 OR username = $3", user.Username, strings.ToLower(user.Username), strings.ToUpper(user.Username))
	nameInDatabase.Scan(&userToCheck.Username)

	// Je vérifie si l'email voulu existe déjà dans la DB, et si oui, je l'ajoute dans userToCheck :
	emailInDatabase := database.Db.QueryRow("SELECT email FROM user WHERE email = ?", user.Email)
	emailInDatabase.Scan(&userToCheck.Email)

	// Si un nom ou email a été ajouté à userToCheck, cela veut dire que le nom ou email est déjà pris :
	if userToCheck.Username != "" && userToCheck.Email != "" {
		errorData.Account = errors.New("account already existing")
		return errorData
	}

	if userToCheck.Username != "" {
		errorData.Username = errors.New("username unavailable")
		return errorData
	}

	if userToCheck.Email != "" {
		errorData.Email = errors.New("email already registered")
		return errorData
	}

	return errorData
}

func RegisterSecret(w http.ResponseWriter, r *http.Request, user database.User) {
	switch r.Method {

	// 🍔 Méthode 'GET' — Lorsqu'on arrive sur la page pour la 1ère fois :
	case "GET":
		err := MyTemplates.ExecuteTemplate(w, "register-secret", nil)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}

	// 🍔 Méthode 'POST' — Lorsqu'on sur le bouton pour s'enregistrer :
	case "POST":
		// Je récupère la question et la réponse secrètes :
		question := r.FormValue("secret-question")
		answer := r.FormValue("secret-answer")

		user.SecretQuestion = question
		user.SecretAnswer = answer

		// (3) Modification de l'utilisateur dans la base de données :
		err := user.UpdateInDatabase("secretQuestion")
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}

		err = user.UpdateInDatabase("secretAnswer")
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}

		// Modification du statut de l'utilisateur (de INCOMPLETE à NORMAL) :
		user.State = database.NORMAL
		err = user.UpdateInDatabase("state")
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}

		// (4) Ajouter la session de l'utilisateur à la base de données :
		err = database.AddSessionToDatabase(w, r, user)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		}

		// Redirection vers la page de login :
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		log.Println(user.SecretQuestion)
		log.Println(user.SecretAnswer)
	}
}
