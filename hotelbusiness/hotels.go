//go:build !solution

package hotelbusiness

import "math"

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func getMinAndMaxDates(m map[int]int, guests []Guest) (map[int]int, int, int) {
	var maxDate = -1
	var minDate = math.MaxInt32

	for _, guest := range guests {
		m[guest.CheckInDate]++
		m[guest.CheckOutDate]--

		minDate = int(math.Min(float64(minDate), float64(guest.CheckInDate)))
		minDate = int(math.Min(float64(minDate), float64(guest.CheckOutDate)))

		maxDate = int(math.Max(float64(maxDate), float64(guest.CheckInDate)))
		maxDate = int(math.Max(float64(maxDate), float64(guest.CheckOutDate)))
	}

	return m, minDate, maxDate
}

func ComputeLoad(guests []Guest) []Load {
	var m = make(map[int]int)
	var nguests []int
	count := 0

	m, minDate, maxDate := getMinAndMaxDates(m, guests)

	for i := minDate; i <= maxDate; i++ {
		count += m[i]
		nguests = append(nguests, count)
	}

	var result = make([]Load, 0)

	for i, amount := range nguests {
		if i > 0 {
			if amount == result[len(result)-1].GuestCount {
				continue
			}
		}

		load := Load{i + minDate, amount}
		result = append(result, load)
	}

	return result
}
