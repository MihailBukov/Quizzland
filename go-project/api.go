package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var manager *Manager

type APIServer struct {
	listenAddr string
}

func NewApiServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	manager = NewManager()
	router := mux.NewRouter()
	router.HandleFunc("/api/users/{username}", Auth(handleUser)).Methods("GET", "DELETE", "PUT")
	router.HandleFunc("/api/users", Auth(handleUser)).Methods("POST")
	router.HandleFunc("/api/deposit", Auth(handleDeposit)).Methods("POST")
	router.HandleFunc("/quizzes/sell", Auth(handleSellQuiz)).Methods("POST")
	router.HandleFunc("/quizzes/buy", Auth(handleBuyQuiz)).Methods("POST")
	router.HandleFunc("/api/quizzes/{id}", Auth(handleQuizzes)).Methods("GET", "DELETE")
	router.HandleFunc("/api/quizzes", Auth(handleQuizzes)).Methods("POST", "PUT")
	router.HandleFunc("/api/register", handleRegister).Methods("POST")
	router.HandleFunc("/api/login", handleLogin).Methods("POST")
	router.HandleFunc("/api/logout", Auth(handleLogout)).Methods("GET")
	router.HandleFunc("/quiz/comments", Auth(handleComments)).Methods("POST")
	router.HandleFunc("/quiz/comments/{id}", Auth(handleComments)).Methods("GET", "DELETE", "PATCH")
	router.HandleFunc("/quiz/ratings", Auth(handleRatings)).Methods("POST")
	router.HandleFunc("/quiz/ratings/{id}", Auth(handleRatings)).Methods("GET", "DELETE", "PATCH")
	router.HandleFunc("/game/create", Auth(handleCreateGame)).Methods("POST")
	router.HandleFunc("/game/{gameCode}/join", Auth(handleJoinGame)).Methods("POST")
	router.HandleFunc("/game/{gameCode}/start", Auth(handleStartGame)).Methods("POST")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	http.ListenAndServe(s.listenAddr, handler)
}

func handleStartGame(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(UserContext)
	if !ok {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		vars := mux.Vars(r)
		gameCode := vars["gameCode"]

		go StartGame(gameCode, user.UserID)

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleCreateGame(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(UserContext)
	if !ok {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		var body CreateGameRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		code, err := CreateGame(user.UserID, body.QuizId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		manager.ServeWS(w, r, code)

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleJoinGame(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(UserContext)
	if !ok {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		vars := mux.Vars(r)
		gameCode := vars["gameCode"]

		if err := JoinGame(gameCode, user.UserID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		manager.ServeWS(w, r, gameCode)

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleSellQuiz(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(UserContext)
	if !ok {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		var body SellQuizRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := SellQuiz(&body, user.UserID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleBuyQuiz(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(UserContext)
	if !ok {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		var body BuyQuizRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := BuyQuiz(&body, user.UserID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleRatings(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(UserContext)
	if !ok {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}

	if r.Method == "GET" {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		quizzes, err := GetRatingsForProduct(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(quizzes)
		return
	} else if r.Method == "DELETE" {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = DeleteRating(uint(id), user.UserID, user.Role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	} else if r.Method == "POST" {
		var body CreateRatingRequest
		userID, ok := r.Context().Value("userID").(uint)
		if !ok {
			http.Error(w, "userID not found in context", http.StatusInternalServerError)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := CreateRating(&body, userID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	} else if r.Method == "PATCH" {
		var body ModifyRatingRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := ModifyRating(&body, user.UserID, user.Role); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleComments(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(UserContext)
	if !ok {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}

	if r.Method == "GET" {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		quizzes, err := GetCommentsForProduct(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(quizzes)
		return
	} else if r.Method == "DELETE" {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = DeleteComment(uint(id), user.UserID, user.Role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	} else if r.Method == "POST" {
		var body CreateCommentRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := CreateComment(&body, user.UserID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	} else if r.Method == "PATCH" {
		var body ModifyCommentRequest
		userID, err := r.Context().Value("userID").(uint)
		if !err {
			http.Error(w, "userID not found in context", http.StatusInternalServerError)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := ModifyComment(&body, userID, user.Role); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	if !DoesUserHaveValidCookie(r) {
		http.Error(w, "not logged in", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	})
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(UserContext)
	if !ok {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}

	if user.Role == Ruser {
		http.Error(w, "already logged in", http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		var body CreateAccountRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := CreateAccount(&body, user.Role); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode("Account created")
		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if DoesUserHaveValidCookie(r) {
		http.Error(w, "Already logged in", http.StatusBadRequest)
            return
	}

	if r.Method == "POST" {
		var body LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tokenString, expTime, err := Login(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expTime,
			Path:    "/",
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			Secure: true,
		})

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleQuizzes(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(UserContext)
	if !ok {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		var body CreateQuizRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := CreateQuiz(&body, user.UserID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	} else if r.Method == "GET" {
		quizzes, err := GetQuizzesForSale()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(quizzes)
		return
	} else if r.Method == "DELETE" {
		if user.Role != Admin {
			http.Error(w, "Missing permission", http.StatusForbidden)
			return
		}

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := DeleteQuiz(id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	} else if r.Method == "PUT" {
		var body ModifyQuizRequest
		userId, ok := r.Context().Value("userID").(uint)
		if !ok {
			http.Error(w, "userID not found in context", http.StatusInternalServerError)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := ModifyQuiz(&body, userId, user.Role); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleDeposit(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(UserContext)
	if !ok {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		var body DepositRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := Deposit(&body, user.UserID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(UserContext)
	if !ok {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}

	if r.Method == "GET" {
		vars := mux.Vars(r)
		username := vars["username"]

		acc, err := GetAccountDto(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(acc)
		return
	} else if r.Method == "DELETE" {
		if user.Role != Admin {
			http.Error(w, "Missing permission", http.StatusForbidden)
			return
		}

		vars := mux.Vars(r)
		username := vars["username"]

		err := DeleteAccount(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	} else if r.Method == "PUT" {
		vars := mux.Vars(r)
		username := vars["username"]

		var body ModifyAccountRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := EditAccount(username, user.UserID, user.Role, body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func Auth(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("token")
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        tokenString := cookie.Value

        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return []byte("SECRET_KEY"), nil
        })
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        if !token.Valid {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
		}

		userCtx := UserContext{
            Role:   claims.Role,
            UserID: claims.Id,
        }

		ctx := context.WithValue(r.Context(), "user", userCtx)

		HandlerFunc.ServeHTTP(w, r.WithContext(ctx))
	}
}
