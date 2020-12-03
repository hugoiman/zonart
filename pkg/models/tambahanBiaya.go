package models

// TambahanBiaya is class
type TambahanBiaya struct {
	IDTambahanBiaya int    `json:"idTambahanBiaya"`
	Item            string `json:"item"`
	Nominal         int    `json:"nominal"`
	Berat           int    `json:"berat"`
}
