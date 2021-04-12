package common

const (
	inchToMMRatio = 25.4
)

// ScanArea represents the coordinates of a scan area.
type ScanArea struct {
	TLX interface{}
	TLY interface{}
	BRX interface{}
	BRY interface{}
}

// PixelsToMillimeters returns a new ScanArea containing the value of the current instance
// converted into millimeters, using the given resolution in DPI (dots per inch).
func (sa *ScanArea) PixelsToMillimeters(resolution int) *ScanArea {
	return &ScanArea{
		TLX: pxToMM(sa.TLX.(int), resolution),
		TLY: pxToMM(sa.TLY.(int), resolution),
		BRX: pxToMM(sa.BRX.(int), resolution),
		BRY: pxToMM(sa.BRY.(int), resolution),
	}
}

// pxToMM calculates a value in millimeters from the given value in pixels and resolution
// in DPI (dots per inch).
func pxToMM(pxValue int, resolution int) float64 {
	return float64(pxValue) * inchToMMRatio / float64(resolution)
}
