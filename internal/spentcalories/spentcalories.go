package spentcalories

import (
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
		return 0, "", 0, nil
	}

	// 1 элемент - преобразовываем в тип int
	steps, err := strconv.Atoi(dataString[0])
	if err != nil || steps <= 0 {
		return 0, "", 0, nil
	}

	// 2 элемент - вид активности
	tupeOfActivity := dataString[1]

	// 3 элемент - преобразовываем в тип time.Duration
	durationOfTheTraning, err := time.ParseDuration(dataString[2])
	if err != nil {
		return 0, "", 0, nil
	}

	return steps, tupeOfActivity, durationOfTheTraning, err
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
	if steps <= 0 || duration <= 0 {
		return 0, nil
	}

	// Расчет средней скорости при беге
	meanSpeedRunning := meanSpeed(steps, height, duration)

	// Переводим продолжительность бега в минуты
	durationOfTheRunningInMinuters := duration.Minutes()

	//
	runningCalories := (weight * meanSpeedRunning * durationOfTheRunningInMinuters) / minInH

	return runningCalories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка на корректность входящих данных
	if steps <= 0 || duration <= 0 {
		return 0, nil
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
	//
	steps, tupeOfActivity, durationOfTheTraning, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	durationOfThe := distance(steps, height)
	meanSpeed := meanSpeed(steps, height, durationOfTheTraning)
	var caloriesExpended float64

	switch tupeOfActivity {
	case "Ходьба":
		caloriesExpended, err = WalkingSpentCalories(steps, weight, height, durationOfTheTraning)
	case "Бег":
		caloriesExpended, err = RunningSpentCalories(steps, weight, height, durationOfTheTraning)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	bottomLine := fmt.Sprintf(` Тип тренировки: %s
	                            Длительность: %.2f ч.
								Дистанция: %.2f км.
								Скорость: %.2f км/ч
								Сожгли калорий: %.2f`, tupeOfActivity, durationOfTheTraning.Hours(), durationOfThe, meanSpeed, caloriesExpended)

	return bottomLine, nil
}
