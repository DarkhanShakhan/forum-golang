package authentication

import (
	"database/sql"
	"encoding/json"
	"forum/pkg/sqlite3"
	"log"
	"net/http"
)

func Run() {
	db, err := sqlite3.New()
	if err != nil {
		log.Fatal(err)
	}
	service := NewService(db)
	http.HandleFunc("/session", service.GetUserId)
	http.ListenAndServe("localhost:7777", nil)
}

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db}
}

type Message struct {
	Session    string `json:"session"`
	UserId     int    `json:"user_id,omitempty"`
	ExpiryDate string `json:"exp_date,omitempty"`
}

func (s *Service) GetUserId(w http.ResponseWriter, r *http.Request) {
	msg := Message{}
	json.NewDecoder(r.Body).Decode(&msg)

	tx, _ := s.db.Begin()
	stmt, _ := tx.Prepare("SELECT * FROM sessions WHERE session = ?;")
	rows, _ := stmt.Query(msg.Session)
	if rows.Next() {
		rows.Scan(&msg.Session, &msg.UserId, &msg.ExpiryDate)
	}
	out, _ := json.Marshal(msg)
	w.Header().Set("cookies", "heres")
	w.Write(out)

}
