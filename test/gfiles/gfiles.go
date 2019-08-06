package gfiles

import (
	"archive/zip"
	"crypto/md5"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type fc func(string) (map[string]struct{}, error)

func getAllFiles() (f fc) {
	files := make(map[string]struct{}, 0)
	f = func(root string) (map[string]struct{}, error) {
		rd, err := ioutil.ReadDir(root)
		for _, fi := range rd {
			newRoot := root + "/" + fi.Name()
			if fi.IsDir() {
				f(newRoot)
			} else {
				files[newRoot] = struct{}{}
			}
		}

		return files, err
	}
	return f
}

//GetAllFiles 获取目录下所有文件
func GetAllFiles(root string) (map[string]struct{}, error) {
	fc := getAllFiles()
	return fc(root)
}

//DeCompressZip 解压文件
func DeCompressZip(path string) ([]byte, error) {
	rc, err := zip.OpenReader(path)
	defer rc.Close()
	bs := make([]byte, 0)
	if err != nil {
		return bs, err
	}
	for _, file := range rc.File {
		f, err := file.Open()
		if err != nil {
			return bs, err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return bs, err
		}

		bs = append(bs, b...)
	}
	return bs, err
}

// srcFile could be a single file or a directory
func Zip(srcFile string, destZip string) error {
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+"/")
		// header.Name = path
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}


		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//计算文件MD5值
func GetFileMD5(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return nil, err
	}

	return md5hash.Sum(nil), nil
}

//计算文件MD5值
func GetFileMD51(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return nil, err
	}

	return md5hash.Sum(nil), nil
}
