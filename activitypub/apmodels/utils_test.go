package apmodels

import (
	"errors"
	"net/url"
	"testing"

	"github.com/go-fed/activity/streams"
)

func TestGetIRIFromActorPropertyWithNil(t *testing.T) {
	_, err := GetIRIFromActorProperty(nil)
	if err == nil {
		t.Error("GetIRIFromActorProperty(nil) should return error")
	}
	if !errors.Is(err, ErrMissingIRI) {
		t.Errorf("GetIRIFromActorProperty(nil) error = %v, want ErrMissingIRI", err)
	}
}

func TestGetIRIFromActorPropertyWithEmpty(t *testing.T) {
	actor := streams.NewActivityStreamsActorProperty()
	_, err := GetIRIFromActorProperty(actor)
	if err == nil {
		t.Error("GetIRIFromActorProperty with empty actor should return error")
	}
	if !errors.Is(err, ErrMissingIRI) {
		t.Errorf("GetIRIFromActorProperty error = %v, want ErrMissingIRI", err)
	}
}

func TestGetIRIFromActorPropertyWithValidIRI(t *testing.T) {
	actor := streams.NewActivityStreamsActorProperty()
	iri, _ := url.Parse("https://example.com/user/test")
	actor.AppendIRI(iri)

	result, err := GetIRIFromActorProperty(actor)
	if err != nil {
		t.Errorf("GetIRIFromActorProperty with valid IRI should not return error, got %v", err)
	}
	if result.String() != "https://example.com/user/test" {
		t.Errorf("GetIRIFromActorProperty = %v, want %v", result.String(), "https://example.com/user/test")
	}
}

func TestGetIRIStringFromActorPropertyWithValidIRI(t *testing.T) {
	actor := streams.NewActivityStreamsActorProperty()
	iri, _ := url.Parse("https://example.com/user/test")
	actor.AppendIRI(iri)

	result, err := GetIRIStringFromActorProperty(actor)
	if err != nil {
		t.Errorf("GetIRIStringFromActorProperty with valid IRI should not return error, got %v", err)
	}
	if result != "https://example.com/user/test" {
		t.Errorf("GetIRIStringFromActorProperty = %v, want %v", result, "https://example.com/user/test")
	}
}

func TestGetIRIFromObjectPropertyWithNil(t *testing.T) {
	_, err := GetIRIFromObjectProperty(nil)
	if err == nil {
		t.Error("GetIRIFromObjectProperty(nil) should return error")
	}
	if !errors.Is(err, ErrMissingIRI) {
		t.Errorf("GetIRIFromObjectProperty(nil) error = %v, want ErrMissingIRI", err)
	}
}

func TestGetIRIFromObjectPropertyWithEmpty(t *testing.T) {
	object := streams.NewActivityStreamsObjectProperty()
	_, err := GetIRIFromObjectProperty(object)
	if err == nil {
		t.Error("GetIRIFromObjectProperty with empty object should return error")
	}
	if !errors.Is(err, ErrMissingIRI) {
		t.Errorf("GetIRIFromObjectProperty error = %v, want ErrMissingIRI", err)
	}
}

func TestGetIRIFromObjectPropertyWithValidIRI(t *testing.T) {
	object := streams.NewActivityStreamsObjectProperty()
	iri, _ := url.Parse("https://example.com/post/123")
	object.AppendIRI(iri)

	result, err := GetIRIFromObjectProperty(object)
	if err != nil {
		t.Errorf("GetIRIFromObjectProperty with valid IRI should not return error, got %v", err)
	}
	if result.String() != "https://example.com/post/123" {
		t.Errorf("GetIRIFromObjectProperty = %v, want %v", result.String(), "https://example.com/post/123")
	}
}

func TestGetIRIStringFromObjectPropertyWithValidIRI(t *testing.T) {
	object := streams.NewActivityStreamsObjectProperty()
	iri, _ := url.Parse("https://example.com/post/123")
	object.AppendIRI(iri)

	result, err := GetIRIStringFromObjectProperty(object)
	if err != nil {
		t.Errorf("GetIRIStringFromObjectProperty with valid IRI should not return error, got %v", err)
	}
	if result != "https://example.com/post/123" {
		t.Errorf("GetIRIStringFromObjectProperty = %v, want %v", result, "https://example.com/post/123")
	}
}

