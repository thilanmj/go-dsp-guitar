package resample

import (
	"math"
)

/*
 * The Lanczos kernel function L(x, a).
 */
func lanczosKernel(x float64, a float64) float64 {

	/*
	 * Calculate the sections of the Lanczos kernel.
	 */
	if x == 0 {
		return 1.0
	} else if (-a < x) && (x < a) {
		piX := math.Pi * x
		piXa := piX / a
		piXsquared := piX * piX
		xSin := math.Sin(piX)
		xaSin := math.Sin(piXa)
		prodSins := xSin * xaSin
		arg := a * prodSins
		result := arg / piXsquared
		return result
	} else {
		return 0.0
	}

}

/*
 * The Lanczos interpolation function S(s, x, a).
 */
func lanczosInterpolate(s []float64, x float64, a uint16) float64 {
	floorX := math.Floor(x)
	idx := int(floorX)
	idxInc := idx + 1
	aInt := int(a)
	lBound := idxInc - aInt
	uBound := idxInc + aInt
	n := len(s)
	aFloat := float64(a)
	sum := float64(0.0)

	/*
	 * Calculate the Lanczos sum.
	 */
	for i := lBound; i < uBound; i++ {

		/*
		 * Check if we are still within the bounds of the slice.
		 */
		if (i >= 0) && (i < n) {
			iFloat := float64(i)
			diff := x - iFloat
			fac := s[i]
			val := lanczosKernel(diff, aFloat)
			sum += fac * val
		}

	}

	return sum
}

/*
 * Resample time series data from a source to a target sampling rate using the
 * Lanczos resampling method.
 */
func Time(samples []float64, sourceRate uint32, targetRate uint32) []float64 {
	inputLength := len(samples)
	inputLengthFloat := float64(inputLength)
	sourceRateFloat := float64(sourceRate)
	targetRateFloat := float64(targetRate)
	expansion := targetRateFloat / sourceRateFloat
	outputLengthFloat := inputLengthFloat * expansion
	outputLengthFloor := math.Floor(outputLengthFloat)
	outputLength := int(outputLengthFloor)

	/*
	 * If we exactly hit the last sample, do not expand the sequence.
	 */
	if outputLengthFloor == outputLengthFloat {
		outputLength--
	}

	outputBuffer := make([]float64, outputLength)
	dx := sourceRateFloat / targetRateFloat

	/*
	 * Calculate output samples using Lanczos interpolation.
	 */
	for i, _ := range outputBuffer {
		iFloat := float64(i)
		x := iFloat * dx
		val := lanczosInterpolate(samples, x, 3)
		outputBuffer[i] = val
	}

	return outputBuffer
}

/*
 * Resample frequency domain data to a different number of target bins using
 * the Lanczos resampling method.
 */
func Frequency(bins []complex128, numTargetBins uint32) []complex128 {
	numSourceBins := len(bins)
	sourceReal := make([]float64, numSourceBins)
	sourceImag := make([]float64, numSourceBins)

	/*
	 * Extract real and imaginary sequences from complex sequence.
	 */
	for i, elem := range bins {
		elemReal := real(elem)
		sourceReal[i] = elemReal
		elemImag := imag(elem)
		sourceImag[i] = elemImag
	}

	targetBins := make([]complex128, numTargetBins)
	numSourceBinsFloat := float64(numSourceBins)
	numTargetBinsFloat := float64(numTargetBins)
	dx := numSourceBinsFloat / numTargetBinsFloat

	/*
	 * Calculate output samples using Lanczos interpolation.
	 */
	for i, _ := range targetBins {
		iFloat := float64(i)
		x := iFloat * dx
		targetReal := lanczosInterpolate(sourceReal, x, 3)
		targetImag := lanczosInterpolate(sourceImag, x, 3)
		targetComplex := complex(targetReal, targetImag)
		targetBins[i] = targetComplex
	}

	return targetBins
}

/*
 * Oversample the source signal by a factor and write the result into the
 * target buffer.
 */
func Oversample(source []float64, target []float64, factor uint32) {
	factorFloat := float64(factor)
	dx := 1.0 / factorFloat
	factorInt := int(factor)

	/*
	 * Calculate output samples using Lanczos interpolation.
	 */
	for i, _ := range target {
		m := i % factorInt

		/*
		 * If we hit an input sample exactly, we can save some
		 * computation.
		 */
		if m == 0 {
			idx := i / factorInt
			val := source[idx]
			target[i] = val
		} else {
			iFloat := float64(i)
			x := iFloat * dx
			val := lanczosInterpolate(source, x, 3)
			target[i] = val
		}

	}

}
