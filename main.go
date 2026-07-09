package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"strconv"
)

type Navigation struct {
	IdNavigationPlan int
	IndexValue       string
	Description      string
	IndexFather      string
	ItemOrder        int
	DispositionId    int
}

func main() {
	idNavigationPlan := 1

	file, err := excelize.OpenFile("NavigationPlan/CredAluga_AD_Preventivo_Painel de Navegação_v1.0.xlsx")

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
		var navigation Navigation

		for j := 1; j < len(row); j++ {
			var cell = ""

			if allow {
				switch j {
				case 1:
					navigation.IdNavigationPlan = idNavigationPlan
					navigation.Description = row[j]
				case 2:
					navigation.IndexValue = row[j]
				case 3:
					if row[j] == "-" {
						navigation.IndexFather = "Null"
					}

					if row[j] != "-" {
						navigation.IndexFather = fmt.Sprintf("'%s'", row[j])
					}
				case 4:
					v, _ := strconv.Atoi(row[j])
					navigation.ItemOrder = v
				case 5:
					v, _ := strconv.Atoi(row[j])
					navigation.DispositionId = v
				}
				continue
			}

			if !allow {
				cell, _ = file.GetCellValue("Painel de Navegação", fmt.Sprintf("%s%d", "B", i+1))

				if cell == "Description" {
					allow = true
				}
			}

			break
		}

		if allow {
			fmt.Printf("insert into SysConfiguration..NavigationDetail (IdNavigationPlan, IndexValue, Description, Indexfather, ItemOrder, DispositionId) values (%d, '%s', '%s', %s, %d, %d);", navigation.IdNavigationPlan, navigation.IndexValue, navigation.Description, navigation.IndexFather, navigation.ItemOrder, navigation.DispositionId)
		}

		if allow {
			fmt.Println()
		}
	}
}
