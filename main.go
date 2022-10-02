package main

import (
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Comment struct {
	Message string `validate:"required,min=1,max=140"`
	UserName string `validate:"required,min=1,max=15"`
}

type Book struct {
	Title string `validate:"required"`
	Price *int `validate:"required"`
}

// type ResponseWriter interface{
// 	Header() Header
// 	Write([]byte) (int, error)
// 	WtiteHeader(statusCode int)
// }

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}


func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	http.Handle("/healthz", MiddlewareLogging(http.HandlerFunc(Healthz)))
	http.ListenAndServe(":8888", nil)
}

func MiddlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("start %s\n", r.URL)
		next.ServeHTTP(w, r)
		log.Printf("finish %s\n", r.URL)
	})
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter{
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func wrapHandlerWithLogging(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, req)

		statusCode := lrw.statusCode

		log.Printf("%d %s", statusCode, http.StatusText(statusCode))
	})
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode >= 400 {
		log.Printf("Response Body: %s", b)
	}
	return lrw.ResponseWriter.Write(b)
}


// func main() {
// 	err := r.ParseForm()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	//指定されたパラーメーターを返す(GETクエリーとPOST/PUT/PATCHのx-www-form-urlencoded)
// 	word := r.FormValue("searchword")
// 	log.Printf("searchword = %s\n", word)

// 	//mapとしてアクセス
// 	words, ok := r.Form["searchword"]
// 	log.Printf("search words = %v has values %v\n", words, ok)

// 	//全部のクエリーをループでアクセス
// 	log.Printf("search words = %v has values %v\n", words, ok)

// 	//全部のクエリーをループでアクセス
// 	log.Print("all queries")
// 	for key, values := range r.Form {
// 		log.Printf("  %s: %v\n", key, values)
// 	}
// }
// 	s := `{"Title":"Real World HTTP ミニ版", "Price":0}`
// 	var b Book
// 	if err := json.Unmarshal([]byte(s), &b); err != nil {
// 		log.Fatal(err)
// 	}

// 	if err := validator.New().Struct(b); err != nil {
// 		var ve validator.ValidationErrors
// 			if errors.As(err, &ve) {
// 				for _, fe := range ve {
// 					fmt.Printf("フィールド%s が %s違反です(値:%v)\n", fe.Field(), fe.Tag(),fe.Value())
// 				}
			
// 		}
// 	}
// 	var mutex = &sync.RWMutex{}
// 	comments := make([]Comment, 0, 100)

// 	http.HandleFunc("/comments", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")

// 		switch r.Method {
// 		case http.MethodGet:
// 			mutex.RLock()

// 			if err := json.NewEncoder(w).Encode(comments); err != nil {
// 				http.Error(w, fmt.Sprintf(`{"status":"%s"}`, err), http.StatusInternalServerError)
// 				return 
// 			}
// 			mutex.RUnlock()
		
// 		case http.MethodPost:
// 			var c Comment
// 			if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
// 				http.Error(w, fmt.Sprintf(`{"status":"%s"}`, err), http.StatusInternalServerError)
// 				return
// 			}
			
// 			validate := validator.New()
// 			if err := validate.Struct(c); err != nil {
// 				var out []string
// 				var ve validator.ValidationErrors
// 				if errors.As(err, &ve) {
// 					for _, fe := range ve {
// 						switch fe.Field() {
// 						case "Message":
// 							out = append(out, fmt.Sprintf("Messageは1~140文字です"))
// 						case "UserName":
// 							out = append(out, fmt.Sprintf("Messageは1~15文字です"))
// 						}
// 					}
// 				}
// 				http.Error(w, fmt.Sprintf(`{"status":"%s"}`, strings.Join(out, ",")), http.StatusBadRequest)
// 				return
// 			}
// 			mutex.Lock()
// 			comments = append(comments, c)
// 			mutex.Unlock()

// 			w.WriteHeader(http.StatusCreated)
// 			w.Write([]byte(`{"status": "created"}`))

// 		default:
// 			http.Error(w, `{"status":!permits only GET or POST"}`, http.StatusMethodNotAllowed)
		
// 		}
// 	})
// 	http.ListenAndServe(":8888", nil)
// }

// func init() {
// 	pgxDriver = &Driver {
// 		configs: make(map[string]*pgx.ConnConfig),
// 	}
// 	fakeTxConns = make(map[*pgx.Conn]*sql.Tx)
// 	sql.Register("pgx", pgxDriver)
// }
// type User struct {
// 	UserID string
// 	UserName string
// 	CreatedAt time.Time
// }


// type Hendler interface {
// 	ServeHTTP(ResponseWriter, *Request)
// }

// func Hello(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("hello world!"))
// }

// func main() {
// 	ctx, _ := context.WithCancel(context.Background())
// 	db ,err := sql.Open("pgx", "host=localhost port=5432 user=testuser dbname=testdb password=pass sslmode=disable")
// 	err = db.PingContext(ctx)
// 	if nil != err {
// 	}

// 	//　全取得
// 	rows, err := db.QueryContext(ctx, `SELECT user_id, user_name, created_at FROM users ORDER BY user_id;`)
// 	if err != nil {
// 		log.Fatalf("query all users: %v", err)
// 	}
// 	defer rows.Close()
	
// 	var users []*User
// 	for rows.Next() {
// 		var (
// 			userID, userName string
// 			createdAt time.Time
// 		)
// 		if err := rows.Scan(&userID, &userName, &createdAt); err != nil {
// 			log.Fatalf("scan the user:%v", err)
// 		}
// 		users = append(users, &User{
// 			UserID: userID,
// 			UserName: userName,
// 			CreatedAt: createdAt,
// 		})
// 	}
// 	if err := rows.Close(); err != nil {
// 		log.Fatalf("rows close: %v", err)
// 	}
	
// 	if err := rows.Err(); err != nil {
// 		log.Fatalf("scan users: %v", err)
// 	}
// 	// fmt.Printf("%v\n", rows)

// 	var (
// 		userID, userName string
// 		createdAt time.Time
// 	)

// 	row := db.QueryRowContext(ctx, `SELECT user_name, created_at FROM users WHERE user_id = $1;`, userID)
// 	err2 := row.Scan(&userID, &userName, &createdAt)
// 	if err != nil {
// 		log.Fatalf("query row(user_id=%s): %v", userID, err2)
// 	}
	
// 	u := User{
// 		UserID: userID,
// 		UserName: userName,
// 		CreatedAt: createdAt,
// 	}
// 	fmt.Printf("%v\n", u)

// }
