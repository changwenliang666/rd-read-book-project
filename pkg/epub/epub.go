package epub

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"io"
	"path/filepath"
	"strings"
)

type Metadata struct {
	Title       string
	Author      string
	Description string
	CoverImg    []byte
	RemoteUrl   string
}

// 从 *zip.Reader 中解析 EPUB 元数据
func ParseEPUBFromZipReader(r *zip.Reader) (*Metadata, error) {
	var metadata Metadata

	// =========================
	// 1. 找 META-INF/container.xml
	// =========================
	var opfPath string

	for _, f := range r.File {
		if f.Name == "META-INF/container.xml" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			type Rootfile struct {
				FullPath string `xml:"full-path,attr"`
			}

			type Container struct {
				Rootfiles []Rootfile `xml:"rootfiles>rootfile"`
			}

			var container Container
			if err := xml.NewDecoder(rc).Decode(&container); err != nil {
				return nil, err
			}

			if len(container.Rootfiles) == 0 {
				return nil, errors.New("container.xml 中没有 rootfile")
			}

			opfPath = container.Rootfiles[0].FullPath
			break
		}
	}

	if opfPath == "" {
		return nil, errors.New("找不到 content.opf")
	}

	opfDir := filepath.Dir(opfPath)

	// =========================
	// 2. 解析 content.opf
	// =========================
	var coverID string
	var coverHref string

	for _, f := range r.File {
		if f.Name != opfPath {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		// ---- XML 结构定义 ----
		type Meta struct {
			Name     string `xml:"name,attr"`
			Content  string `xml:"content,attr"`
			Property string `xml:"property,attr"`
			Value    string `xml:",chardata"`
		}

		type Item struct {
			ID        string `xml:"id,attr"`
			Href      string `xml:"href,attr"`
			MediaType string `xml:"media-type,attr"`
		}

		type Manifest struct {
			Items []Item `xml:"item"`
		}

		type MetadataXML struct {
			Title       string `xml:"http://purl.org/dc/elements/1.1/ title"`
			Creator     string `xml:"http://purl.org/dc/elements/1.1/ creator"`
			Description string `xml:"http://purl.org/dc/elements/1.1/ description"`
			Metas       []Meta `xml:"meta"`
		}

		type Package struct {
			Metadata MetadataXML `xml:"metadata"`
			Manifest Manifest    `xml:"manifest"`
		}

		var pkg Package
		if err := xml.NewDecoder(rc).Decode(&pkg); err != nil {
			return nil, err
		}

		// ---- 基本信息 ----
		metadata.Title = strings.TrimSpace(pkg.Metadata.Title)
		metadata.Author = strings.TrimSpace(pkg.Metadata.Creator)

		// ---- 简介（多兜底）----
		metadata.Description = strings.TrimSpace(pkg.Metadata.Description)

		if metadata.Description == "" {
			for _, m := range pkg.Metadata.Metas {
				if strings.ToLower(m.Name) == "description" {
					metadata.Description = strings.TrimSpace(m.Content)
					break
				}
				if strings.Contains(strings.ToLower(m.Property), "description") {
					metadata.Description = strings.TrimSpace(m.Value)
					break
				}
			}
		}

		// ---- 找封面 ID ----
		for _, m := range pkg.Metadata.Metas {
			if strings.ToLower(m.Name) == "cover" {
				coverID = m.Content
				break
			}
		}

		// ---- 找封面 href ----
		if coverID != "" {
			for _, item := range pkg.Manifest.Items {
				if item.ID == coverID {
					coverHref = item.Href
					break
				}
			}
		}

		break
	}

	// =========================
	// 3. 读取封面文件
	// =========================
	if coverHref != "" {
		coverPath := filepath.Clean(filepath.Join(opfDir, coverHref))

		for _, f := range r.File {
			if f.Name == coverPath {
				rc, err := f.Open()
				if err != nil {
					break
				}
				defer rc.Close()

				data, err := io.ReadAll(rc)
				if err == nil {
					metadata.CoverImg = data
				}
				break
			}
		}
	}

	return &metadata, nil
}
