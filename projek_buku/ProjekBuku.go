package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	cls "github.com/MasterDimmy/go-cls"
	"github.com/go-pdf/fpdf"
)

type Buku struct {
	Kode, Judul, Pengarang, Penerbit string
	Halaman, Tahun                   int
	Waktu                            time.Time
}

var ListBuku []Buku

// var KodeBukuTerpakai = make(map[string]bool)

func TambahBuku() {
	KodeB := ""
	JudulB := ""
	PengarangB := ""
	PenerbitB := ""
	JumlahHalaman := 0
	TahunTerbit := 0

	fmt.Println("")
	fmt.Println("Tambahkan Buku")
	fmt.Println("")
	draftBuku := []Buku{}

	for {

		fmt.Print("Masukan Kode Buku : ")
		_, err := fmt.Scanln(&KodeB)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}

		for _, book := range draftBuku {
			if book.Kode == "book-"+KodeB {
				fmt.Println("Kode Buku Sudah Terpakai. Silahkan Masukkan Kode Lain.")
				return
			}
		}

		listJsonBuku, err := os.ReadDir("books")
		if err != nil {
			fmt.Println("Terjadi Error: ", err)
			return
		}

		for _, bookJson := range listJsonBuku {
			if bookJson.Name() == "book-"+KodeB+".json" {
				fmt.Println("Kode Buku Sudah Terpakai. Silahkan Masukkan Kode Lain.")
				return
			}
		}

		fmt.Print("Masukan Judul Buku : ")
		_, err = fmt.Scanln(&JudulB)
		if err != nil {
			fmt.Println("Terjadi Error: ", err)
		}

		fmt.Print("Masukan Nama Pengarang : ")
		_, err = fmt.Scanln(&PengarangB)
		if err != nil {
			fmt.Println("Terjadi Error: ", err)
		}

		fmt.Print("Masukan Nama Penerbit : ")
		_, err = fmt.Scanln(&PenerbitB)
		if err != nil {
			fmt.Println("Terjadi Error: ", err)
		}

		fmt.Print("Masukan Jumlah Halaman : ")
		_, err = fmt.Scanln(&JumlahHalaman)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}

		fmt.Print("Masukan Tahun Terbit : ")
		_, err = fmt.Scanln(&TahunTerbit)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}

		draftBuku = append(draftBuku, Buku{
			Kode:      fmt.Sprintf("book-%s", KodeB),
			Judul:     JudulB,
			Pengarang: PengarangB,
			Penerbit:  PenerbitB,
			Halaman:   JumlahHalaman,
			Tahun:     TahunTerbit,
			Waktu:     time.Now(),
		})

		pilihanTambahBuku := 0
		fmt.Println("Ketik 1 untuk tambah buku lain, ketik 0 untuk selesai")
		_, err = fmt.Scanln(&pilihanTambahBuku)
		if err != nil {
			fmt.Println("Terjadi Error: ", err)
			return
		}

		if pilihanTambahBuku == 0 {
			break
		}
	}

	fmt.Println("Menambah Buku...")

	_ = os.Mkdir("books", 0777)

	ch := make(chan Buku)

	wg := sync.WaitGroup{}

	jumlahPustakawan := 5

	for i := 0; i < jumlahPustakawan; i++ {
		wg.Add(1)
		go SimpanBuku(ch, &wg, i)
	}

	for _, bukuTersimpan := range draftBuku {
		ch <- bukuTersimpan
	}

	close(ch)

	wg.Wait()

	fmt.Println("Berhasil Menambahkan Buku")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func SimpanBuku(ch <-chan Buku, wg *sync.WaitGroup, noPustakawan int) {

	for bukuTersimpan := range ch {
		dataJson, err := json.Marshal(bukuTersimpan)
		if err != nil {
			fmt.Println("Terjadi Error: ", err)
			return
		}

		err = os.WriteFile(fmt.Sprintf("books/%s.json", bukuTersimpan.Kode), dataJson, 0644)
		if err != nil {
			fmt.Println("Terjadi Error: ", err)
			return
		}
		fmt.Printf("Pekerja No %d Memproses Kode Buku : %s!\n", noPustakawan, bukuTersimpan.Kode)
	}
	wg.Done()
}

