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

	for _, row := range rows {
		for _, col := range row {
			fmt.Print(col, "\t")
		}
		fmt.Println()
	}
}
