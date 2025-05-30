package daysteps

import (
	"fmt"
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
		return 0, 0, nil
	}

	// 1 элемент - преобразовываем в тип int
	steps, err := strconv.Atoi(dataString[0])
	if err != nil || steps <= 0 {
		return 0, 0, nil
	}

	// 2 элемент - преобразовываем в тип time.Duration
	durationOfTheWalk, err := time.ParseDuration(dataString[1])
	if err != nil {
		return 0, 0, nil
	}

	return steps, durationOfTheWalk, err
}

func DayActionInfo(data string, weight, height float64) string {
	steps, durationOfTheWalk, err := parsePackage(data)
	if err != nil {
		return fmt.Sprintf("Ошибка: %v", err)
	}

	// Дистанция в метрах
	distanceInM := float64(steps) * stepLength

	// Дистанция в километрах
	distanceInKm := distanceInM / mInKm

	// Калории потраченые на прогулке
	caloriesExpended, err := spentcalories.WalkingSpentCalories(steps, weight, height, durationOfTheWalk)
	if err != nil {
		return fmt.Sprintf("Ошибка: %v", err)
	}

	activityOutput := fmt.Sprintf(`Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.`, steps, distanceInKm, caloriesExpended)
	return activityOutput
}
