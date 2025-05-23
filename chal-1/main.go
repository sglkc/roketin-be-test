package main

import (
	"fmt"
	"os"
)

/*
On earth 1 day consists of 24 hours, 1 hour = 60 minutes, 1 minute 60 seconds.
On Roketin Planet 1 day consists of 10 hours, 1 hour 100 minutes, 1 minute 100 seconds.
Make a conversion from clock on earth to clock on Roketin planet
Input in the form of 3 integers, namely hours, minutes, seconds
Output example:
on earth 00:00:00, on planet Roketin Planet .. : .. : ..
*/
func main() {
	var h, m, s int
	fmt.Print("Enter earth time (hh mm ss): ")
	_, err := fmt.Scanf("%d %d %d", &h, &m, &s)
	if err != nil {
		fmt.Println("Invalid input")
		os.Exit(1)
	}

	earthSeconds := h*3600 + m*60 + s

	// Earth: 24 * 60 * 60 = 86400 seconds/day
	// Roketin: 10 * 100 * 100 = 100000 seconds/day
	// Ratio: 100000 / 86400
	roketinSeconds := earthSeconds * (100000 / 86400)

	roketinHours := roketinSeconds / (100 * 100)
	roketinMinutes := (roketinSeconds / 100) % 100
	roketinSeconds = roketinSeconds % 100

	// Output
	fmt.Printf("On earth %02d:%02d:%02d, on planet Roketin Planet %03d:%03d:%03d\n",
		h, m, s, roketinHours, roketinMinutes, roketinSeconds)
}
