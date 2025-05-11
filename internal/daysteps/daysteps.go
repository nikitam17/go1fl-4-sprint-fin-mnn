package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	spclr "github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// 1. разделить строку на слайс строк
	slice := strings.Split(data, ",")

	// 2. проверить, чтобы длина слайса была равна 2, так как в строке данных у нас количество шагов и продолжительность
	if len(slice) != 2 {
		return 0, 0, fmt.Errorf("длина слайса не равна 2")
	}

	// 3. преобразовать первый элемент слайса (количество шагов) в тип int. Обработать возможные ошибки. При их возникновении
	// из функции вернуть 0 шагов, 0 продолжительность и ошибку
	steps, err := strconv.Atoi(slice[0])
	if err != nil {
		return 0, 0, fmt.Errorf("не удалось преобразовать количество шагов: %v", err)
	}

	// 4. проверить: количество шагов должно быть больше 0. Если это не так, вернуть нули и ошибку
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	// 5. преобразовать второй элемент слайса в time.Duration. В пакете time есть метод для парсинга строки в time.Duration.
	// Обработать возможные ошибки. При их возникновении из функции вернуть 0 шагов, 0 продолжительность и ошибку
	duration, err := time.ParseDuration(slice[1])
	if err != nil {
		return 0, 0, fmt.Errorf("не удалось преобразовать продолжительность: %v", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	// 6. если всё прошло без ошибок, верните количество шагов, продолжительность и nil (для ошибки)
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// 1.получить данные о количестве шагов и продолжительности прогулки с помощью функции parsePackage().
	// В случае возникновения ошибки вывести её на экран и вернуть пустую строку
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// 2. проверить, чтобы количество шагов было больше 0. В противном случае вернуть пустую строку
	if steps <= 0 {
		return ""
	}

	// 3. вычислить дистанцию в метрах. Дистанция равна произведению количества шагов на длину шага.
	// Константа stepLength (длина шага) уже определена в коде
	dist := float64(steps) * stepLength

	// 4. перевести дистанцию в километры, разделив её на число метров в километре (константа mInKm, определена в пакете)
	dist /= mInKm

	// 5. вычислить количество калорий, потраченных на прогулке. Функция для вычисления калорий WalkingSpentCalories()
	// будет определена в пакете spentcalories, которую вы тоже реализуете.
	calories, err1 := spclr.WalkingSpentCalories(steps, weight, height, duration)
	if err1 != nil {
		return ""
	}

	// 6. сформировать строку вида
	// Количество шагов: 792.
	// Дистанция составила 0.51 км.
	// Вы сожгли 221.33 ккал.
	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, dist, calories)
	return result
}

// Добрый день!
// При запуске тестов для daysteps.go почему-то выдаёт ошибки вида:
//--- FAIL: TestDayStepsSuite/TestDayActionInfo/нулевая_продолжительность (0.00s)
//daysteps_test.go:342:
//			Error Trace:    C:/Users/mnn/Documents/Golang/Dev/go1fl-4-sprint-fin-mnn/internal/daysteps/daysteps_test.go:342
//													C:/Users/mnn/go/pkg/mod/github.com/stretchr/testify@v1.10.0/suite/suite.go:115
//			Error:          Should NOT be empty, but was
//			Test:           TestDayStepsSuite/TestDayActionInfo/нулевая_продолжительность
//			Messages:       Ожидался вывод в лог, но его нет

// Не понятно, что проверялось. Причем здесь это?
