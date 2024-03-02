package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Buku struct {
	Kode, Judul, Pengarang, Penerbit string
	Halaman, Tahun                   int
}

var ListBuku []Buku

func TambahBuku() {
	KodeB := bufio.NewReader(os.Stdin)
	JudulB := bufio.NewReader(os.Stdin)
	PengarangB := bufio.NewReader(os.Stdin)
	PenerbitB := bufio.NewReader(os.Stdin)
	JumlahHalaman := 0
	TahunTerbit := 0
	fmt.Println("")
	fmt.Println("Tambahkan Buku")
	fmt.Println("")
	fmt.Print("Masukan Kode Buku : ")
	KodeBuku, err := KodeB.ReadString('\n')
	if err != nil {
		fmt.Println("Terjadi Error:", err)
		return
	}
	KodeBuku = strings.TrimSpace(KodeBuku)

	fmt.Print("Masukan Judul Buku : ")
	JudulBuku, err := JudulB.ReadString('\n')
	if err != nil {
		fmt.Println("Terjadi Error:", err)
		return
	}
	JudulBuku = strings.TrimSpace(JudulBuku)

	fmt.Print("Masukan Nama Pengarang : ")
	NamaPengarang, err := PengarangB.ReadString('\n')
	if err != nil {
		fmt.Println("Terjadi Error:", err)
		return
	}
	NamaPengarang = strings.TrimSpace(NamaPengarang)

	fmt.Print("Masukan Nama Penerbit : ")
	NamaPenerbit, err := PenerbitB.ReadString('\n')
	if err != nil {
		fmt.Println("Terjadi Error:", err)
		return
	}
	NamaPenerbit = strings.TrimSpace(NamaPenerbit)

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

	ListBuku = append(ListBuku, Buku{
		Kode:      KodeBuku,
		Judul:     JudulBuku,
		Pengarang: NamaPengarang,
		Penerbit:  NamaPenerbit,
		Halaman:   JumlahHalaman,
		Tahun:     TahunTerbit,
	})

	fmt.Println("Berhasil Menambahkan Buku")
}

func LihatList() {
	fmt.Println("")
	fmt.Println("Lihat List Buku")
	fmt.Println("")

	if len(ListBuku) == 0 {
		fmt.Println("Tidak Ada Buku Yang Tersimpan")
		return
	}
	for urutan, Buku := range ListBuku {
		fmt.Printf("\n %d. Kode Buku : %s \n- Judul : %s \n- Pengarang : %s \n- Penerbit : %s \n- Jumlah Halaman : %d  \n- Tahun Rilis : %d \n",
			urutan+1,
			Buku.Kode,
			Buku.Judul,
			Buku.Pengarang,
			Buku.Penerbit,
			Buku.Halaman,
			Buku.Tahun)
	}
}

func main() {
	PilihanBuku := 0

	fmt.Println("")
	fmt.Println("Aplikasi Manajemen Daftar Buku Perpustakaan")
	fmt.Println("")
	fmt.Println("Silahkan Pilih : ")
	fmt.Println("1. Tambah Buku")
	fmt.Println("2. Liat List Buku")
	fmt.Println("3. Hapus Buku")
	fmt.Println("4. Ubah/Edit Buku")
	fmt.Println("5. Keluar")
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
		HapusBuku()
	case 4:
		EditBuku()
	case 5:
		fmt.Println("\nSelesai")
		os.Exit(0)
	}

	main()
}

func HapusBuku() {
	var UrutanBuku string
	fmt.Println("")
	fmt.Println("Hapus Buku")
	fmt.Println("")
	LihatList()
	fmt.Println("")

	fmt.Print("Masukan Kode Buku : ")
	fmt.Scanln(&UrutanBuku)

	for urutan, j := range ListBuku {
		if j.Kode == UrutanBuku {
			ListBuku = append(
				ListBuku[:urutan],
				ListBuku[urutan+1:]...,
			)
			fmt.Println("Buku Berhasil Dihapus!")
		} else {
			fmt.Println("Terjadi Kesalahan")
		}
	}
}

func EditBuku() {
	//inputedit := bufio.NewReader(os.Stdin)
	var KodeEdit string

	fmt.Println("")
	fmt.Println("List Buku")
	fmt.Println("")
	LihatList()

	fmt.Print("\nMasukan Kode Buku Yang Ingin Di Edit : ")
	fmt.Scanln(&KodeEdit)

	for i, Buku := range ListBuku {
		if Buku.Kode == KodeEdit {
			inputedit := bufio.NewReader(os.Stdin)

			fmt.Print("Masukan Judul Buku : ")
			ListBuku[i].Judul, _ = inputedit.ReadString('\n')

			fmt.Print("Masukan Nama Pengarang : ")
			ListBuku[i].Pengarang, _ = inputedit.ReadString('\n')

			fmt.Print("Masukan Nama Penerbit : ")
			ListBuku[i].Penerbit, _ = inputedit.ReadString('\n')

			fmt.Print("Masukan Halaman Buku : ")
			fmt.Scanln(&ListBuku[i].Halaman)

			fmt.Print("Masukan Tahun Terbit : ")
			fmt.Scanln(&ListBuku[i].Tahun)

			fmt.Print("\nBuku Berhasil Di Edit")

			return
		} else {
			fmt.Println("Terjadi Error")

		}
	}
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
