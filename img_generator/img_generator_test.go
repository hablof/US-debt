package imggenerator

import (
	"os"
	"testing"
)

func TestGenerateImage(t *testing.T) {
	ig := ImageGenerator{}

	ig.GenerateImage(31462147316778000, "test_31462147316778000.png")
	ig.GenerateImage(3146214731677800, "test_3146214731677800.png")
	ig.GenerateImage(314621473167780, "test_314621473167780.png")
	ig.GenerateImage(31462147316778, "test_31462147316778.png") // original
	ig.GenerateImage(3146214731677, "test_3146214731677.png")
	ig.GenerateImage(314621473167, "test_314621473167.png")
	ig.GenerateImage(31462147316, "test_31462147316.png")
}

func TestErase(t *testing.T) {
	os.Remove("test_31462147316778000.png")
	os.Remove("test_3146214731677800.png")
	os.Remove("test_314621473167780.png")
	os.Remove("test_31462147316778.png")
	os.Remove("test_3146214731677.png")
	os.Remove("test_314621473167.png")
	os.Remove("test_31462147316.png")
}
