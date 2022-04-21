package users

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"

	"auth_api/sessions"

	"github.com/go-redis/redis/extra/redisotel/v8"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel/codes"
)

type Redis struct {
	client   *redis.Client
	sessions *sessions.Redis
}

func NewRedis() *Redis {
	registerMetrics()
	r := &Redis{
		client: redis.NewClient(&redis.Options{
			Addr:     "192.168.1.45:6379",
			Password: "",
			DB:       0,
		}),
		sessions: sessions.NewSessions(),
	}

	r.client.AddHook(&redisotel.TracingHook{})
	return r
}

type credsUnmarshal struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (r *Redis) getUser(ctx context.Context, username string) bool {
	_, span := tr.Start(ctx, "Check user exists")
	log.Printf("msg=\"checking user exists...\", app=\"auth_api\", trace_id=\"%s\", span_id=\"%s\", level=\"debug\"", span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
	defer span.End()
	// Returns whether the user exists in the database or not
	_, err := r.client.Get(context.Background(), username).Result()
	// true if it does exist, false if theres an error implying it doesn't exist
	if err != nil {
		log.Printf("msg=\"user doesn't exist\", app=\"auth_api\", trace_id=\"%s\", span_id=\"%s\", level=\"debug\"", span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
		return false
	}

	log.Printf("msg=\"user already exists\", app=\"auth_api\", trace_id=\"%s\", span_id=\"%s\", level=\"debug\"", span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
	return true
}

func (r *Redis) CreateUser(ctx context.Context, username string, password string, email string) error {
	_, span := tr.Start(ctx, "Create a user")
	log.Printf("msg=\"creating a user...\", app=\"auth_api\", trace_id=\"%s\", span_id=\"%s\", level=\"debug\"", span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
	defer span.End()
	h := sha1.New()
	h.Write([]byte(password))
	sha1_hash := hex.EncodeToString(h.Sum(nil))

	cu := credsUnmarshal{
		Username: username,
		Password: sha1_hash,
		Email:    email,
	}
	json, err := json.Marshal(cu)
	if err != nil {
		span.RecordError(errors.New("unable to marshal struct into json"))
		span.SetStatus(codes.Error, "Error marshalling struct into json")
		log.Printf("msg=\"error whilst trying to marshal a struct into json\", app=\"auth_api\", err=\"%s\", trace_id=\"%s\", span_id=\"%s\", level=\"error\"", err, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
		return err
	}

	// nil error means it works
	err = r.client.Set(context.Background(), username, json, 0).Err()
	if err != nil {
		span.RecordError(errors.New("error setting value in database"))
		span.SetStatus(codes.Error, "Error setting value in database")
		log.Printf("msg=\"error setting value in the database\", app=\"auth_api\", err=\"%s\", trace_id=\"%s\", span_id=\"%s\", level=\"error\"", err, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
		return err
	}

	log.Printf("msg=\"user was created successfully\", app=\"auth_api\", trace_id=\"%s\", span_id=\"%s\", level=\"debug\"", span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
	return nil
}

func (r *Redis) confirmCreds(username string, password string) bool {
	val, err := r.client.Get(context.Background(), username).Result()
	if err != nil {
		return false
	}

	var cu credsUnmarshal

	json.Unmarshal([]byte(val), &cu)
	return cu.Password == password
}

func (r *Redis) Login(ctx context.Context, username string, password string, sessionID string) (string, error) {
	// this will log the user in, and if successful return a session id and nil, otherwise "" and an error
	ctx, span := tr.Start(ctx, "Login User")
	defer span.End()
	userExists := r.getUser(ctx, username)
	if !userExists {
		// we don't have the user
		return "", errors.New("a user doesn't exist with that username")
	}

	sid := sessionID
	// TODO : refactor below
	if sid != "" {
		sid, _ = r.sessions.GetSession(sid)
		if sid != "" {
			// we have a session id
			return sid, nil
		}
	}

	// we can assume that we do have a user, and we don't currently have a session, so login and create session
	h := sha1.New()
	h.Write([]byte(password))
	sha1_hash := hex.EncodeToString(h.Sum(nil))

	if !r.confirmCreds(username, sha1_hash) {
		// not authd
		return "", errors.New("either the username or password was incorrect")
	}

	// we are authd at this point so create a session
	sid, err := r.sessions.CreateSession(username)
	if err != nil {
		return "", err
	}

	return sid, nil
}
