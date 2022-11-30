package utils

func CalcPerc(diff int64, relativeTime int64) float64 {
	/*	svOld := "05:00:00"
		svNew := "15:00:00"

		timeOld, _ := time.Parse(config.TimeFormat, svOld)
		timeNew, _ := time.Parse(config.TimeFormat, svNew)
		diff := timeNew.Sub(timeOld)
	*/

	return float64(diff) / float64(relativeTime) * 100

}
