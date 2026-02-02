package segmentation

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

// DiagramDetector handles the identification of non-text regions in exam papers
type DiagramDetector struct{}

func NewDiagramDetector() *DiagramDetector {
	return &DiagramDetector{}
}

// DetectRegions identifies potential diagram regions in an image
func (d *DiagramDetector) DetectRegions(imgBytes []byte) ([]image.Rectangle, error) {
	img, err := gocv.IMDecode(imgBytes, gocv.IMReadColor)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	// Blur to reduce noise
	blurred := gocv.NewMat()
	defer blurred.Close()
	gocv.GaussianBlur(gray, &blurred, image.Pt(5, 5), 0, 0, gocv.BorderDefault)

	// Edge detection
	edges := gocv.NewMat()
	defer edges.Close()
	gocv.Canny(blurred, &edges, 50, 150)

	// Dilate to connect nearby edges
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(5, 5))
	defer kernel.Close()
	dilated := gocv.NewMat()
	defer dilated.Close()
	gocv.Dilate(edges, &dilated, kernel)

	// Find contours
	contours := gocv.FindContours(dilated, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	
	var regions []image.Rectangle
	for i := 0; i < contours.Size(); i++ {
		area := gocv.ContourArea(contours.At(i))
		// Filter out small noise or very thin lines (likely text)
		if area > 5000 { // Threshold for "diagram" size
			rect := gocv.BoundingRect(contours.At(i))
			regions = append(regions, rect)
		}
	}

	return regions, nil
}

// ExtractRegion crops a specific rectangle from the image
func (d *DiagramDetector) ExtractRegion(imgBytes []byte, rect image.Rectangle) ([]byte, error) {
	img, err := gocv.IMDecode(imgBytes, gocv.IMReadColor)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	cropped := img.Region(rect)
	defer cropped.Close()

	buf, err := gocv.IMEncode(".png", cropped)
	if err != nil {
		return nil, err
	}
	return buf.GetBytes(), nil
}

// Helper to draw detected regions (useful for debugging/UI)
func (d *DiagramDetector) DrawRegions(imgBytes []byte, regions []image.Rectangle) ([]byte, error) {
	img, err := gocv.IMDecode(imgBytes, gocv.IMReadColor)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	for _, rect := range regions {
		gocv.Rectangle(&img, rect, color.RGBA{255, 0, 0, 0}, 3)
	}

	buf, err := gocv.IMEncode(".jpg", img)
	if err != nil {
		return nil, err
	}
	return buf.GetBytes(), nil
}
