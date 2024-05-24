package helpers

import (
	"database/sql"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var DummyText = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris sit amet ipsum vitae risus ullamcorper facilisis in eu ex. Suspendisse porta rhoncus eros, sed feugiat risus rhoncus id. Cras dapibus urna eu pharetra semper. Morbi commodo est non nisi aliquam molestie. Integer pharetra nibh quis aliquet facilisis. Vivamus volutpat hendrerit velit, id pharetra magna blandit sit amet. Aenean diam sem, pellentesque nec ornare auctor, vestibulum eu magna. Cras placerat justo fermentum, iaculis nisl commodo, gravida arcu. Proin non purus eget velit maximus imperdiet sodales a turpis. Sed laoreet quam ut luctus scelerisque. Integer ut velit urna. Sed vel sapien fringilla, porta nunc in, bibendum dolor."
var Storage = "./views/assets/img/posts"

type FinalData struct {
	Posts      []ForumPost
	IsLogged   bool
	UserData   map[string]interface{}
	Categories []map[string]interface{}
}

type ForumPost struct {
	Id             string
	AuthorId       string
	AuthorName     string
	AuthorUsername string
	AuthorPicture  string
	Title          string
	Content        string
	Categories     string
	PostPicture    string
	CreationDate   string
	Comments       []PostComment
	Likes          int
	Dislikes       int
	NumComments    int
	Liked          bool
	Disliked       bool
}

type PostComment struct {
	Id             string
	AuthorName     string
	AuthorUsername string
	AuthorPicture  string
	Content        string
	Likes          int
	Dislikes       int
	Liked          bool
	Disliked       bool
}

type ValidationRules map[string]*regexp.Regexp

var rules = ValidationRules{
	"name":     regexp.MustCompile(`^.{2,20}$`),
	"email":    regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,25}$`),
	"password": regexp.MustCompile(`^.{8,30}$`),
	"username": regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`),
	"content":  regexp.MustCompile(`^.{0,500}$`),
	"comment":  regexp.MustCompile(`^.{1,500}$`),
	"category": regexp.MustCompile(`^[a-zA-Z0-9]{1,15}$`),
	"title":    regexp.MustCompile(`^.{2,100}$`),
}

func ValidateInput(champ, value string) bool {
	if strings.TrimSpace(value) != "" {
		if rule, ok := rules[champ]; ok {
			return rule.MatchString(value)
		}
	}

	return false
}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func UUID() string {
	// 1st fake category
	newID, err := uuid.NewV4()
	ErrorSimple(err)

	return newID.String()
}

func RandN(min, max int) int {
	return rand.Intn(max-min) + min
}

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	ErrorSimple(err)

	return string(hashedPassword)
}

func EmptyStringIfNil(w interface{}) string {
	if w == nil {
		return ""
	}
	return InterfaceToString(w)
}

func DefaultIfNil(w interface{}) string {
	if w == nil {
		return "default-user.png"
	}
	return InterfaceToString(w)
}

func InterfaceToString(w interface{}) string {
	return fmt.Sprintf("%s", w)
}

func InterfaceToSlice(data interface{}) []string {
	// Vérifier si data est une tranche vide
	if data == nil {
		return []string{}
	}

	// Tentative de conversion de data en tranche d'interface
	sliceData, ok := data.([]interface{})
	if !ok {
		// Si la conversion échoue, retourner une tranche vide
		return []string{}
	}

	// Initialiser une tranche de chaînes de caractères pour stocker les valeurs
	result := make([]string, len(sliceData))

	// Parcourir les éléments de la tranche d'interface et les convertir en chaînes de caractères
	for i, v := range sliceData {
		// Vérifier si l'élément est de type string
		if str, ok := v.(string); ok {
			result[i] = str
		} else {
			// Si l'élément n'est pas de type string, convertir en une représentation textuelle
			result[i] = fmt.Sprintf("%v", v)
		}
	}

	return result
}

func ExtractSQLRows(rows *sql.Rows) []map[string]interface{} {
	// Initialiser une liste vide pour stocker les résultats
	results := make([]map[string]interface{}, 0)

	// Récupérer les noms de colonnes
	cols, err := rows.Columns()
	ErrorSimple(err)

	// Parcourir chaque ligne
	for rows.Next() {
		// Initialiser une map pour stocker les valeurs de cette ligne
		rowValues := make(map[string]interface{})

		// Initialiser une tranche d'interfaces pour stocker les valeurs
		rowInterface := make([]interface{}, len(cols))
		for i := range rowInterface {
			rowInterface[i] = new(interface{})
		}

		// Scanner les valeurs de la ligne dans la tranche d'interfaces
		err := rows.Scan(rowInterface...)
		ErrorSimple(err)

		// Copier les valeurs de la tranche d'interfaces dans la map de valeurs
		for i, colName := range cols {
			rowValues[colName] = *(rowInterface[i].(*interface{}))
		}

		// Ajouter la map de valeurs à la liste des résultats
		results = append(results, rowValues)
	}

	return results
}
