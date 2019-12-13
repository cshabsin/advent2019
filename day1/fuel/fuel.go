package fuel

func Fuel(mass int) int {
	return mass/3-2
}

func AllFuel(mass int) int {
	base := Fuel(mass)
	if base <= 0 {
		return 0
	}
	additional := AllFuel(base)
	return base + additional
}
