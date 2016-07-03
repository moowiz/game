package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func readOBJ(filename string) (*Poly, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	filePath := filepath.Dir(filename)

	poly, err := parsePoly(file, filePath)
	if err != nil {
		return nil, err
	}
	return poly, err
}

func parsePoly(r io.Reader, dataDir string) (*Poly, error) {
	scanner := bufio.NewScanner(r)

	poly := &Poly{}
	verts := [][]float32{}
	uvs := [][]float32{}
	normals := [][]float32{}
	vertIndices := [][]int{}
	uvIndices := [][]int{}
	normalIndices := [][]int{}
	var matMap materialMap
	var err error

	for scanner.Scan() {
		innerScan := bufio.NewScanner(bytes.NewReader(scanner.Bytes()))
		innerScan.Split(bufio.ScanWords)
		if !innerScan.Scan() {
			if !strings.HasPrefix(innerScan.Text(), "#") {
				fmt.Printf("ignoring %s", innerScan.Text())
			}
			continue
		}
		command := innerScan.Text()

		switch command {
		case "#":
			continue
		case "v":
			new, err := parseFloats(innerScan, scanner.Text())
			if err != nil {
				return nil, err
			}
			verts = append(verts, new)
		case "vt":
			new, err := parseFloats(innerScan, scanner.Text())
			if err != nil {
				return nil, err
			}
			uvs = append(uvs, new)
		case "vn":
			new, err := parseFloats(innerScan, scanner.Text())
			if err != nil {
				return nil, err
			}
			normals = append(normals, new)
		case "f":
			tmpVIs := []int{}
			tmpUVIs := []int{}
			tmpNIs := []int{}
			for innerScan.Scan() {
				bySlash := strings.Split(innerScan.Text(), "/")
				parsed, err := strconv.ParseInt(bySlash[0], 10, 32)
				if err != nil {
					return nil, fmt.Errorf(
						"while parsing '%s': %s", scanner.Text(), err)
				}
				tmpVIs = append(tmpVIs, int(parsed))
				if len(bySlash) == 1 {
					continue
				}
				if len(bySlash[1]) > 0 {
					parsed, err = strconv.ParseInt(bySlash[1], 10, 32)
					if err != nil {
						return nil, fmt.Errorf(
							"while parsing '%s': %s", scanner.Text(), err)
					}
					tmpUVIs = append(tmpUVIs, int(parsed))
				}
				parsed, err = strconv.ParseInt(bySlash[2], 10, 32)
				if err != nil {
					return nil, fmt.Errorf(
						"while parsing '%s': %s", scanner.Text(), err)
				}
				tmpNIs = append(tmpNIs, int(parsed))
			}
			vertIndices = append(vertIndices, tmpVIs)
			uvIndices = append(uvIndices, tmpUVIs)
			normalIndices = append(normalIndices, tmpNIs)
		case "mtllib":
			if !innerScan.Scan() {
				return nil, fmt.Errorf(
					"expected material lib name while reading material")
			}
			matMap, err = readMaterialLibrary(innerScan.Text(), dataDir)
			if err != nil {
				return nil, err
			}
		case "usemtl":
			// fake shit. fix me
			if !innerScan.Scan() {
				return nil, fmt.Errorf(
					"expected material name while reading material")
			}
			name := innerScan.Text()

			poly.material = matMap[name]
		default:
			fmt.Printf("Ignoring line %s\n", scanner.Text())
		}
	}

	for i := range vertIndices {
		face := vertIndices[i]
		faceUVs := uvIndices[i]
		faceNormals := normalIndices[i]
		for ind, vertexIndex := range face {
			vert := verts[vertexIndex-1]
			poly.verts = append(poly.verts, vert...)
			if len(faceUVs) > 0 {
				poly.uvs = append(poly.uvs, uvs[faceUVs[ind]-1]...)
			}
			if len(faceNormals) > 0 {
				poly.normals = append(poly.normals, normals[faceNormals[ind]-1]...)
			}
		}
	}
	return poly, nil
}

type materialMap map[string]*material

func readMaterialLibrary(filename, dataDir string) (materialMap, error) {
	file, err := os.Open(filepath.Join(dataDir, filename))
	if err != nil {
		return nil, err
	}
	matMap := make(map[string]*material)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		innerScan := bufio.NewScanner(bytes.NewReader(scanner.Bytes()))
		innerScan.Split(bufio.ScanWords)
		if !innerScan.Scan() {
			continue
			//return nil, fmt.Errorf("bad material %q", scanner.Text())
		}

		switch innerScan.Text() {
		case "#":
			continue
		case "newmtl":
			if !innerScan.Scan() {
				return nil, fmt.Errorf(
					"expected material name while reading material")
			}
			name := innerScan.Text()
			fmt.Printf("name %s\n", name)
			mat, err := parseMaterial(scanner, dataDir)
			if err != nil {
				return nil, err
			}
			matMap[name] = mat
		}
	}
	return matMap, nil
}

func parseFloats(scanner *bufio.Scanner, line string) ([]float32, error) {
	tmp := []float32{}
	for scanner.Scan() {
		parsed, err := strconv.ParseFloat(scanner.Text(), 32)
		if err != nil {
			return nil, fmt.Errorf("while parsing '%s': %s", scanner.Text(), err)
		}
		tmp = append(tmp, float32(parsed))
	}
	return tmp, nil
}

func newTexture(file string) (uint32, error) {
	fmt.Println("reading texture", file)
	imgFile, err := os.Open(file)
	if err != nil {
		panic(fmt.Sprintf("texture %s not found on disk: %s", file, err))
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}
