package securityext

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Secret Manager

type Secret struct {
	Name      string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SecretManager struct {
	secrets map[string]*Secret
	key     []byte
	mu      sync.RWMutex
}

func NewSecretManager(encryptionKey string) *SecretManager {
	hash := sha256.Sum256([]byte(encryptionKey))
	return &SecretManager{
		secrets: make(map[string]*Secret),
		key:     hash[:],
	}
}

func (sm *SecretManager) Set(name, value string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	encrypted, err := sm.encrypt(value)
	if err != nil {
		return err
	}

	now := time.Now()
	if existing, ok := sm.secrets[name]; ok {
		existing.Value = encrypted
		existing.UpdatedAt = now
	} else {
		sm.secrets[name] = &Secret{
			Name:      name,
			Value:     encrypted,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return nil
}

func (sm *SecretManager) Get(name string) (string, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	secret, ok := sm.secrets[name]
	if !ok {
		return "", false
	}

	decrypted, err := sm.decrypt(secret.Value)
	if err != nil {
		return "", false
	}
	return decrypted, true
}

func (sm *SecretManager) Delete(name string) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, ok := sm.secrets[name]; ok {
		delete(sm.secrets, name)
		return true
	}
	return false
}

func (sm *SecretManager) List() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	names := make([]string, 0, len(sm.secrets))
	for name := range sm.secrets {
		names = append(names, name)
	}
	return names
}

func (sm *SecretManager) Exists(name string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, ok := sm.secrets[name]
	return ok
}

func (sm *SecretManager) encrypt(value string) (string, error) {
	block, err := aes.NewCipher(sm.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(value), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (sm *SecretManager) decrypt(encrypted string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(sm.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Input Sanitizer

type Sanitizer struct {
	patterns map[string]*regexp.Regexp
	mu       sync.RWMutex
}

func NewSanitizer() *Sanitizer {
	s := &Sanitizer{
		patterns: make(map[string]*regexp.Regexp),
	}

	s.patterns["html"] = regexp.MustCompile(`<[^>]*>`)
	s.patterns["sql"] = regexp.MustCompile(`['";\\]`)
	s.patterns["xss"] = regexp.MustCompile(`<script[^>]*>.*?</script>`)
	s.patterns["path"] = regexp.MustCompile(`\.\./`)
	s.patterns["email"] = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return s
}

func (s *Sanitizer) SanitizeHTML(input string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.patterns["html"].ReplaceAllString(input, "")
}

func (s *Sanitizer) SanitizeSQL(input string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.patterns["sql"].ReplaceAllString(input, "")
}

func (s *Sanitizer) SanitizeXSS(input string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.patterns["xss"].ReplaceAllString(input, "")
}

func (s *Sanitizer) SanitizePath(input string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.patterns["path"].ReplaceAllString(input, "")
}

func (s *Sanitizer) ValidateEmail(email string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.patterns["email"].MatchString(email)
}

func (s *Sanitizer) SanitizeAll(input string) string {
	result := input
	result = s.SanitizeHTML(result)
	result = s.SanitizeXSS(result)
	result = s.SanitizePath(result)
	return result
}

// Hash

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func VerifyPassword(password, hash string) bool {
	return HashPassword(password) == hash
}

// Token Generator

func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Validator

type Validator struct {
	rules map[string]func(string) error
	mu    sync.RWMutex
}

func NewValidator() *Validator {
	return &Validator{
		rules: make(map[string]func(string) error),
	}
}

func (v *Validator) AddRule(name string, rule func(string) error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.rules[name] = rule
}

func (v *Validator) Validate(name, value string) error {
	v.mu.RLock()
	defer v.mu.RUnlock()

	rule, ok := v.rules[name]
	if !ok {
		return fmt.Errorf("rule not found: %s", name)
	}
	return rule(value)
}

func (v *Validator) ValidateAll(values map[string]string) []error {
	var errors []error

	for name, value := range values {
		if err := v.Validate(name, value); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// Default Validator Rules

func RequiredRule(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("value is required")
	}
	return nil
}

func MinLengthRule(min int) func(string) error {
	return func(value string) error {
		if len(value) < min {
			return fmt.Errorf("value must be at least %d characters", min)
		}
		return nil
	}
}

func MaxLengthRule(max int) func(string) error {
	return func(value string) error {
		if len(value) > max {
			return fmt.Errorf("value must be at most %d characters", max)
		}
		return nil
	}
}

func PatternRule(pattern string) func(string) error {
	return func(value string) error {
		matched, err := regexp.MatchString(pattern, value)
		if err != nil {
			return err
		}
		if !matched {
			return fmt.Errorf("value does not match pattern")
		}
		return nil
	}
}