func TestGetIRIFromJSONLDIdPropertyWithNil(t *testing.T) {
	_, err := GetIRIFromJSONLDIdProperty(nil)
	if err == nil {
		t.Error("GetIRIFromJSONLDIdProperty(nil) should return error")
	}
	if !errors.Is(err, ErrMissingIRI) {
		t.Errorf("GetIRIFromJSONLDIdProperty(nil) error = %v, want ErrMissingIRI", err)
	}
}

func TestGetIRIFromJSONLDIdPropertyWithValidIRI(t *testing.T) {
	id := streams.NewJSONLDIdProperty()
	iri, _ := url.Parse("https://example.com/activity/456")
	id.SetIRI(iri)

	result, err := GetIRIFromJSONLDIdProperty(id)
	if err != nil {
		t.Errorf("GetIRIFromJSONLDIdProperty with valid IRI should not return error, got %v", err)
	}
	if result.String() != "https://example.com/activity/456" {
		t.Errorf("GetIRIFromJSONLDIdProperty = %v, want %v", result.String(), "https://example.com/activity/456")
	}
}

func TestGetIRIStringFromJSONLDIdPropertyWithValidIRI(t *testing.T) {
	id := streams.NewJSONLDIdProperty()
	iri, _ := url.Parse("https://example.com/activity/456")
	id.SetIRI(iri)

	result, err := GetIRIStringFromJSONLDIdProperty(id)
	if err != nil {
		t.Errorf("GetIRIStringFromJSONLDIdProperty with valid IRI should not return error, got %v", err)
	}
	if result != "https://example.com/activity/456" {
		t.Errorf("GetIRIStringFromJSONLDIdProperty = %v, want %v", result, "https://example.com/activity/456")
	}
}

func TestGetIRIFromJSONLDIdPropertyWithNoIRI(t *testing.T) {
	id := streams.NewJSONLDIdProperty()
	// Set a non-IRI value (using Set instead of SetIRI)
	iri, _ := url.Parse("https://example.com/activity/456")
	id.Set(iri)

	// When using Set() the IRI should still be retrievable via GetIRI()
	result, err := GetIRIFromJSONLDIdProperty(id)
	if err != nil {
		t.Errorf("GetIRIFromJSONLDIdProperty should work with Set(), got error %v", err)
	}
	if result == nil {
		t.Error("GetIRIFromJSONLDIdProperty result should not be nil")
	}
}

func TestGetPublicKeyPemWithNil(t *testing.T) {
	_, err := GetPublicKeyPem(nil)
	if err == nil {
		t.Error("GetPublicKeyPem(nil) should return error")
	}
}

func TestGetPublicKeyPemWithValidKey(t *testing.T) {
	publicKey := streams.NewW3IDSecurityV1PublicKey()
	pemProp := streams.NewW3IDSecurityV1PublicKeyPemProperty()
	pemProp.Set("-----BEGIN PUBLIC KEY-----\ntest\n-----END PUBLIC KEY-----")
	publicKey.SetW3IDSecurityV1PublicKeyPem(pemProp)

	result, err := GetPublicKeyPem(publicKey)
	if err != nil {
		t.Errorf("GetPublicKeyPem with valid key should not return error, got %v", err)
	}
	if result != "-----BEGIN PUBLIC KEY-----\ntest\n-----END PUBLIC KEY-----" {
		t.Errorf("GetPublicKeyPem = %v, want PEM string", result)
	}
}

func TestGetPublicKeyPemWithNoPem(t *testing.T) {
	publicKey := streams.NewW3IDSecurityV1PublicKey()
	// Don't set PEM property

	_, err := GetPublicKeyPem(publicKey)
	if err == nil {
		t.Error("GetPublicKeyPem with no PEM should return error")
	}
}

