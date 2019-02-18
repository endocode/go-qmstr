package service

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func NewFileNode(path string, fileType FileNode_Type) *FileNode {
	filename := filepath.Base(path)
	hash, err := HashFile(path)
	broken := false
	if err != nil {
		hash = "nohash" + path
		broken = true
	}
	return &FileNode{Name: filename, FileType: fileType, Path: path, Hash: hash, Broken: broken}
}

func CreateInfoNode(infoType string, dataNodes ...*InfoNode_DataNode) *InfoNode {
	return &InfoNode{
		Type:      infoType,
		DataNodes: dataNodes,
	}
}

func CreateWarningNode(warning string) *InfoNode {
	return CreateInfoNode("warning", &InfoNode_DataNode{
		Type: "warning_message",
		Data: warning,
	})
}

func CreateErrorNode(errorMes string) *InfoNode {
	return CreateInfoNode("error", &InfoNode_DataNode{
		Type: "error_message",
		Data: errorMes,
	})
}

func (in *InfoNode) Describe(indent string) string {
	describe := []string{fmt.Sprintf("%s|- Type: %s, Confidence score: %v", indent, in.Type, in.ConfidenceScore)}
	indent = indent + "\t"
	for _, anode := range in.Analyzer {
		describe = append(describe, fmt.Sprintf("%s|- Name: %s", indent, anode.Name))
	}
	for _, dnode := range in.DataNodes {
		describe = append(describe, fmt.Sprintf("%s|- Type: %s, Data: %s", indent, dnode.Type, dnode.Data))
	}
	return strings.Join(describe, "\n")
}

func (fn *FileNode) Describe(less bool, indent string) string {
	describe := []string{fmt.Sprintf("%s|- Name: %s, Path: %s, Hash: %s", indent, fn.Name, fn.Path, fn.Hash)}
	indent = indent + "\t"
	if !less {
		for _, inode := range fn.AdditionalInfo {
			describe = append(describe, inode.Describe(indent))
		}
	}
	for _, fnode := range fn.DerivedFrom {
		describe = append(describe, fnode.Describe(less, indent))
	}
	return strings.Join(describe, "\n")
}

func (pn *PackageNode) Describe(less bool) string {
	describe := []string{fmt.Sprintf("|- Name: %s", pn.Name)}
	if !less {
		for _, inode := range pn.AdditionalInfo {
			describe = append(describe, inode.Describe("\t"))
		}
	}
	for _, fnode := range pn.Targets {
		describe = append(describe, fnode.Describe(less, "\t"))
	}
	return strings.Join(describe, "\n")
}

func HashFile(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	return Hash(f)
}

func Hash(r io.Reader) (string, error) {
	h := sha1.New()
	buf := make([]byte, 0, 4*1024)
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			return "", err
		}
		h.Write(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
