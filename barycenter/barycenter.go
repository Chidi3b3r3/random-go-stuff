package barycenter

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type MassPoint struct {
	x, y, z, mass float32
}

func addMassPoint(a, b MassPoint) MassPoint {
	return MassPoint{
		a.x + b.x,
		a.y + b.y,
		a.z + b.z,
		a.mass + b.mass,
	}
}

func averageMassPoint(a, b MassPoint) MassPoint {
	sum := addMassPoint(a, b)
	return MassPoint{
		sum.x / 2,
		sum.y / 2,
		sum.z / 2,
		sum.mass,
	}
}

func toWeightedSubspace(a MassPoint) MassPoint {
	return MassPoint{
		a.x * a.mass,
		a.y * a.mass,
		a.z * a.mass,
		a.mass,
	}
}

func fromWeightedSubspace(a MassPoint) MassPoint {
	return MassPoint{
		a.x / a.mass,
		a.y / a.mass,
		a.z / a.mass,
		a.mass,
	}
}

func averageMassPointWeighted(a, b MassPoint) MassPoint {
	aWeightedSubspace := toWeightedSubspace(a)
	bWeightedSubspace := toWeightedSubspace(b)

	return fromWeightedSubspace(averageMassPoint(aWeightedSubspace, bWeightedSubspace))
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func closeFile(fi *os.File) {
	err := fi.Close()

	handleError(err)
}

func Compute() {
	if len(os.Args) < 2 {
		fmt.Println("incomplete args")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	handleError(err)
	defer closeFile(file)

	var masspoints []MassPoint

	startLoading := time.Now()

	for {
		var newMasspoint MassPoint
		_, err := fmt.Fscanf(file, "%f:%f:%f:%f", &newMasspoint.x, &newMasspoint.y, &newMasspoint.z, &newMasspoint.mass)
		if err == io.EOF {
			break
		} else if err != nil {
			continue
		}
		masspoints = append(masspoints, newMasspoint)
	}

	fmt.Printf("Loaded %d values from file in %s", len(masspoints), time.Since(startLoading))

	if len(masspoints) <= 1 {
		handleError(errors.New("insufficient masspoints"))
	}

	startCalculation := time.Now()

	for len(masspoints) > 1 {
		var newMasspoints []MassPoint

		for i := 0; i < len(masspoints)-1; i += 2 {
			newMasspoints = append(newMasspoints, averageMassPointWeighted(masspoints[i], masspoints[i+1]))
		}

		if len(masspoints)%2 != 0 {
			newMasspoints = append(newMasspoints, masspoints[len(masspoints)-1])
		}

		masspoints = newMasspoints
	}

	systemAvergage := masspoints[0]

	fmt.Printf("System barycenter is at (%f, %f, %f) and the system's mass is %f.\n",
		systemAvergage.x,
		systemAvergage.y,
		systemAvergage.z,
		systemAvergage.mass,
	)

	fmt.Printf("Calculation took %s.\n", time.Since(startCalculation))
}

func DataGeneration() {
	if len(os.Args) < 2 {
		fmt.Println("missing arguments")
		os.Exit(1)
	}

	nBodies, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	rand.Seed(time.Now().Unix())
	maxPos := 100
	maxMass := 5

	for i := 0; i <= nBodies; i++ {
		posX := rand.Intn(maxPos*2) - maxPos
		posY := rand.Intn(maxPos*2) - maxPos
		posZ := rand.Intn(maxPos*2) - maxPos
		mass := rand.Intn(maxMass-1) + 1

		fmt.Printf("%d:%d:%d:%d\n", posX, posY, posZ, mass)
	}
}

func NaiveImplementation() string {
	return "naive"
}
