package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/davidjchavez/url-shortener/internal/model"
	"github.com/davidjchavez/url-shortener/internal/repository"
)

type URLService struct {
	repo    *repository.URLRepository
	baseURL string
}

var errCodeGeneration = errors.New("failed to generate unique code")

const safeURLCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode() string {
	result := make([]byte, 6)
	charsetLength := big.NewInt(int64(len(safeURLCharset)))
	for i := range result {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return ""
		}
		result[i] = safeURLCharset[randomIndex.Int64()]
	}
	return string(result)
}

func NewURLService(repo *repository.URLRepository, baseURL string) *URLService {
	return &URLService{
		repo:    repo,
		baseURL: baseURL,
	}
}

func (s *URLService) CreateShortURL(originalURL string) (*model.CreateURLResponse, error) {
	var code string

	for i := 0; i < 10; i++ {
		code = generateCode()
		exists, err := s.repo.CodeExists(code)
		if err != nil {
			return nil, err
		}
		if !exists {
			break
		}
		if i == 9 {
			return nil, errCodeGeneration
		}
	}

	if code == "" {
		return nil, errCodeGeneration
	}

	url := model.URL{
		Code:        code,
		OriginalURL: originalURL,
	}

	err := s.repo.Create(&url)
	if err != nil {
		return nil, err
	}

	return &model.CreateURLResponse{
		ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, code),
		OriginalURL: originalURL,
		Code:        code,
	}, nil
}

func (s *URLService) GetOriginalURL(code string) (string, error) {
	url, err := s.repo.GetByCode(code)
	if err != nil {
		return "", err
	}

	go s.repo.IncrementClicks(code)

	return url.OriginalURL, nil
}

func (s *URLService) GetStats(code string) (*model.StatsResponse, error) {
	url, err := s.repo.GetByCode(code)
	if err != nil {
		return nil, err
	}
	return &model.StatsResponse{
		Clicks:    url.Clicks,
		CreatedAt: url.CreatedAt,
	}, nil
}
