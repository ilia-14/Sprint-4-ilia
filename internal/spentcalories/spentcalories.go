package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Разбивам строку на элементы
	dataString := strings.Split(data, ",")

	// Проверяем длину слайса dateString
	if len(dataString) != 3 {
		return 0, "", 0, errors.New("неверный формат данных")
	}

	// 1 элемент - преобразовываем в тип int
	steps, err := strconv.Atoi(dataString[0])
	if err != nil {
		return 0, "", 0, errors.New("неверный формат количества шагов")
	} else if steps <= 0 {
		return 0, "", 0, errors.New("количнство шагов должно быть больше 0")
	}

	// 2 элемент - вид активности
	tupeOfActivity := strings.TrimSpace(dataString[1])
	if tupeOfActivity != "Ходьба" && tupeOfActivity != "Бег" {
		return 0, "", 0, errors.New("неверный тип тренировки")
	}

	// 3 элемент - преобразовываем в тип time.Duration
	durationOfTheTraning, err := time.ParseDuration(dataString[2])
	if err != nil {
		return 0, "", 0, errors.New("неверный формат продолжительности")
	} else if durationOfTheTraning <= 0 {
		return 0, "", 0, errors.New("продолжительность должно быть больше 0")
	}

	return steps, tupeOfActivity, durationOfTheTraning, nil
}

func distance(steps int, height float64) float64 {
	// Расчет длины шага
	estimatedStepLength := height * stepLengthCoefficient

	// Дистанция в метрах
	estimatedDistanceInM := float64(steps) * estimatedStepLength

	// Дистанция в километрах
	estimatedDistanceInKm := estimatedDistanceInM / mInKm

	return estimatedDistanceInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверка на корректность входящих данных
	if duration <= 0 {
		return 0
	}

	// Дистанция в километрах
	estimatedDistanceInKm := distance(steps, height)

	// Средняя скорость
	meanSpeedTraining := estimatedDistanceInKm / duration.Hours()

	return meanSpeedTraining
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка на корректность входящих данных
	if steps <= 0 {
		return 0, errors.New("количнство шагов должно быть больше 0")
	}

	if weight <= 0 {
		return 0, errors.New("вес должно быть больше 0")
	}

	if duration <= 0 {
		return 0, errors.New("продолжительность должно быть больше 0")
	}

	// Расчет средней скорости при беге
	meanSpeedRunning := meanSpeed(steps, height, duration)

	// Переводим продолжительность бега в минуты
	durationOfTheRunningInMinuters := duration.Minutes()

	runningCalories := (weight * meanSpeedRunning * durationOfTheRunningInMinuters) / minInH

	return runningCalories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка на корректность входящих данных
	if steps <= 0 {
		return 0, errors.New("количнство шагов должно быть больше 0")
	}

	if weight <= 0 {
		return 0, errors.New("вес должно быть больше 0")
	}

	if height <= 0 {
		return 0, errors.New("рост должно быть больше 0")
	}

	if duration <= 0 {
		return 0, errors.New("продолжительность должно быть больше 0")
	}

	// Расчет средней скорости при ходьбе
	meanSpeedWalking := meanSpeed(steps, height, duration)

	// Переводим продолжительность ходьбы в минуты
	durationOfTheWalkingInMinuters := duration.Minutes()

	//
	walkingCalories := (weight * meanSpeedWalking * durationOfTheWalkingInMinuters * walkingCaloriesCoefficient) / minInH

	return walkingCalories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, tupeOfActivity, durationOfTheTraning, err := parseTraining(data)
	if tupeOfActivity != "Ходьба" && tupeOfActivity != "Бег" {
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", errors.New("неверный тип данных")
	}

	durationOfThe := distance(steps, height)
	meanSpeed := meanSpeed(steps, height, durationOfTheTraning)
	var caloriesExpended float64

	switch tupeOfActivity {
	case "Ходьба":
		caloriesExpended, err = WalkingSpentCalories(steps, weight, height, durationOfTheTraning)
		if err != nil {
			return "", errors.New("неверный тип данных")
		}
	case "Бег":
		caloriesExpended, err = RunningSpentCalories(steps, weight, height, durationOfTheTraning)
		if err != nil {
			return "", errors.New("неверный тип данных")
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	bottomLine := fmt.Sprintf(`Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f
`, tupeOfActivity, durationOfTheTraning.Hours(), durationOfThe, meanSpeed, caloriesExpended)

	return bottomLine, nil
}
