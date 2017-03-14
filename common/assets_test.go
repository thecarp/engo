package common

import(
	"testing"
	"encoding/base64"
	"strings"
	"bytes"
	"image"
	"engo.io/engo"
	_ "image/png"
)

var PngTestB64 = `iVBORw0KGgoAAAANSUhEUgAAAHgAAAAoBAMAAADHxCMWAAAABGdBTUEAALGPC/xhBQAAACBjSFJN
AAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAGFBMVEUAAAAAbQAAAFUAAKoA
AP8AkgAA/wAAAAASr1Q4AAAAAXRSTlMAQObYZgAAAAFiS0dEAIgFHUgAAAAHdElNRQfgDB4KOBZB
evDzAAABC0lEQVRIx+2UPQ7CMAyFG6ruyZKZoQfIwAlKDlCW7hnI/Y9A89c4toEiWJDq8cmfnq0X
p+uO2kp8BStQB3zAP4fVf479GXwy+2AhKduP1/MeWCxuJr4XayfzHhaL95gObEvzsNB+LddMfhqt
RTQLB99IA+8+s9bWvTl40D7XffNOM2fvTIsbhYsvnByyZXKhHYRl3belWzbRq0cDhx2hb9kbs2Hv
mAeEtZsxG+nRkprCfA27dmpPyjHOk0kuSisVmwLLO4OcalqDpn2chumUdNuZUuU0lLNhMp2fa5Cu
77OmWl8yp1W6uQxN+zit0OiiF3wDvMbcc/bBfZwWrhKz8deQezTyh6VO+VJ7AJzqxGwintLZAAAA
AElFTkSuQmCC`

// TestAssetProposal 
func TestAssetProposal(t *testing.T) {
	engo.SetHeadless(true)
	// Load test data into an image.Image
	buff, err := base64.StdEncoding.DecodeString(strings.TrimSpace(PngTestB64) )
	img, s, err := image.Decode(bytes.NewReader(buff))
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	if err != nil {
		t.Errorf("Got error in decode can't demo bad data", err)
	}
	if s != "png" {
		t.Errorf("Data not png? Did not expect:", s)
	}
	if width != 120 {
		t.Errorf("Data failed to load 120 cols:", width)
	}
	if height != 40 {
		t.Errorf("Data failed to load 40 rows:", height)
	}

	if err := engo.Files.Insert("bar.png", img); err != nil {
		t.Errorf("Data failed to insert: %v", err)
	}

	bargen, err := engo.Files.Resource("bar.png")
	bar, _ := bargen.(TextureResource)

	if err != nil {
		t.Errorf("Error loading bar.png", err)
	}

	if bar.Width != float32(width) {
		t.Errorf("Data failed to load back %v cols: %v", width, bar.Width)
	}

}
