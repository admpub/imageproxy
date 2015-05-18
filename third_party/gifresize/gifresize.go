// Copyright 2013 Daniel Pupius. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package gifresize resizes animated gifs.
//
// Frames in an animated gif aren't necessarily the same size, subsequent
// frames are overlayed on previous frames. Therefore, resizing the frames
// individually may cause problems due to aliasing of transparent pixels. This
// package tries to avoid this by building frames from all previous frames and
// resizing the frames as RGB.
package gifresize

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io"
)

// TransformFunc is a function that transforms an image.
type TransformFunc func(image.Image) image.Image

// Process the GIF read from r, applying transform to each frame, and writing
// the result to w.
func Process(w io.Writer, r io.Reader, transform TransformFunc) error {
	// Decode the original gif.
	im, err := gif.DecodeAll(r)
	if err != nil {
		return err
	}

	// Create a new RGBA image to hold the incremental frames.
	firstFrame := im.Image[0].Bounds()
	b := image.Rect(0, 0, firstFrame.Dx(), firstFrame.Dy())
	img := image.NewRGBA(b)

	// Resize each frame.
	for index, frame := range im.Image {
		bounds := frame.Bounds()
		draw.Draw(img, bounds, frame, bounds.Min, draw.Over)
		im.Image[index] = imageToPaletted(transform(img))
	}

	return gif.EncodeAll(w, im)
}

func imageToPaletted(img image.Image) *image.Paletted {
	b := img.Bounds()
	pm := image.NewPaletted(b, palette.Plan9)
	draw.FloydSteinberg.Draw(pm, b, img, image.ZP)
	return pm
}