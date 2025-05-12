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
	// 1. разделить строку на слайс строк
	slice := strings.Split(data, ",")

	// 2. проверить, чтобы длина слайса была равна 3, так как в строке данных у нас количество шагов, вид активности и
	// продолжительность
	if len(slice) != 3 {
		return 0, "", 0, fmt.Errorf("длина слайса не равна 3")
	}

	// 3. преобразовать первый элемент слайса (количество шагов) в тип int. Обработать возможные ошибки. При их возникновении из
	// функции вернуть 0 шагов, 0 продолжительность и ошибку
	steps, err := strconv.Atoi(slice[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать количество шагов: %v", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	// 4. преобразовать третий элемент слайса в time.Duration. В пакете time есть метод для парсинга строки в time.Duration.
	// Обработать возможные ошибки. При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку
	duration, err := time.ParseDuration(slice[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать продолжительность: %v", err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	// 5. если всё прошло без ошибок, верните количество шагов, вид активности, продолжительность и nil (для ошибки)
	return steps, slice[1], duration, nil
}

func distance(steps int, height float64) float64 {
	// рассчитайте длину шага. Для этого умножьте высоту пользователя на коэффициент длины шага stepLengthCoefficient.
	// cоответствующая константа уже определена в пакете
	lenStep := height * stepLengthCoefficient

	// умножьте пройденное количество шагов на длину шага
	dist := float64(steps) * lenStep

	// разделите полученное значение на число метров в километре (mInKm, константа определена в пакете)
	return dist / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// 1. проверить, что продолжительность duration больше 0. Если это не так, вернуть 0
	if duration <= 0 {
		return 0
	}

	// 2. вычислить дистанцию с помощью distance()
	dist := distance(steps, height)

	// 3. вычислить и вернуть среднюю скорость. Для этого разделите дистанцию на продолжительность в часах.
	// Чтобы перевести продолжительность в часы, воспользуйтесь функцией из пакета time
	return dist / duration.Hours()
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// 1. проверить входные параметры на корректность. Если параметры некорректны, вернуть 0 калорий и соответствующую ошибку
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	// 2. рассчитать среднюю скорость с помощью meanSpeed()
	speed := meanSpeed(steps, height, duration)

	// 3. рассчитать и вернуть количество калорий. Для этого:
	// 3.1. переведите продолжительность в минуты с помощью функции из пакета time
	minutes := duration.Minutes()
	// 3.2. умножьте вес пользователя на среднюю скорость и продолжительность в минутах
	calories := weight * speed * float64(minutes)
	// 3.3. разделите результат на число минут в часе для получения количества потраченных калорий
	return calories / 60, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// 1. проверить входные параметры на корректность. Если параметры некорректны, вернуть 0 калорий и соответствующую ошибку
	// 2. рассчитать среднюю скорость с помощью meanSpeed()
	// 3. рассчитать и вернуть количество калорий. Для этого:
	// 3.1. переведите продолжительность в минуты с помощью функции из пакета time
	// 3.2. умножьте вес пользователя на среднюю скорость и продолжительность в минутах
	// 3.3. разделите результат на число минут в часе для получения количества потраченных калорий
	calories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}

	// 4. умножить полученное число калорий на корректирующий коэффициент walkingCaloriesCoefficient.
	// Соответствующая константа объявляена в пакете. Вернуть полученное значение
	return calories * walkingCaloriesCoefficient, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// 1. получить значения из строки данных с помощью функции parseTraining(), обработать возможные ошибки
	steps, vid, duration, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	var calories float64

	// 2. проверить, какой вид тренировки был передан в строке, которую парсили (лучше использовать switch).
	// Для каждого из видов тренировки рассчитать дистанцию, среднюю скорость и калории
	// 3. для каждого вида тренировки сформировать и вернуть строку вида:
	// Тип тренировки: Бег
	// Длительность: 0.75 ч.
	// Дистанция: 10.00 км.
	// Скорость: 13.34 км/ч
	// Сожгли калорий: 18621.75
	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	switch vid {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	// 4. если был передан неизвестный тип тренировки, вернуть ошибку с текстом неизвестный тип тренировки
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}
	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\n", vid, duration.Hours(), dist)
	result += fmt.Sprintf("Скорость: %.2f км/ч\nСожгли калорий: %.2f\n", speed, calories)
	return result, nil
}
