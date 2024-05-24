package controllers

import (
	"errors"
	"fmt"
	"forum/global"
	"forum/helpers"
	"forum/models"
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

// Vérification de l'existence de la session
// Si l'utilisateur est déjà connecté, on le redirige vers la page principale
// et s'il n'est pas connecté, on lui affiche la page de connexion
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	IsLogged, _ := Auth(r, w)
	if IsLogged {
		// Contenu HTML stocké dans une variable
        htmlContent := `<html><script>window.location.href = "/";</script></html>`

        // Création d'une nouvelle instance de template
        tmpl, err := template.New("monTemplate").Parse(htmlContent)
        helpers.ErrorSimple(err)

        // Exécution du template et écriture dans la réponse HTTP
        err = tmpl.Execute(w, nil)
        helpers.ErrorSimple(err)

		return
	}
	global.W = w

	if !IsLogged {
		if r.Method == http.MethodPost {
			Login(w, r)
			return
		}

		tmpl, err := template.ParseFiles("./views/auth.html")
		helpers.ErrorSimple(err)
		tmpl.Execute(w, nil)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	global.W = w
	helpers.AllowedMethod(w, "POST", r)
	if r.URL.Path != "/login" {
		return
	}

	UsernameOrEmail := html.EscapeString(r.FormValue("username/email"))

	if strings.Contains(UsernameOrEmail, "@") {
		if !helpers.ValidateInput("email", UsernameOrEmail) {
			helpers.ErrorHandler(w, errors.New("INVALID EMAIL"), "INVALID EMAIL", http.StatusBadRequest)
			return
		}
	} else {
		if !helpers.ValidateInput("username", UsernameOrEmail) {
			helpers.ErrorHandler(w, errors.New("INVALID USERNAME"), "INVALID USERNAME", http.StatusBadRequest)
			return
		}
	}

	password := html.EscapeString(r.FormValue("password"))
	if !helpers.ValidateInput("password", password) {
		helpers.ErrorHandler(w, errors.New("INVALID PASSWORD"), "INVALID PASSWORD", http.StatusBadRequest)
		return
	}

	UserEmail := models.User.Fetch([]interface{}{"email"}, []interface{}{UsernameOrEmail})
	UserName := models.User.Fetch([]interface{}{"username"}, []interface{}{UsernameOrEmail})
	var user []models.Model
	if len(UserEmail) == 1 {
		user = UserEmail
	} else {
		user = UserName
	}

	if len(user) == 1 {
		hashedPassword := fmt.Sprintf("%s", user[0].Values["password"])

		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err == nil {
			CreateSession(w, r, fmt.Sprintf("%s", user[0].Values["id"]))

			helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
			return

		} else {
			helpers.ErrorHandler(w, errors.New("INVALID PASSWORD"), "INVALID PASSWORD", http.StatusNotFound)
			return
		}

	} else {
		helpers.ErrorHandler(w, errors.New("USER NOT FOUND"), "USER NOT FOUND", http.StatusNotFound)
		return
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	global.W = w
	helpers.AllowedMethod(w, "GET", r)
	Logout(w, r)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	global.W = w
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	sessionID := cookie.Value

	DeleteSession(w, r, sessionID)
	http.Redirect(w, r, "/login", http.StatusFound)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	global.W = w
	IsLogged, _ := Auth(r, w)
	if IsLogged {
		// Contenu HTML stocké dans une variable
        htmlContent := `<html><script>window.location.href = "/";</script></html>`

        // Création d'une nouvelle instance de template
        tmpl, err := template.New("monTemplate").Parse(htmlContent)
        helpers.ErrorSimple(err)

        // Exécution du template et écriture dans la réponse HTTP
        err = tmpl.Execute(w, nil)
        helpers.ErrorSimple(err)

		return
	}

	if !IsLogged {
		if r.Method == http.MethodPost {
			Register(w, r)
			return
		}

		tmpl, err := template.ParseFiles("./views/auth.html")
		helpers.ErrorSimple(err)
		tmpl.Execute(w, nil)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	global.W = w
	helpers.AllowedMethod(w, "POST", r)

	name := html.EscapeString(r.FormValue("name"))
	if !helpers.ValidateInput("name", name) {
		helpers.ErrorHandler(w, errors.New("INVALID NAME"), "INVALID NAME", http.StatusBadRequest)
		return
	}

	email := html.EscapeString(r.FormValue("email"))
	if !helpers.ValidateInput("email", email) {
		helpers.ErrorHandler(w, errors.New("INVALID EMAIL"), "INVALID EMAIL", http.StatusBadRequest)
		return
	}

	username := html.EscapeString(r.FormValue("username"))
	if !helpers.ValidateInput("username", username) {
		helpers.ErrorHandler(w, errors.New("INVALID USERNAME"), "INVALID USERNAME", http.StatusBadRequest)
		return
	}

	password := html.EscapeString(r.FormValue("password"))
	if !helpers.ValidateInput("password", password) || len(password) < 8 || len(password) > 20 {
		helpers.ErrorHandler(w, errors.New("INVALID PASSWORD"), "INVALID PASSWORD", http.StatusBadRequest)
		return
	}

	// Vérifiez si l'email ou le nom d'utilisateur est déjà utilisé
	resEmail := models.User.Fetch([]interface{}{"email"}, []interface{}{email})
	resUsername := models.User.Fetch([]interface{}{"username"}, []interface{}{username})
	if len(resEmail) > 0 || len(resUsername) > 0 {
		// Renvoyez une erreur si l'email ou le nom d'utilisateur est déjà utilisé
		helpers.ErrorHandler(w, errors.New("USERNAME OR EMAIL ALREADY USED"), "USERNAME OR EMAIL ALREADY USED", http.StatusConflict)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	helpers.ErrorSimple(err)

	picture, handler, _ := r.FormFile("picture")
	// fmt.Println(picture)

	// Si non, créez un nouvel utilisateur
	newUser := models.User
	newUser.Values = map[string]interface{}{
		"name":     name,
		"username": username,
		"email":    email,
		"password": helpers.HashPassword(password),
	}

	if picture == nil {
		newUser.Values["picture"] = "default-user.png"
	} else {
		fileBytes, err := io.ReadAll(picture)
		if err != nil {
			helpers.ErrorSimple(err)
			helpers.JSONResponse(w, errors.New("FAILED TO READ IMAGE CONTENT"), http.StatusInternalServerError)
			return
		}

		contentType := http.DetectContentType(fileBytes)

		if contentType != "image/png" && contentType != "image/jpeg" && contentType != "image/gif" {
			helpers.ErrorHandler(w, errors.New("UNSUPPORTED FILE FORMAT"), contentType, http.StatusBadRequest)
			return
		}

		os.MkdirAll("./views/assets/img/profile", os.ModePerm)

		fileName := fmt.Sprintf("%s-%s", helpers.UUID(), handler.Filename)
		picturePath := filepath.Join("./views/assets/img/profile", fileName)

		out, err := os.Create(picturePath)
		if err != nil {
			helpers.ErrorSimple(err)
			helpers.JSONResponse(w, errors.New("FAILED TO CREATE IMAGE FILE"), http.StatusInternalServerError)
			return
		}
		defer out.Close()

		_, err = out.Write(fileBytes)
		if err != nil {
			helpers.ErrorSimple(err)
			helpers.JSONResponse(w, errors.New("FAILED TO WRITE IMAGE CONTENT"), http.StatusInternalServerError)
			return
		}

		newUser.Values["picture"] = fileName

	}

	res := newUser.Create()

	// Connectez l'utilisateur après l'inscription
	CreateSession(w, r, fmt.Sprintf("%s", res.Values["id"]))

	helpers.JSONResponse(w, errors.New("SUCCESS"), http.StatusAccepted)
}

func Auth(r *http.Request, w http.ResponseWriter) (bool, []models.Model) {
	global.W = w
	var IsLogged bool = true
	Cookie, err := r.Cookie("sessionID")
	if err != nil {
		IsLogged = false
	} else {
		sessions := models.Session.Fetch([]interface{}{"Id"}, []interface{}{Cookie.Value})
		if len(sessions) == 0 {
			return false, nil
		}

		userID, OK := sessions[0].Values["idUser"]
		if !OK {
			return false, nil
		}

		users := models.User.Fetch([]interface{}{"Id"}, []interface{}{userID})
		if len(users) > 0 {
			IdUser = users[0].Values["id"].(string)
			return true, users
		} else {
			return false, nil
		}
	}
	return IsLogged, nil
}
