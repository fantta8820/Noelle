package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
)

func main() {
	file, err := excelize.OpenFile("NavigationPlan/MLGOMES_AD_MLGOMES_COBRANÇA_HONDA_Painel de Navegação_v1.0.xlsx")

	if err != nil {
		log.Fatal("Houve um erro ao abrir o arquivo: ", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal("Houve um erro ao fechar o arquivo: ", err)
		}
	}()

	rows, err := file.GetRows("Painel de Navegação")

	if err != nil {
		log.Fatal("Houve um erro ao carregar as informações do seu arquivo: ", err)
	}

	allow := false

	for i, row := range rows {
		for _, col := range row {
			var cell = ""

			if !allow {
				cell, err = file.GetCellValue("Painel de Navegação", fmt.Sprintf("%s%d", "B", i+1))
			}

			if allow {
				fmt.Print(col, "\t")
				continue
			}

			if cell == "Description" {
				allow = true
			}

			break
		}

		if allow {
			fmt.Println()
		}
	}
}