func LihatDaftarBuku(ch <-chan string, chBuku chan Buku, wg *sync.WaitGroup) {
	var buku Buku
	for kodeBuku := range ch {
		dataJson, err := os.ReadFile(fmt.Sprintf("books/%s", kodeBuku))
		if err != nil {
			fmt.Println("Terjadi Error: ", err)
		}

		err = json.Unmarshal(dataJson, &buku)
		if err != nil {
			fmt.Println("Terjadi Error: ", err)
		}

		chBuku <- buku
	}
	wg.Done()
}

func LihatList() {
	fmt.Println("")
	fmt.Println("Lihat List Buku")
	fmt.Println("")
	ListBuku = []Buku{}

	listJsonBuku, err := os.ReadDir("books")
	if err != nil {
		fmt.Println("Terjadi Error: ", err)
	}

	wg := sync.WaitGroup{}

	ch := make(chan string)
	chBuku := make(chan Buku, len(listJsonBuku))

	jumlahPustakawan := 5

	for i := 0; i < jumlahPustakawan; i++ {
		wg.Add(1)
		go LihatDaftarBuku(ch, chBuku, &wg)
	}

	for _, fileBuku := range listJsonBuku {
		ch <- fileBuku.Name()
	}

	close(ch)

	wg.Wait()

	close(chBuku)

	for dataBuku := range chBuku {
		ListBuku = append(ListBuku, dataBuku)
	}

	sort.Slice(ListBuku, func(i, j int) bool {
		return ListBuku[i].Waktu.Before(ListBuku[j].Waktu)
	})

	if len(ListBuku) < 1 {
		fmt.Println("===== Tidak ada buku =====")
	}

	for i, v := range ListBuku {
		i++
		fmt.Printf("%d. Kode Buku : %s, Judul Buku : %s, Pengarang Buku : %s, Penerbit Buku : %s, Jumlah Halaman : %d, Tahun Terbit %d\n",
			i, v.Kode, v.Judul, v.Pengarang, v.Penerbit, v.Halaman, v.Tahun)
	}

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func DetailBuku(kode string) {
	fmt.Println("")
	fmt.Println("Detail Buku")
	fmt.Println("")

	var isiBuku bool

	for _, buku := range ListBuku {
		if buku.Kode == kode {
			isiBuku = true
			fmt.Printf("Kode Buku : %s\n", buku.Kode)
			fmt.Printf("Judul Buku : %s\n", buku.Judul)
			fmt.Printf("Pengarang Buku : %s\n", buku.Pengarang)
			fmt.Printf("Penerbit Buku : %s\n", buku.Penerbit)
			fmt.Printf("Jumlah Halaman Buku : %d\n", buku.Halaman)
			fmt.Printf("Tahun Terbit Buku : %d\n", buku.Tahun)
			break
		}
	}

	if !isiBuku {
		fmt.Println("Kode Buku Salah Atau Tidak Ada")
	}

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func HapusBuku(kode string) {

	var isiBuku bool
	for i, buku := range ListBuku {
		if buku.Kode == kode {
			isiBuku = true
			err := os.Remove(fmt.Sprintf("books/%s.json", ListBuku[i].Kode))
			if err != nil {
				fmt.Println("Terjadi Error: ", err)
			}

			fmt.Print("\n")
			fmt.Println("Buku Berhasil Dihapus")
			break
		}
	}

	if !isiBuku {

		fmt.Print("\n")
		fmt.Println("Kode Buku Salah Atau Tidak Ada")
	}

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func GeneratedPdfBuku() {

	LihatList()
	fmt.Println("===== Membuat Daftar Buku =====")
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 12)
	pdf.SetLeftMargin(10)
	pdf.SetRightMargin(10)

	for i, buku := range ListBuku {
		bukuText := fmt.Sprintf(
			"Buku #%d:\nKode Buku : %s\nJudul : %s\nPengarang : %s\nPenerbit : %s\nJumlah Halaman : %d\nTahun Terbit : %d\nWaktu : %s\n",
			i+1, buku.Kode, buku.Judul, buku.Pengarang, buku.Penerbit, buku.Halaman, buku.Tahun,
			buku.Waktu.Format("2006-01-02 15:04:05"))

		pdf.MultiCell(0, 10, bukuText, "0", "L", false)
		pdf.Ln(5)
	}

	err := pdf.OutputFileAndClose(
		fmt.Sprintf("daftar_buku_%s.pdf",
			time.Now().Format("2006-01-02 15-04-05")))

	if err != nil {
		fmt.Println("Terjadi Error: ", err)
	}

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func PrintSelectedBook() {
	LihatList()
	fmt.Print("Masukkan nomor urut buku yang ingin dicetak: ")
	var selectedNumber int
	_, err := fmt.Scanln(&selectedNumber)
	if err != nil {
		fmt.Println("Terjadi Error: ", err)
		return
	}
	if selectedNumber < 1 || selectedNumber > len(ListBuku) {
		fmt.Println("Nomor Urut buku tidak Valid.")
		return
	}
	selectedBook := ListBuku[selectedNumber-1]
	Selected(selectedBook)
}

func Selected(selectedBook Buku) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 12)
	pdf.SetLeftMargin(10)
	pdf.SetRightMargin(10)

	bukuText := fmt.Sprintf(
		"====================================\nKodeBuku : %s\nJudulBuku : %s\nPengarang : %s\nPenerbit : %s\nJumlahHalaman : %d\nTahunTerbit :  %d\nTanggal : %s\n====================================\n",
		selectedBook.Kode, selectedBook.Judul, selectedBook.Pengarang, selectedBook.Penerbit, selectedBook.Halaman, selectedBook.Tahun,
		selectedBook.Waktu.Format("2006-01-02 15:04:05"))

	pdf.MultiCell(0, 10, bukuText, "0", "L", false)

	err := pdf.OutputFileAndClose(
		fmt.Sprintf("data_buku_%s,pdf",
			time.Now().Format("2006-01-02-15-04-05")))

	if err != nil {
		fmt.Println("Terjadi Error: ", err)
	}
	fmt.Println("Buku Berhasil dicetak dalam file PDF.")
}

func EditBuku(kode string) {

	DetailBuku(kode)

	fmt.Println("")
	fmt.Println("Edit Buku")
	fmt.Println("")

	var buku Buku
	fmt.Print("Masukan Kode Buku : ")
	_, err := fmt.Scanln(&buku.Kode)
	if err != nil {
		fmt.Println("Terjadi Error : ", err)
	}

	listJsonBuku, err := os.ReadDir("books")
	if err != nil {
		fmt.Println("Terjadi Error : ", err)
	}

	for _, bukuJson := range listJsonBuku {
		if bukuJson.Name() == "book-"+buku.Kode+".json" {
			fmt.Println("Kode buku sudah ada. masukkan kode yang lain.")
			return
		}
	}

	fmt.Print("Masukan Judul Buku : ")
	_, err = fmt.Scanln(&buku.Judul)
	if err != nil {
		fmt.Println("Terjadi Error : ", err)
	}

	fmt.Print("Masukan Nama Pengarang : ")
	_, err = fmt.Scanln(&buku.Pengarang)
	if err != nil {
		fmt.Println("Terjadi Error : ", err)
	}

	fmt.Print("Masukan Nama Penerbit : ")
	_, err = fmt.Scanln(&buku.Penerbit)
	if err != nil {
		fmt.Println("Terjadi Error : ", err)
	}

	fmt.Print("Masukan Halaman Buku : ")
	_, err = fmt.Scanln(&buku.Halaman)
	if err != nil {
		fmt.Println("Terjadi Error : ", err)
	}

	fmt.Print("Masukan Tahun Terbit : ")
	_, err = fmt.Scanln(&buku.Tahun)
	if err != nil {
		fmt.Println("Terjadi Error : ", err)
	}

	fmt.Print("\nBuku Berhasil Di Edit")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	buku.Kode = "book-" + buku.Kode
	fmt.Println(buku)

	for i, b := range ListBuku {
		if b.Kode == kode {
			ListBuku[i] = buku
			dataJson, err := json.Marshal(ListBuku[i])
			if err != nil {
				fmt.Println("Terjadi Error : ", err)
			}

			err = os.WriteFile(fmt.Sprintf("books/%s.json", ListBuku[i].Kode), dataJson, 0644)
			if err != nil {
				fmt.Println("Terjadi Error : ", err)
			}

			err = os.Remove(fmt.Sprintf("books/%s.json", kode))
			if err != nil {
				fmt.Println("Terjadi Error : ", err)
			}

			break
		}
	}
}

func main() {
	PilihanBuku := 0

	cls.CLS()
	fmt.Println("")
	fmt.Println("Aplikasi Manajemen Daftar Buku Perpustakaan")
	fmt.Println("")
	fmt.Println("Silahkan Pilih : ")
	fmt.Println("1. Tambah Buku")
	fmt.Println("2. Liat List Buku")
	fmt.Println("3. Detail Buku")
	fmt.Println("4. Ubah/Edit Buku")
	fmt.Println("5. Hapus Buku")
	fmt.Println("6. Print semua buku ke pdf")
	fmt.Println("7. Print salah satu buku ke pdf")
	fmt.Println("8. Keluar")
	fmt.Println("")

	fmt.Print("Masukan Pilihan : ")
	_, err := fmt.Scanln(&PilihanBuku)
	if err != nil {
		fmt.Println("Terjadi error:", err)
	}

	switch PilihanBuku {
	case 1:
		TambahBuku()
	case 2:
		LihatList()
	case 3:
		var pilihanDetail string
		LihatList()
		fmt.Print("Masukkan Kode Buku : ")
		_, err := fmt.Scanln(&pilihanDetail)
		if err != nil {
			fmt.Println("Terjadi Error : ", err)
			return
		}
		DetailBuku(pilihanDetail)
	case 4:
		var pilihanEdit string
		LihatList()
		fmt.Print("Masukkan Kode Buku Yang akan diedit : ")
		_, err := fmt.Scanln(&pilihanEdit)
		if err != nil {
			fmt.Println("Terjadi Error : ", err)
			return
		}
		EditBuku(pilihanEdit)
	case 5:
		var pilihanHapus string
		LihatList()
		fmt.Print("masukkan kode yang akan dihapus : ")
		_, err := fmt.Scanln(&pilihanHapus)
		if err != nil {
			fmt.Println("Terjadi error: ", err)
			return
		}
		HapusBuku(pilihanHapus)
	case 6:
		GeneratedPdfBuku()
	case 7:
		PrintSelectedBook()
	case 8:
		fmt.Println("\nSelesai")
		os.Exit(0)
	default:
		fmt.Println("\ntidak ada opsi")
	}

	main()
}

// Fitur pada apliaksi :
// 1.Menambah buku baru dengan informasi
// -Kode Buku string
// -Judul Buku string
// -Pengarang string
// -Penerbit string
// -Jumlah Halaman int
// -Tahun Terbit int
// 2.Menampilkan semua list pada daftar buku di perpustakaan
// 3.Dapat menghapus buku menggunakan Kode Buku
// 4.Dapat mengubah/mengedit buku berdasarkan Kode Buku

// // Fitur pada apliaksi :
// // 1.Menambah buku baru dengan informasi
// // -Kode Buku string
// // -Judul Buku string
// // -Pengarang string
// // -Penerbit string
// // -Jumlah Halaman int
// // -Tahun Terbit int
// // 2.Menampilkan semua list pada daftar buku di perpustakaan
// // 3.Dapat menghapus buku menggunakan Kode Buku
// // 4.Dapat mengubah/mengedit buku berdasarkan Kode Buku
