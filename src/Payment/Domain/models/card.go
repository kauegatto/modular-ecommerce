package models

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CardNumber struct {
	Value string
}

func NewCardNumber(pan string) (*CardNumber, error) {
	pan = strings.ReplaceAll(pan, " ", "")
	pan = strings.ReplaceAll(pan, "-", "")

	if !regexp.MustCompile(`^\d{13,19}$`).MatchString(pan) {
		return nil, errors.New("invalid card number format")
	}

	if !isLuhnValid(pan) {
		return nil, errors.New("invalid card number checksum")
	}

	return &CardNumber{Value: pan}, nil
}

func (cn *CardNumber) String() string {
	return fmt.Sprintf("****-%s", cn.LastFourDigits())
}

func (cn *CardNumber) LastFourDigits() string {
	return cn.Value[len(cn.Value)-4:]
}

type SecurityCode struct {
	Value string
}

func (sc *SecurityCode) String() string {
	return "***"
}

func NewSecurityCode(code string) (*SecurityCode, error) {
	if !regexp.MustCompile(`^\d{3,4}$`).MatchString(code) {
		return nil, errors.New("security code must be 3 or 4 digits")
	}
	return &SecurityCode{Value: code}, nil
}

type ExpirationDate struct {
	month int
	year  int
}

func NewExpirationDate(month, year string) (*ExpirationDate, error) {
	m, err := strconv.Atoi(month)
	if err != nil || m < 1 || m > 12 {
		return nil, errors.New("invalid month")
	}

	y, err := strconv.Atoi(year)
	if err != nil {
		return nil, errors.New("invalid year")
	}

	if y < 100 {
		y += 2000
	}

	now := time.Now()
	if y < now.Year() || (y == now.Year() && m < int(now.Month())) {
		return nil, errors.New("card has expired")
	}

	return &ExpirationDate{
		month: m,
		year:  y,
	}, nil
}

func (ed *ExpirationDate) String() string {
	return fmt.Sprintf("%02d/%d", ed.month, ed.year)
}

type CardHolderName struct {
	Value string
}

func NewCardHolderName(name string) (*CardHolderName, error) {
	name = strings.TrimSpace(name)
	if len(name) < 2 || len(name) > 100 {
		return nil, errors.New("invalid cardholder name length")
	}

	if !regexp.MustCompile(`^[a-zA-Z\s\-']+$`).MatchString(name) {
		return nil, errors.New("cardholder name contains invalid characters")
	}

	return &CardHolderName{Value: name}, nil
}

type Card struct {
	number       *CardNumber
	securityCode *SecurityCode
	expiration   *ExpirationDate
	holderName   *CardHolderName
}

func NewCard(pan, securityCode, expirationMonth, expirationYear, cardHolderName string) (*Card, error) {
	number, err := NewCardNumber(pan)
	if err != nil {
		return nil, fmt.Errorf("invalid card number: %w", err)
	}

	code, err := NewSecurityCode(securityCode)
	if err != nil {
		return nil, fmt.Errorf("invalid security code: %w", err)
	}

	expiration, err := NewExpirationDate(expirationMonth, expirationYear)
	if err != nil {
		return nil, fmt.Errorf("invalid expiration date: %w", err)
	}

	holder, err := NewCardHolderName(cardHolderName)
	if err != nil {
		return nil, fmt.Errorf("invalid cardholder name: %w", err)
	}

	return &Card{
		number:       number,
		securityCode: code,
		expiration:   expiration,
		holderName:   holder,
	}, nil
}

func (c *Card) Number() *CardNumber {
	return c.number
}

func (c *Card) LastFourDigits() string {
	return c.number.LastFourDigits()
}

func (c *Card) ExpirationDate() *ExpirationDate {
	return c.expiration
}

func (c *Card) CardHolderName() *CardHolderName {
	return c.holderName
}

func (c *Card) SecurityCode() *SecurityCode {
	return c.securityCode
}

func isLuhnValid(pan string) bool {
	var sum int
	nDigits := len(pan)
	parity := nDigits % 2

	for i := 0; i < nDigits; i++ {
		digit, _ := strconv.Atoi(string(pan[i]))
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
