package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

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
	Init()
	manager = NewManager()
	router := mux.NewRouter()
	router.HandleFunc("/api/users/{username}", Auth(RoleAssignment(handleUser))).Methods("GET", "DELETE", "PUT")
	router.HandleFunc("/api/users", Auth(RoleAssignment(handleUser))).Methods("POST")
	router.HandleFunc("/api/deposit", Auth(RoleAssignment(handleDeposit))).Methods("POST")
	router.HandleFunc("/quizzes/sell", Auth(RoleAssignment(handleSellQuiz))).Methods("POST")
	router.HandleFunc("/quizzes/buy", Auth(RoleAssignment(handleBuyQuiz))).Methods("POST")
	router.HandleFunc("/api/quizzes/{id}", Auth(RoleAssignment(handleQuizzes))).Methods("GET", "DELETE")
	router.HandleFunc("/api/quizzes", Auth(RoleAssignment(handleQuizzes))).Methods("POST", "PUT")
	router.HandleFunc("/api/register", RoleAssignment(handleRegister)).Methods("POST")
	router.HandleFunc("/api/login", handleLogin).Methods("POST")
	router.HandleFunc("/api/logout", Auth(RoleAssignment(handleLogout))).Methods("GET")
	router.HandleFunc("/quiz/comments", Auth(RoleAssignment(handleComments))).Methods("POST")
	router.HandleFunc("/quiz/comments/{id}", Auth(RoleAssignment(handleComments))).Methods("GET", "DELETE", "PATCH")
	router.HandleFunc("/quiz/ratings", Auth(RoleAssignment(handleRatings))).Methods("POST")
	router.HandleFunc("/quiz/ratings/{id}", Auth(RoleAssignment(handleRatings))).Methods("GET", "DELETE", "PATCH")
	router.HandleFunc("/game/create", Auth(RoleAssignment(handleCreateGame))).Methods("POST")
	router.HandleFunc("/game/{gameCode}/join", Auth(RoleAssignment(handleJoinGame))).Methods("POST")
	router.HandleFunc("/game/{gameCode}/start", Auth(RoleAssignment(handleStartGame))).Methods("POST")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	http.ListenAndServe(s.listenAddr, handler)
}

