package object

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

	verts := [][]float32{}
	uvs := [][]float32{}
	normals := [][]float32{}
	vertIndices := [][]int{}
	uvIndices := [][]int{}
	normalIndices := [][]int{}
	var pMat *material
	var matMap materialMap
	var err error

	for scanner.Scan() {
		innerScan := bufio.NewScanner(bytes.NewReader(scanner.Bytes()))
		innerScan.Split(bufio.ScanWords)
		if !innerScan.Scan() {
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
			num := 0
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
				if num > 3 {
					fmt.Println("more than 3")
				}
				num++
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
			if !innerScan.Scan() {
				return nil, fmt.Errorf(
					"expected material name while reading material")
			}
			name := innerScan.Text()

			pMat = matMap[name]
		default:
			//fmt.Printf("Ignoring line %s\n", scanner.Text())
		}
	}

	pVerts := []float32{}
	pUvs := []float32{}
	pNormals := []float32{}
	for i := range vertIndices {
		face := vertIndices[i]
		faceUVs := uvIndices[i]
		faceNormals := normalIndices[i]
		for ind := range face {
			vert := verts[face[ind]-1]
			pVerts = append(pVerts, vert...)
			if len(faceUVs) > 0 {
				uv := uvs[faceUVs[ind]-1]
				//uv[1] = 1 - uv[1]
				pUvs = append(pUvs, uv...)
			}
			if len(faceNormals) > 0 {
				pNormals = append(pNormals, normals[faceNormals[ind]-1]...)
			}
		}
	}
	return NewPoly(pMat, pVerts, pUvs, pNormals)
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
