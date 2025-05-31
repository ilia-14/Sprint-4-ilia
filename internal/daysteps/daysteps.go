package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// Разбивам строку на элементы
	dataString := strings.Split(data, ",")

	// Проверяем длину слайса dateString
	if len(dataString) != 2 {
		return 0, 0, errors.New("неверный формат данных")
	}

	// 1 элемент - преобразовываем в тип int
	steps, err := strconv.Atoi(dataString[0])
	if err != nil {
		return 0, 0, errors.New("неверный формат количества шагов")
	} else if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше 0")
	}

	// 2 элемент - преобразовываем в тип time.Duration
	durationOfTheWalk, err := time.ParseDuration(dataString[1])
	if err != nil {
		return 0, 0, errors.New("неверный формат продолжительности")
	} else if durationOfTheWalk <= 0 {
		return 0, 0, errors.New("продолжительность должно быть больше 0")
	}

	return steps, durationOfTheWalk, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, durationOfTheWalk, err := parsePackage(data)
	if err != nil {
		log.Printf("Ошибка: %v", err)
		return ""
	}

	if steps <= 0 {
		log.Printf("Ошибка: количество шагов должно быть больше 0")
		return ""
	}

	if durationOfTheWalk <= 0 {
		log.Printf("Ошибка: продолжительность должно быть больше 0")
		return ""
	}

	// Дистанция в метрах
	distanceInM := float64(steps) * stepLength

	// Дистанция в километрах
	distanceInKm := distanceInM / mInKm

	// Калории потраченые на прогулке
	caloriesExpended, err := spentcalories.WalkingSpentCalories(steps, weight, height, durationOfTheWalk)
	if err != nil {
		return ""
	}

	activityOutput := fmt.Sprintf(`Количество шагов: %d.
Дистанция составила %.2f км.
Вы сожгли %.2f ккал.
`, steps, distanceInKm, caloriesExpended)
	return activityOutput
}