// Test that safe accessors don't panic on nil/empty inputs
func TestSafeAccessorsNoPanic(t *testing.T) {
	// These should not panic
	_, _ = GetIRIFromActorProperty(nil)
	_, _ = GetIRIStringFromActorProperty(nil)
	_, _ = GetIRIFromObjectProperty(nil)
	_, _ = GetIRIStringFromObjectProperty(nil)
	_, _ = GetIRIFromJSONLDIdProperty(nil)
	_, _ = GetIRIStringFromJSONLDIdProperty(nil)
	_, _ = GetPublicKeyPem(nil)
	_ = GetImageFromIcon(nil)
	_ = GetHostnameFromJSONLDId(nil)
	_ = IsFirstObjectActivityStreamsPerson(nil)

	// Empty properties should also not panic
	_, _ = GetIRIFromActorProperty(streams.NewActivityStreamsActorProperty())
	_, _ = GetIRIFromObjectProperty(streams.NewActivityStreamsObjectProperty())
	_ = GetImageFromIcon(streams.NewActivityStreamsIconProperty())
	_ = IsFirstObjectActivityStreamsPerson(streams.NewActivityStreamsObjectProperty())
}

func TestGetImageFromIconWithNil(t *testing.T) {
	result := GetImageFromIcon(nil)
	if result != nil {
		t.Error("GetImageFromIcon(nil) should return nil")
	}
}

func TestGetImageFromIconWithEmpty(t *testing.T) {
	icon := streams.NewActivityStreamsIconProperty()
	result := GetImageFromIcon(icon)
	if result != nil {
		t.Error("GetImageFromIcon with empty icon should return nil")
	}
}

func TestGetImageFromIconWithValidImage(t *testing.T) {
	icon := streams.NewActivityStreamsIconProperty()
	image := streams.NewActivityStreamsImage()
	urlProp := streams.NewActivityStreamsUrlProperty()
	imageURL, _ := url.Parse("https://example.com/avatar.png")
	urlProp.AppendIRI(imageURL)
	image.SetActivityStreamsUrl(urlProp)
	icon.AppendActivityStreamsImage(image)

	result := GetImageFromIcon(icon)
	if result == nil {
		t.Error("GetImageFromIcon with valid image should not return nil")
	}
	if result.String() != "https://example.com/avatar.png" {
		t.Errorf("GetImageFromIcon = %v, want %v", result.String(), "https://example.com/avatar.png")
	}
}

func TestGetHostnameFromJSONLDIdWithNil(t *testing.T) {
	result := GetHostnameFromJSONLDId(nil)
	if result != "" {
		t.Errorf("GetHostnameFromJSONLDId(nil) = %v, want empty string", result)
	}
}

func TestGetHostnameFromJSONLDIdWithValidId(t *testing.T) {
	id := streams.NewJSONLDIdProperty()
	iri, _ := url.Parse("https://example.com/user/test")
	id.SetIRI(iri)

	result := GetHostnameFromJSONLDId(id)
	if result != "example.com" {
		t.Errorf("GetHostnameFromJSONLDId = %v, want %v", result, "example.com")
	}
}

func TestIsFirstObjectActivityStreamsPersonWithNil(t *testing.T) {
	result := IsFirstObjectActivityStreamsPerson(nil)
	if result {
		t.Error("IsFirstObjectActivityStreamsPerson(nil) should return false")
	}
}

func TestIsFirstObjectActivityStreamsPersonWithEmpty(t *testing.T) {
	object := streams.NewActivityStreamsObjectProperty()
	result := IsFirstObjectActivityStreamsPerson(object)
	if result {
		t.Error("IsFirstObjectActivityStreamsPerson with empty object should return false")
	}
}

func TestIsFirstObjectActivityStreamsPersonWithPerson(t *testing.T) {
	object := streams.NewActivityStreamsObjectProperty()
	person := streams.NewActivityStreamsPerson()
	object.AppendActivityStreamsPerson(person)

	result := IsFirstObjectActivityStreamsPerson(object)
	if !result {
		t.Error("IsFirstObjectActivityStreamsPerson with person should return true")
	}
}

func TestIsFirstObjectActivityStreamsPersonWithNonPerson(t *testing.T) {
	object := streams.NewActivityStreamsObjectProperty()
	note := streams.NewActivityStreamsNote()
	object.AppendActivityStreamsNote(note)

	result := IsFirstObjectActivityStreamsPerson(object)
	if result {
		t.Error("IsFirstObjectActivityStreamsPerson with note should return false")
	}
}
