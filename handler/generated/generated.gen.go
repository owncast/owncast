// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package generated

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Disconnect inbound stream
	// (GET /admin/disconnect)
	DisconnectInboundConnection(w http.ResponseWriter, r *http.Request)
	// Get the current server config
	// (GET /admin/serverconfig)
	GetServerConfig(w http.ResponseWriter, r *http.Request)
	// Get current inboard broadcaster
	// (GET /admin/status)
	GetAdminStatus(w http.ResponseWriter, r *http.Request)
	// Get viewer count over time
	// (GET /admin/viewersOverTime)
	GetViewersOverTime(w http.ResponseWriter, r *http.Request, params GetViewersOverTimeParams)
	// Gets a list of chat messages
	// (GET /chat)
	GetChatList(w http.ResponseWriter, r *http.Request)

	// (OPTIONS /chat/register)
	OptionsChatRegister(w http.ResponseWriter, r *http.Request)
	// Registers an anonymous chat user
	// (POST /chat/register)
	RegisterAnonymousChatUser(w http.ResponseWriter, r *http.Request, params RegisterAnonymousChatUserParams)
	// Get the web config
	// (GET /config)
	GetConfig(w http.ResponseWriter, r *http.Request)
	// Get list of custom emojis supported in the chat
	// (GET /emoji)
	GetEmoji(w http.ResponseWriter, r *http.Request)
	// Gets the list of followers
	// (GET /followers)
	GetFollowers(w http.ResponseWriter, r *http.Request, params GetFollowersParams)
	// Send a system message to the chat
	// (POST /integrations/chat/system)
	SendSystemMessage(w http.ResponseWriter, r *http.Request)
	// Save video playback metrics for future video health recording
	// (POST /metrics/playback)
	PostMetricsPlayback(w http.ResponseWriter, r *http.Request)
	// Register for notifications
	// (POST /notifications/register)
	PostNotificationsRegister(w http.ResponseWriter, r *http.Request, params PostNotificationsRegisterParams)
	// Tell the backend you're an active viewer
	// (GET /ping)
	Ping(w http.ResponseWriter, r *http.Request)
	// Request remote follow
	// (POST /remotefollow)
	RemoteFollow(w http.ResponseWriter, r *http.Request)
	// Get all social platforms
	// (GET /socialplatforms)
	GetSocialPlatforms(w http.ResponseWriter, r *http.Request)
	// Get the status of the server
	// (GET /status)
	GetStatus(w http.ResponseWriter, r *http.Request)
	// Get a list of video variants available
	// (GET /video/variants)
	GetVideoVariants(w http.ResponseWriter, r *http.Request)
	// Get the YP protocol data
	// (GET /yp)
	GetYP(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Disconnect inbound stream
// (GET /admin/disconnect)
func (_ Unimplemented) DisconnectInboundConnection(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get the current server config
// (GET /admin/serverconfig)
func (_ Unimplemented) GetServerConfig(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get current inboard broadcaster
// (GET /admin/status)
func (_ Unimplemented) GetAdminStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get viewer count over time
// (GET /admin/viewersOverTime)
func (_ Unimplemented) GetViewersOverTime(w http.ResponseWriter, r *http.Request, params GetViewersOverTimeParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Gets a list of chat messages
// (GET /chat)
func (_ Unimplemented) GetChatList(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (OPTIONS /chat/register)
func (_ Unimplemented) OptionsChatRegister(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Registers an anonymous chat user
// (POST /chat/register)
func (_ Unimplemented) RegisterAnonymousChatUser(w http.ResponseWriter, r *http.Request, params RegisterAnonymousChatUserParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get the web config
// (GET /config)
func (_ Unimplemented) GetConfig(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get list of custom emojis supported in the chat
// (GET /emoji)
func (_ Unimplemented) GetEmoji(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Gets the list of followers
// (GET /followers)
func (_ Unimplemented) GetFollowers(w http.ResponseWriter, r *http.Request, params GetFollowersParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Send a system message to the chat
// (POST /integrations/chat/system)
func (_ Unimplemented) SendSystemMessage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Save video playback metrics for future video health recording
// (POST /metrics/playback)
func (_ Unimplemented) PostMetricsPlayback(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Register for notifications
// (POST /notifications/register)
func (_ Unimplemented) PostNotificationsRegister(w http.ResponseWriter, r *http.Request, params PostNotificationsRegisterParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Tell the backend you're an active viewer
// (GET /ping)
func (_ Unimplemented) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Request remote follow
// (POST /remotefollow)
func (_ Unimplemented) RemoteFollow(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get all social platforms
// (GET /socialplatforms)
func (_ Unimplemented) GetSocialPlatforms(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get the status of the server
// (GET /status)
func (_ Unimplemented) GetStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get a list of video variants available
// (GET /video/variants)
func (_ Unimplemented) GetVideoVariants(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get the YP protocol data
// (GET /yp)
func (_ Unimplemented) GetYP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// DisconnectInboundConnection operation middleware
func (siw *ServerInterfaceWrapper) DisconnectInboundConnection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DisconnectInboundConnection(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetServerConfig operation middleware
func (siw *ServerInterfaceWrapper) GetServerConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetServerConfig(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetAdminStatus operation middleware
func (siw *ServerInterfaceWrapper) GetAdminStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetAdminStatus(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetViewersOverTime operation middleware
func (siw *ServerInterfaceWrapper) GetViewersOverTime(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	ctx = context.WithValue(ctx, BasicAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetViewersOverTimeParams

	// ------------- Optional query parameter "windowStart" -------------

	err = runtime.BindQueryParameter("form", true, false, "windowStart", r.URL.Query(), &params.WindowStart)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "windowStart", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetViewersOverTime(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetChatList operation middleware
func (siw *ServerInterfaceWrapper) GetChatList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetChatList(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// OptionsChatRegister operation middleware
func (siw *ServerInterfaceWrapper) OptionsChatRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.OptionsChatRegister(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// RegisterAnonymousChatUser operation middleware
func (siw *ServerInterfaceWrapper) RegisterAnonymousChatUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params RegisterAnonymousChatUserParams

	headers := r.Header

	// ------------- Optional header parameter "X-Forwarded-User" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Forwarded-User")]; found {
		var XForwardedUser string
		n := len(valueList)
		if n != 1 {
			siw.ErrorHandlerFunc(w, r, &TooManyValuesForParamError{ParamName: "X-Forwarded-User", Count: n})
			return
		}

		err = runtime.BindStyledParameterWithOptions("simple", "X-Forwarded-User", valueList[0], &XForwardedUser, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
		if err != nil {
			siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "X-Forwarded-User", Err: err})
			return
		}

		params.XForwardedUser = &XForwardedUser

	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RegisterAnonymousChatUser(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetConfig operation middleware
func (siw *ServerInterfaceWrapper) GetConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetConfig(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetEmoji operation middleware
func (siw *ServerInterfaceWrapper) GetEmoji(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetEmoji(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetFollowers operation middleware
func (siw *ServerInterfaceWrapper) GetFollowers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetFollowersParams

	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", r.URL.Query(), &params.Offset)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "offset", Err: err})
		return
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetFollowers(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// SendSystemMessage operation middleware
func (siw *ServerInterfaceWrapper) SendSystemMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SendSystemMessage(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostMetricsPlayback operation middleware
func (siw *ServerInterfaceWrapper) PostMetricsPlayback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostMetricsPlayback(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostNotificationsRegister operation middleware
func (siw *ServerInterfaceWrapper) PostNotificationsRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params PostNotificationsRegisterParams

	// ------------- Required query parameter "accessToken" -------------

	if paramValue := r.URL.Query().Get("accessToken"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "accessToken"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "accessToken", r.URL.Query(), &params.AccessToken)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "accessToken", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostNotificationsRegister(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// Ping operation middleware
func (siw *ServerInterfaceWrapper) Ping(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Ping(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// RemoteFollow operation middleware
func (siw *ServerInterfaceWrapper) RemoteFollow(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RemoteFollow(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetSocialPlatforms operation middleware
func (siw *ServerInterfaceWrapper) GetSocialPlatforms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetSocialPlatforms(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetStatus operation middleware
func (siw *ServerInterfaceWrapper) GetStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetStatus(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetVideoVariants operation middleware
func (siw *ServerInterfaceWrapper) GetVideoVariants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetVideoVariants(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetYP operation middleware
func (siw *ServerInterfaceWrapper) GetYP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetYP(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/admin/disconnect", wrapper.DisconnectInboundConnection)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/admin/serverconfig", wrapper.GetServerConfig)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/admin/status", wrapper.GetAdminStatus)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/admin/viewersOverTime", wrapper.GetViewersOverTime)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/chat", wrapper.GetChatList)
	})
	r.Group(func(r chi.Router) {
		r.Options(options.BaseURL+"/chat/register", wrapper.OptionsChatRegister)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/chat/register", wrapper.RegisterAnonymousChatUser)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/config", wrapper.GetConfig)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/emoji", wrapper.GetEmoji)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/followers", wrapper.GetFollowers)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/integrations/chat/system", wrapper.SendSystemMessage)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/metrics/playback", wrapper.PostMetricsPlayback)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/notifications/register", wrapper.PostNotificationsRegister)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/ping", wrapper.Ping)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/remotefollow", wrapper.RemoteFollow)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/socialplatforms", wrapper.GetSocialPlatforms)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/status", wrapper.GetStatus)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/video/variants", wrapper.GetVideoVariants)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/yp", wrapper.GetYP)
	})

	return r
}
