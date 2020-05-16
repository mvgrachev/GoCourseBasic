package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/sirupsen/logrus"
        "go.mongodb.org/mongo-driver/mongo"	
	"go.mongodb.org/mongo-driver/mongo/options"
	"fmt"
)

// NewLogger - Создаёт новый логгер
func NewLogger() *logrus.Logger {
        lg := logrus.New()
        lg.SetReportCaller(false)
        lg.SetFormatter(&logrus.TextFormatter{})
        lg.SetLevel(logrus.DebugLevel)
        return lg
}

// TestPost - объект поста
type TestPost struct {
        Id string `json:"id"`
        Title   string `json:"title"`
        Date    string `json:"date"`
        Summary string `json:"summary"`
        Body    interface{} `json:"body"`
        Status  int `json:"status"`
}

func (post *TestPost) Read(ctx context.Context, db *mongo.Database) error {
	return nil
}

// TestServer - объект сервера
type TestServer struct {
	lg            *logrus.Logger
	db            *mongo.Database
	post          *TestPost
}

func TestPostPostHandler( t *testing.T ) {
	lg := NewLogger()

        client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
        if err != nil {
                lg.WithError(err).Fatal("can't connect to db")
        }
        ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
        defer cancel()
        err = client.Connect(ctx)
        if err != nil {
                lg.WithError(err).Fatal("can't connect to db")
        }

        db := client.Database("blog")

	testPost := &TestPost{
	}

	testMessage := map[string]interface{}{
		"id": "12345",
		"title": "Test title",
		"date": "16.05.2020",
		"summary": "Test for post",
		"body": []string{
			"Head1",
			"Head2",
			"Head3",
		},
	}

	bytesRepresentation, err := json.Marshal(testMessage)
	if err != nil {
		log.Fatalln(err)
	}

	req, err := http.NewRequest("POST", "/api/v1/posts", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	serv := &TestServer{
		lg: lg,
		db: db,
		post: testPost,
	}

	// Мы создаем ResponseRecorder(реализует интерфейс http.ResponseWriter)
	// и используем его для получения ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serv.postPostHandler)

	// Наш хендлер соответствует интерфейсу http.Handler, а значит
	// мы можем использовать ServeHTTP и напрямую указать 
	// Request и ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Проверяем код
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var result map[string]interface{}

	json.NewDecoder(rr.Body).Decode(&result)
	fmt.Println(result)

	// Проверяем тело ответа
	//expected := `{"alive": true}`
	//if rr.Body.String() != expected {
	//	t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	//}
}
