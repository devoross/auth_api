package users

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-playground/validator"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type loginPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type signupPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

type loginResponse struct {
	SID      string `json:"sid"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

var tr = otel.Tracer("AuthAPI")

func (r *Redis) LoginHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var span trace.Span

	ctx, span = tr.Start(ctx, "Perform Login Request")
	defer span.End()

	// make sure it's POST
	if req.Method != "POST" {
		_, childSpan := tr.Start(ctx, "Response")
		span.RecordError(errors.New("http method not allowed"))
		span.SetStatus(codes.Error, "incorrect http method provided")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		log.Printf("msg=\"incorrect method was used on http request\", method=\"%s\", remote_addr=\"%s\", request_uri=\"%s\", trace_id=\"%s\", span_id=\"%s\"", req.Method, req.RemoteAddr, req.RequestURI, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
		childSpan.End()
		return
	}

	// use the Login function, which will return either an error, or the session id
	var lp loginPayload

	val := req.Header.Get("x-session-id")

	_, childSpan := tr.Start(ctx, "Decode JSON Body")
	err := json.NewDecoder(req.Body).Decode(&lp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("msg=\"bad request provided\", method=\"%s\", remote_addr=\"%s\", request_uri=\"%s\", trace_id=\"%s\", span_id=\"%s\"", req.Method, req.RemoteAddr, req.RequestURI, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
		span.RecordError(errors.New("bad request provided"))
		span.SetStatus(codes.Error, "bad payload provided on the request")
		span.AddEvent("tag", "helloworld")
		childSpan.End()
		return
	}
	childSpan.End()

	_, childSpan = tr.Start(ctx, "Validate JSON Body")
	validate := validator.New()
	err = validate.Struct(lp)
	if err != nil {
		http.Error(w, "bad payload provided", http.StatusBadRequest)
		log.Printf("msg=\"bad request provided\", method=\"%s\", remote_addr=\"%s\", request_uri=\"%s\", trace_id=\"%s\", span_id=\"%s\"", req.Method, req.RemoteAddr, req.RequestURI, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
		span.RecordError(errors.New("bad request provided"))
		span.SetStatus(codes.Error, "bad payload provided on the request")
		childSpan.End()
		return
	}
	childSpan.End()

	// have this return the username, email and the sid
	sid, err := r.Login(ctx, lp.Username, lp.Password, val)
	if err != nil {
		_, childSpan = tr.Start(ctx, "Response")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		log.Printf("msg=\"bad credentials provided\", method=\"%s\", remote_addr=\"%s\", request_uri=\"%s\", trace_id=\"%s\", span_id=\"%s\"", req.Method, req.RemoteAddr, req.RequestURI, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
		span.RecordError(err)
		span.SetStatus(codes.Error, "Incorrect credentials provided")
		childSpan.End()
		return
	}

	_, childSpan = tr.Start(ctx, "Response")
	defer childSpan.End()
	a := loginResponse{
		SID:      sid,
		Email:    "test@test.com",
		Username: lp.Username,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
	log.Printf("msg=\"request completed successfully\", method=\"%s\", remote_addr=\"%s\", request_uri=\"%s\", trace_id=\"%s\", span_id=\"%s\"", req.Method, req.RemoteAddr, req.RequestURI, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
}

func (r *Redis) SignupHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var span trace.Span

	ctx, span = tr.Start(ctx, "Perform Signup Request")
	defer span.End()

	// make sure it's POST
	if req.Method != "POST" {
		_, childSpan := tr.Start(ctx, "Response")
		span.RecordError(errors.New("http method not allowed"))
		span.SetStatus(codes.Error, "incorrect http method provided")
		http.Error(w, "http method not allowed", http.StatusMethodNotAllowed)
		log.Printf("msg=\"incorrect method was used on http request\", method=\"%s\", remote_addr=\"%s\", request_uri=\"%s\", trace_id=\"%s\", span_id=\"%s\"", req.Method, req.RemoteAddr, req.RequestURI, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
		childSpan.End()
		return
	}

	// use the Login function, which will return either an error, or the session id
	var sp signupPayload

	_, childSpan := tr.Start(ctx, "Decode JSON Body")
	err := json.NewDecoder(req.Body).Decode(&sp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("msg=\"bad request provided\", method=\"%s\", remote_addr=\"%s\", request_uri=\"%s\", trace_id=\"%s\", span_id=\"%s\"", req.Method, req.RemoteAddr, req.RequestURI, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
		span.RecordError(errors.New("bad request provided"))
		span.SetStatus(codes.Error, "bad payload provided on the request")
		childSpan.End()
		return
	}
	childSpan.End()

	_, childSpan = tr.Start(ctx, "Validate JSON Body")
	validate := validator.New()
	err = validate.Struct(sp)
	if err != nil {
		http.Error(w, "bad payload provided", http.StatusBadRequest)
		log.Printf("msg=\"bad request provided\", method=\"%s\", remote_addr=\"%s\", request_uri=\"%s\", trace_id=\"%s\", span_id=\"%s\"", req.Method, req.RemoteAddr, req.RequestURI, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
		span.RecordError(errors.New("bad request provided"))
		span.SetStatus(codes.Error, "bad payload provided on the request")
		childSpan.End()
		return
	}
	childSpan.End()

	// does the user exist
	exists := r.getUser(ctx, sp.Username)
	if exists {
		_, childSpan = tr.Start(ctx, "Response")
		http.Error(w, "username already exists", http.StatusBadRequest)
		log.Printf("msg=\"username already exists\", method=\"%s\", remote_addr=\"%s\", request_uri=\"%s\", trace_id=\"%s\", span_id=\"%s\"", req.Method, req.RemoteAddr, req.RequestURI, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
		span.RecordError(err)
		span.SetStatus(codes.Error, "Username already exists")
		childSpan.End()
		return
	}

	err = r.CreateUser(ctx, sp.Username, sp.Password, sp.Email)
	if err != nil {
		_, childSpan = tr.Start(ctx, "Response")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Printf("msg=\"internal server error\", method=\"%s\", remote_addr=\"%s\", request_uri=\"%s\", trace_id=\"%s\", span_id=\"%s\"", req.Method, req.RemoteAddr, req.RequestURI, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
		span.RecordError(err)
		span.SetStatus(codes.Error, "Internal Server Error")
		childSpan.End()
		return
	}

	_, childSpan = tr.Start(ctx, "Response")
	log.Printf("msg=\"request completed successfully\", method=\"%s\", remote_addr=\"%s\", request_uri=\"%s\", trace_id=\"%s\", span_id=\"%s\"", req.Method, req.RemoteAddr, req.RequestURI, span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
	defer childSpan.End()
}
