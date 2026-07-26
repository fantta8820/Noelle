package script

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Navigation struct {
	IdNavigationPlan string
	IndexValue       string
	Description      string
	IndexFather      string
	ItemOrder        int
	DispositionId    int
}

func RunScript(path string, navigationPlanName string, campaignId int) (string, error) {
	rows, file, err := OpenFile(path)

	if err != "" {
		return "", errors.New(err)
	}

	allow := false
	descriptionIndex := 0

	return GenerateSQL(allow, rows, descriptionIndex, navigationPlanName, campaignId, file), nil
}

func OpenFile(path string) ([][]string, *excelize.File, string) {
	file, err := excelize.OpenFile(path)

	if err != nil {
		return nil, nil, "Houve um erro ao abrir o arquivo."
	}

	rows, err := file.GetRows("Painel de Navegação")

	if err != nil {
		return nil, nil, "Houve um erro ao carregar as informações do arquivo."
	}

	if err := file.Close(); err != nil {
		return nil, nil, "Houve um erro ao fechar o arquivo."
	}

	return rows, file, ""
}

func GenerateSQL(allow bool, rows [][]string, descriptionIndex int, navigationPlanName string, campaignId int, file *excelize.File) string {
	query := `DECLARE @IdNavigationPlan int;
BEGIN TRY
	BEGIN TRANSACTION
		select top 1 @IdNavigationPlan = (IdNavigationPlan + 1) from SysConfiguration..NavigationPlan with(nolock) order by IdNavigationPlan desc;
		insert into SysConfiguration..NavigationPlan (IdNavigationPlan, Name, LastUpdateUserId, LastUpdate, IsActive, InsertDate) values(@IdNavigationPlan, '{Name}', 1, GETDATE(), 1, GETDATE())
		insert into SysConfiguration..NavigationPlanCampaign values({CampaignId}, @IdNavigationPlan);

`

	for i, row := range rows {
		var navigation Navigation

		for j := 1; j < len(row); j++ {
			var cell = ""

			if allow {
				switch j {
				case 1:
					navigation.IdNavigationPlan = "@IdNavigationPlan"
					navigation.Description = row[j]
				case 2:
					if len(row[j]) < 2 {
						navigation.IndexValue = fmt.Sprintf("0%s", row[j])
					} else {
						navigation.IndexValue = row[j]
					}
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
					descriptionIndex = i
					allow = true
				}
			}

			break
		}

		if allow && i > descriptionIndex {
			query += fmt.Sprintf("		insert into SysConfiguration..NavigationDetail (IdNavigationPlan, IndexValue, Description, Indexfather, ItemOrder, DispositionId) values (%s, '%s', '%s', %s, %d, %d);", navigation.IdNavigationPlan, navigation.IndexValue, navigation.Description, navigation.IndexFather, navigation.ItemOrder, navigation.DispositionId)
		}

		if allow && i > descriptionIndex && i < len(rows)-1 {
			query += "\n"
		}
	}

	query += `
	COMMIT;
END TRY
BEGIN CATCH
	ROLLBACK;
	THROW;
END CATCH;`

	query = strings.ReplaceAll(query, "{Name}", navigationPlanName)
	query = strings.ReplaceAll(query, "{CampaignId}", strconv.Itoa(campaignId))

	return query
}