func handleStartGame(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "role not found in context", http.StatusInternalServerError)
		return
	}

	if role == Nuser {
		http.Error(w, "Missing permission", http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		userID, ok := r.Context().Value("userID").(uint)
		if !ok {
			http.Error(w, "userID not found in context", http.StatusInternalServerError)
			return
		}

		vars := mux.Vars(r)
		gameCode := vars["gameCode"]

		go StartGame(gameCode, userID)

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleCreateGame(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "role not found in context", http.StatusInternalServerError)
		return
	}

	if role == Nuser {
		http.Error(w, "Missing permission", http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		userID, ok := r.Context().Value("userID").(uint)
		if !ok {
			http.Error(w, "userID not found in context", http.StatusInternalServerError)
			return
		}

		var body CreateGameRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		code, err := CreateGame(userID, body.QuizId)
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
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "role not found in context", http.StatusInternalServerError)
		return
	}

	if role == Nuser {
		http.Error(w, "Missing permission", http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		userID, ok := r.Context().Value("userID").(uint)
		if !ok {
			http.Error(w, "userID not found in context", http.StatusInternalServerError)
			return
		}

		vars := mux.Vars(r)
		gameCode := vars["gameCode"]

		if err := JoinGame(gameCode, userID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		manager.ServeWS(w, r, gameCode)

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleSellQuiz(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "role not found in context", http.StatusInternalServerError)
		return
	}

	if role == Nuser {
		http.Error(w, "Missing permission", http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		var body SellQuizRequest
		userID, ok := r.Context().Value("userID").(uint)
		if !ok {
			http.Error(w, "userID not found in context", http.StatusInternalServerError)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := SellQuiz(&body, userID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleBuyQuiz(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "role not found in context", http.StatusInternalServerError)
		return
	}

	if role == Nuser {
		http.Error(w, "Missing permission", http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		var body BuyQuizRequest
		userID, ok := r.Context().Value("userID").(uint)
		if !ok {
			http.Error(w, "userID not found in context", http.StatusInternalServerError)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := BuyQuiz(&body, userID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleRatings(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "role not found in context", http.StatusInternalServerError)
		return
	}

	if role == Nuser {
		http.Error(w, "Missing permission", http.StatusForbidden)
		return
	}

	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
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

		err = DeleteRating(uint(id), userID, role)
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

		if err := ModifyRating(&body, userID, role); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleComments(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "role not found in context", http.StatusInternalServerError)
		return
	}

	if role == Nuser {
		http.Error(w, "Missing permission", http.StatusForbidden)
		return
	}

	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		http.Error(w, "userID not found in context", http.StatusInternalServerError)
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

		err = DeleteComment(uint(id), userID, role)
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

		if err := CreateComment(&body, userID); err != nil {
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

		if err := ModifyComment(&body, userID, role); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	delete(session.Values, "userID")
	delete(session.Values, "role")
	session.Save(r, w)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "role not found in context", http.StatusInternalServerError)
		return
	}

	if role != Nuser {
		http.Error(w, "Already logged in", http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		var body CreateAccountRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := CreateAccount(&body, role); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode("Account created")
		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	/*
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "role not found in context", http.StatusInternalServerError)
		return
	}
	*/

	session, _ := store.Get(r, "session")
	role := session.Values["role"]

	if role == nil {
		role = Nuser
	}

	if role != Nuser {
		http.Error(w, "Already logged in", http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		var body LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id, role, err := Login(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		session, _ := store.Get(r, "session")
		session.Values["userID"] = id
		session.Values["role"] = role
		store.Save(r, w, session)

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleQuizzes(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "role not found in context", http.StatusInternalServerError)
		return
	}

	if role == Nuser {
		http.Error(w, "Missing permission", http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		var body CreateQuizRequest
		session, _ := store.Get(r, "session")

		userId, ok := session.Values["userID"].(uint)
		if !ok {
			http.Error(w, "userID not found in context", http.StatusInternalServerError)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := CreateQuiz(&body, userId); err != nil {
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
		if role != Admin {
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
		if err := ModifyQuiz(&body, userId, role); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleDeposit(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "role not found in context", http.StatusInternalServerError)
		return
	}

	if role != Ruser {
		http.Error(w, "Missing permission", http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		var body DepositRequest
		userId, ok := r.Context().Value("userID").(uint)
		if !ok {
			http.Error(w, "userID not found in context", http.StatusInternalServerError)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := Deposit(&body, userId); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		return
	}

	http.Error(w, "The payload is in an unsupported format", http.StatusUnsupportedMediaType)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	role, ok := r.Context().Value("role").(string)
	if !ok {
		http.Error(w, "role not found in context", http.StatusInternalServerError)
		return
	}

	if role == Nuser {
		http.Error(w, "Missing permission", http.StatusForbidden)
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
		if role != Admin {
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

		userId, ok := r.Context().Value("userID").(uint)
		if !ok {
			http.Error(w, "userID not found in context", http.StatusInternalServerError)
			return
		}

		var body ModifyAccountRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := EditAccount(username, userId, role, body)
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
		session, _ := store.Get(r, "session")
		userId, ok := session.Values["userID"]
		if !ok {
			http.Error(w, "invalid access token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userId)

		HandlerFunc.ServeHTTP(w, r.WithContext(ctx))
	}
}

func RoleAssignment(HandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		role, ok := session.Values["role"]
		if !ok {
			role = Nuser
			session.Values["role"] = Nuser
			session.Save(r, w)
		}

		ctx := context.WithValue(r.Context(), "role", role)

		HandlerFunc.ServeHTTP(w, r.WithContext(ctx))
	}
}
