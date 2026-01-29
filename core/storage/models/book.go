package models

import (
	"errors"
	"strings"
	"time"
)

type Book struct {
	ISBN        string    `json:"isbn"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	PublishedAt time.Time `json:"published_at"`
}

func NewBook(isbn, title, author string, publishedAt time.Time) (*Book, error) {
	if err := validateBook(title, author, isbn); err != nil {
		return nil, err
	}

	return &Book{
		ISBN:        isbn,
		Title:       title,
		Author:      author,
		PublishedAt: publishedAt,
	}, nil
}

func validateBook(title string, author string, isbn string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("title cannot be empty")
	}

	if strings.TrimSpace(author) == "" {
		return errors.New("author cannot be empty")
	}

	if strings.TrimSpace(isbn) == "" {
		return errors.New("ISBN cannot be empty")
	}

	if err := validateISBN(isbn); err != nil {
		return err
	}

	return nil
}

func validateISBN(isbn string) error {
	cleaned := strings.ReplaceAll(isbn, "-", "")
	cleaned = strings.ReplaceAll(cleaned, " ", "")

	switch len(cleaned) {
	case 10:
		return validateISBN10(cleaned)
	case 13:
		return validateISBN13(cleaned)
	default:
		return errors.New("ISBN must be 10 or 13 characters (excluding hyphens)")
	}
}

func validateISBN10(isbn string) error {
	sum := 0
	for i := 0; i < 9; i++ {
		digit := int(isbn[i] - '0')
		if digit < 0 || digit > 9 {
			return errors.New("ISBN-10 must contain only digits (and optionally X as last character)")
		}
		sum += digit * (10 - i)
	}

	lastChar := isbn[9]
	var lastDigit int
	if lastChar == 'X' || lastChar == 'x' {
		lastDigit = 10
	} else {
		lastDigit = int(lastChar - '0')
		if lastDigit < 0 || lastDigit > 9 {
			return errors.New("ISBN-10 must contain only digits (and optionally X as last character)")
		}
	}
	sum += lastDigit

	if sum%11 != 0 {
		return errors.New("invalid ISBN-10 checksum")
	}
	return nil
}

func validateISBN13(isbn string) error {
	sum := 0
	for i := 0; i < 13; i++ {
		digit := int(isbn[i] - '0')
		if digit < 0 || digit > 9 {
			return errors.New("ISBN-13 must contain only digits")
		}
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}

	if sum%10 != 0 {
		return errors.New("invalid ISBN-13 checksum")
	}
	return nil
}
