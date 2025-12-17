package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var errVal string = "Not correct card number"

type Bank struct {
	Name    string
	BinFrom int
	BinTo   int
}

func loadBankData(path string) ([]Bank, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var banks []Bank

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, ",")
		if len(words) != 3 {
			return nil, fmt.Errorf("expected 3 words, got %d", len(words))
		}

		name := strings.TrimSpace(words[0])
		binFrom, err := strconv.Atoi(strings.TrimSpace(words[1]))
		if err != nil {
			return nil, err
		}
		binTo, err := strconv.Atoi(strings.TrimSpace(words[2]))
		if err != nil {
			return nil, err
		}
		banks = append(banks, Bank{name, binFrom, binTo})
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return banks, nil
}

func extractBIN(cardNumber string) int {
	num, err := strconv.Atoi(cardNumber[:6])
	if err != nil {
		return 0
	}
	return num
}

func identifyBank(bin int, banks []Bank) string {
	for _, b := range banks {
		if b.BinFrom <= bin && bin <= b.BinTo {
			return b.Name
		}
	}
	return "Неизвестный банк"
}

func validateLuhn(cardNumber string) bool {
	nums := make([]int, len(cardNumber))
	var err error
	for i := 0; i < len(cardNumber); i++ {
		nums[i], err = strconv.Atoi(string(cardNumber[i]))
		if err != nil {
			return false
		}
	}

	double := false
	for i := len(nums) - 1; i >= 0; i-- {
		if double == false {
			double = true
			continue
		}
		nums[i] *= 2
		if nums[i] > 9 {
			nums[i] = nums[i]%10 + nums[i]/10
		}
		double = false
	}

	sum := 0
	for i := 0; i < len(nums); i++ {
		sum += nums[i]
	}
	if sum%10 != 0 {
		return false
	}
	return true
}

func getUserInput() string {
	fmt.Println("Enter your card number: ")

	in := bufio.NewReader(os.Stdin)
	line, err := in.ReadString('\n')
	if err != nil {
		if err == io.EOF && strings.TrimSpace(line) == "" {
			return ""
		}
	}

	return strings.TrimSpace(line)
}

func validateInput(cardNumber string) bool {
	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		return false
	}

	for i := 0; i < len(cardNumber); i++ {
		if cardNumber[i] < '0' || cardNumber[i] > '9' {
			return false
		}
	}

	return true
}

func main() {
	fmt.Println("Добро пожаловать в программу валидации карт!")

	banks, err := loadBankData("banks.txt")
	if err != nil {
		fmt.Println("Ошибка: ", err)
	}

	for {
		cardNumber := getUserInput()
		if cardNumber == "" {
			break
		}

		if !validateInput(cardNumber) {
			fmt.Println(errVal)
			continue
		}

		if !validateLuhn(cardNumber) {
			fmt.Println(errVal)
			continue
		}

		fmt.Println("Номер карты валиден")

		six := extractBIN(cardNumber)

		bank := identifyBank(six, banks)

		if bank != "Неизвестный банк" {
			fmt.Println("Банк: {", bank, "}")
		} else {
			fmt.Println("Эмитент не определен")
		}
	}
}
