package services

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gkazioka/clientsapi/app/domain"
	"github.com/gkazioka/clientsapi/app/repositories"
	"github.com/gkazioka/clientsapi/app/repositories/interfaces"
	"github.com/gkazioka/clientsapi/app/types"
)

type UserService struct {
	userRepository interfaces.UserRepository
}

var once sync.Once
var singleton UserService

var cpfPattern = regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)
var cnpjPattern = regexp.MustCompile(`^\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2}$`)

func (s *UserService) Save(ctx context.Context, user domain.User) error {
	if ok, error := validateDocument(user.Code); !ok {
		return error
	}
	error := s.userRepository.Save(ctx, user)
	if error != nil {
		fmt.Fprintf(os.Stderr, "UserService.Save: Error: %v\n", error)
		return error
	}
	return nil
}

func (s UserService) ListAll(ctx context.Context) []domain.User {
	return s.userRepository.ListAll(ctx)
}

func (s UserService) FindUserByCode(ctx context.Context, userCode string) (*domain.User, error) {
	if ok, error := validateDocument(userCode); !ok {
		return nil, error
	}
	userFound, _ := s.userRepository.FindUserByCode(ctx, userCode)
	return userFound, nil
}

func validateCPF(cpf string) bool {
	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")

	if len(cpf) != 11 {
		return false
	}

	if cpf == strings.Repeat(string(cpf[0]), len(cpf)) {
		return false
	}

	calculateDigit := func(digits []int, factor int) int {
		sum := 0
		for _, digit := range digits {
			sum += digit * factor
			factor--
		}
		rest := sum % 11
		if rest < 2 {
			return 0
		}
		return 11 - rest
	}

	digits := make([]int, 11)
	for i := 0; i < 11; i++ {
		digit, err := strconv.Atoi(string(cpf[i]))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	firstDigit := calculateDigit(digits[:9], 10)
	if firstDigit != digits[9] {
		return false
	}

	secondDigit := calculateDigit(digits[:10], 11)
	return secondDigit == digits[10]
}

func validateCNPJ(cnpj string) bool {
	cnpj = strings.ReplaceAll(cnpj, ".", "")
	cnpj = strings.ReplaceAll(cnpj, "/", "")
	cnpj = strings.ReplaceAll(cnpj, "-", "")

	if len(cnpj) != 14 {
		return false
	}

	calculateDigit := func(digits []int, factor []int) int {
		sum := 0
		for i := 0; i < len(factor); i++ {
			sum += digits[i] * factor[i]
		}
		rest := sum % 11
		if rest < 2 {
			return 0
		}
		return 11 - rest
	}

	digits := make([]int, 14)
	for i := 0; i < 14; i++ {
		digit, err := strconv.Atoi(string(cnpj[i]))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	factor1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	firstDigit := calculateDigit(digits[:12], factor1)
	if firstDigit != digits[12] {
		return false
	}

	factor2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	secondDigit := calculateDigit(digits[:13], factor2)
	return secondDigit == digits[13]
}

func validateDocument(document string) (bool, error) {

	if cpfPattern.MatchString(document) {
		if validateCPF(document) {
			return true, nil
		}
		return false, types.ErrorInvalidCpf
	} else if cnpjPattern.MatchString(document) {
		if validateCNPJ(document) {
			return true, nil
		}
		return false, types.ErrorInvalidCnpj
	}
	return false, types.ErrorInvalidDocument
}

func newUserService(repositoryType string) UserService {
	if repositoryType == "memory" {
		return UserService{&repositories.UserRepositoryMemory{}}
	}
	if repositoryType == "postgres" {
		return UserService{repositories.NewUserRepositoryPostgres()}
	} else {
		return UserService{&repositories.UserRepositoryMemory{}}
	}
}

func GetInstance(repositoryType string) UserService {
	once.Do(func() {
		singleton = newUserService(repositoryType)
	})
	return singleton
}
