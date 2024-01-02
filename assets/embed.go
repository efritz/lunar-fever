package assets

import (
	"embed"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
)

//go:embed **/*
var assets embed.FS

func LoadShader(name string) (string, error) {
	return decodeAsset("shaders", name, "glsl", func(r io.Reader) (string, error) {
		bytes, err := io.ReadAll(r)
		if err != nil {
			return "", err
		}

		return string(bytes), nil
	})
}

func LoadTexture(name string) (*image.RGBA, error) {
	return decodeAsset("textures", name, "png", func(r io.Reader) (*image.RGBA, error) {
		img, err := png.Decode(r)
		if err != nil {
			return nil, err
		}

		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
		return rgba, nil
	})
}

func decodeAsset[E any](assetType, name, assetExt string, reader func(r io.Reader) (E, error)) (val E, _ error) {
	file, err := assets.Open(fmt.Sprintf("%s/%s.%s", assetType, name, assetExt))
	if err != nil {
		return val, err
	}
	defer file.Close()

	return reader(file)
}
