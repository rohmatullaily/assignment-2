package main

import (
	"assignment-2/models"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/antonlindstrom/pgstore"
	"golang.org/x/crypto/scrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var store *pgstore.PGStore
var db *gorm.DB

const (
	DBHost         = "localhost"
	DBUserName     = "postgres"
	DBUserPassword = "postgres"
	DBName         = "postgres"
	DBSchema       = "training"
	DBPort         = "5432"
)

func dbConnection() *gorm.DB {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai search_path=%s", DBHost, DBUserName, DBUserPassword, DBName, DBPort, DBSchema)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(0)
	}

	db.AutoMigrate(&models.User{})

	return db
}

func generateRandomKey(length int) ([]byte) {
	key := make([]byte, length)
	return key
}

func pgStore() *pgstore.PGStore {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DBUserName, DBUserPassword, DBHost, DBPort, DBName)
	authKey := generateRandomKey(32) 
	encryptionKey := generateRandomKey(32)

	var err error
	store, err = pgstore.NewPGStore(url, authKey, encryptionKey)
	if err != nil {
		log.Println("ERROR", err)
		os.Exit(1)
	}

	return store
}

func HashPassword(password string) (string, error) {
    salt := generateRandomKey(32)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    hashedPassword, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
    if err != nil {
        return "", err
    }

    encodedSalt := base64.StdEncoding.EncodeToString(salt)
    encodedPassword := base64.StdEncoding.EncodeToString(hashedPassword)

    return encodedSalt + "-" + encodedPassword, nil
}

func CheckPasswordHash(password, hash string) bool {
    parts := strings.Split(hash, "-")
    if len(parts) != 2 {
        return false
    }

    encodedSalt := parts[0]
    encodedHashedPassword := parts[1]

    salt, err := base64.StdEncoding.DecodeString(encodedSalt)
    if err != nil {
        return false
    }

    hashedPassword, err := base64.StdEncoding.DecodeString(encodedHashedPassword)
    if err != nil {
        return false
    }

    testHash, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
    if err != nil {
        return false
    }

    return bytes.Equal(testHash, hashedPassword)
}

func main() {
	gob.Register(models.User{})

	store = pgStore()
	db = dbConnection()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	http.HandleFunc("/", routeIndex)
	http.HandleFunc("/register", routeRegister)
	http.HandleFunc("/logout", routeLogout)

	fmt.Println("server started at localhost:5555")
	http.ListenAndServe(":5555", nil)
}

func routeIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		session, _ := store.Get(r, "session")
		val := session.Values["user"]
		user, ok := val.(models.User)

		if !ok {
			var tmpl = template.Must(template.New("login").ParseFiles("layout/login.html"))
			var err = tmpl.Execute(w, nil)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		var tmpl = template.Must(template.New("home").ParseFiles("layout/home.html"))
		var err = tmpl.Execute(w, user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var username = r.FormValue("username")
		var password = r.Form.Get("password")

		user := models.User{}

		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				var tmpl = template.Must(template.New("login").ParseFiles("layout/login.html"))
				var err = tmpl.Execute(w, models.Layout{
					Message: "Account not found, please check if the username and password are correct.",
				})

				if err != nil {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
				}

				return
			}
		}

		if ok := CheckPasswordHash(password, user.Password); !ok {
			var tmpl = template.Must(template.New("login").ParseFiles("layout/login.html"))
			var err = tmpl.Execute(w, models.Layout{
				Message: "wrong password",
			})

			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			return
		}

		session, err := store.Get(r, "session")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		session.Values["user"] = user
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func routeRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var tmpl = template.Must(template.New("register").ParseFiles("layout/register.html"))
		var err = tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var firstname = r.FormValue("firstname")
		var lastname = r.Form.Get("lastname")
		var username = r.Form.Get("username")
		var password = r.Form.Get("password")

		hashedPassword, _ := HashPassword(password)

		user := models.User{}

		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				http.Redirect(w, r, "/register", http.StatusSeeOther)
				return
			}
		}

		if user.ID > 0 {
			var tmpl = template.Must(template.New("register").ParseFiles("layout/register.html"))
			var err = tmpl.Execute(w, models.Layout{
				Message: "Account Already Exists",
			})

			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}

			return
		}

		user = models.User{
			FirstName: firstname,
			LastName:  lastname,
			Password:  hashedPassword,
			Username:  username,
		}

		if err := db.Create(&user).Error; err != nil {
			http.Redirect(w, r, "/register", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func routeLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		session, err := store.Get(r, "session")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Options.MaxAge = -1
		session.Save(r, w)

		var tmpl = template.Must(template.New("login").ParseFiles("layout/login.html"))
		err = tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}
