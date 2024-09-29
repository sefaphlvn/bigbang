package resources

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sefaphlvn/bigbang/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func DecodeBase64Config(encodedConfig string) (*models.TypedConfig, error) {
	decodedString, err := base64.StdEncoding.DecodeString(encodedConfig)
	if err != nil {
		return nil, err
	}

	var configData models.TypedConfig
	err = json.Unmarshal(decodedString, &configData)
	if err != nil {
		return nil, err
	}

	return &configData, nil
}

func GetTypedConfigValue(jsonStringStr string, path string, logger *logrus.Logger) *models.TypedConfig {
	value := gjson.Get(jsonStringStr, path).String()

	if value == "" {
		logger.Debugf("typed_config value empty for path: %s", path)
		return nil
	}

	typedConfig, err := DecodeBase64Config(value)
	if err != nil {
		logger.Debugf("Error decoding base64 config: %v", err)
		return nil
	}

	return typedConfig
}

func ProcessTypedConfigs(jsonStringStr string, typedConfigPath models.TypedConfigPath, logger *logrus.Logger) ([]*models.TypedConfig, map[string]*models.TypedConfig) {
	var typedConfigs []*models.TypedConfig
	typedConfigsMap := make(map[string]*models.TypedConfig)
	seenConfigs := make(map[string]struct{})

	if typedConfigPath.IsPerTypedConfig {
		if len(typedConfigPath.ArrayPaths) == 0 {
			result := gjson.Get(jsonStringStr, typedConfigPath.PathTemplate)
			if result.Exists() {
				result.ForEach(func(key, value gjson.Result) bool {
					dynamicKey := key.String()
					dynamicPath := fmt.Sprintf("%s.%s", typedConfigPath.PathTemplate, dynamicKey)
					processPath(jsonStringStr, dynamicPath, &typedConfigs, typedConfigsMap, seenConfigs, logger)
					return true
				})
			}
		} else {
			result := gjson.Get(jsonStringStr, typedConfigPath.ArrayPaths[0].ParentPath)
			if result.IsArray() {
				processPerTypedConfigArray(result.Array(), jsonStringStr, typedConfigPath.PathTemplate, typedConfigPath.ArrayPaths, &typedConfigs, typedConfigsMap, seenConfigs, logger)
			} else if result.Exists() {
				processDynamicKey(result, typedConfigPath.ArrayPaths[0].ParentPath, &typedConfigs, typedConfigsMap, seenConfigs, logger)
			}
		}
	} else {
		if len(typedConfigPath.ArrayPaths) == 0 {
			processPath(jsonStringStr, typedConfigPath.PathTemplate, &typedConfigs, typedConfigsMap, seenConfigs, logger)
		} else {
			result := gjson.Get(jsonStringStr, typedConfigPath.ArrayPaths[0].ParentPath)
			if result.IsArray() {
				processArray(result.Array(), jsonStringStr, typedConfigPath.PathTemplate, typedConfigPath.ArrayPaths, &typedConfigs, typedConfigsMap, seenConfigs, logger)
			} else if result.Exists() {
				processPath(result.String(), typedConfigPath.ArrayPaths[0].ParentPath, &typedConfigs, typedConfigsMap, seenConfigs, logger)
			}
		}
	}

	return typedConfigs, typedConfigsMap
}

// Dynamic key ile path işlemlerini işleyen fonksiyon
func processDynamicKey(result gjson.Result, basePath string, typedConfigs *[]*models.TypedConfig, typedConfigsMap map[string]*models.TypedConfig, seenConfigs map[string]struct{}, logger *logrus.Logger) {
	result.ForEach(func(key, value gjson.Result) bool {
		dynamicKey := key.String()
		dynamicPath := fmt.Sprintf("%s.%s", basePath, dynamicKey)
		processPath(result.String(), dynamicPath, typedConfigs, typedConfigsMap, seenConfigs, logger)
		return true
	})
}

// PerTypedConfig için dizi elemanlarını işleyen yardımcı fonksiyon
func processPerTypedConfigArray(array []gjson.Result, jsonStringStr string, pathTemplate string, arrayPaths []models.ArrayPath, typedConfigs *[]*models.TypedConfig, typedConfigsMap map[string]*models.TypedConfig, seenConfigs map[string]struct{}, logger *logrus.Logger) {
	placeholderCount := strings.Count(pathTemplate, "%d")

	for i := range array {
		combinations := generateIndexCombinations(jsonStringStr, arrayPaths)

		for _, indices := range combinations {
			if len(indices) == placeholderCount {
				indices[0] = i
				finalPath := fmt.Sprintf(pathTemplate, indices...)
				dynamicResult := gjson.Get(jsonStringStr, finalPath)

				if dynamicResult.Exists() {
					dynamicResult.ForEach(func(key, value gjson.Result) bool {
						dynamicKey := key.String()
						dynamicPath := fmt.Sprintf("%s.%s", finalPath, dynamicKey)
						processPath(jsonStringStr, dynamicPath, typedConfigs, typedConfigsMap, seenConfigs, logger)
						return true
					})
				}
			}
		}
	}
}

// Dizi elemanlarını işleyen yardımcı fonksiyon
func processArray(array []gjson.Result, jsonStringStr string, pathTemplate string, arrayPaths []models.ArrayPath, typedConfigs *[]*models.TypedConfig, typedConfigsMap map[string]*models.TypedConfig, seenConfigs map[string]struct{}, logger *logrus.Logger) {
	placeholderCount := strings.Count(pathTemplate, "%d")

	// İlk seviyedeki diziyi işliyoruz; bu seviyenin eleman sayısı kadar döngü yapacağız
	for i := range array {
		// İndeks kombinasyonlarını dinamik olarak üret
		combinations := generateIndexCombinations(jsonStringStr, arrayPaths)

		// Üretilen kombinasyonlar üzerinden path'leri oluşturup işle
		for _, indices := range combinations {
			if len(indices) == placeholderCount {
				// İlk seviyenin indexi i olarak belirleniyor
				indices[0] = i
				// Path template'e indeksleri uygulayarak nihai path'i oluştur
				processPath(jsonStringStr, fmt.Sprintf(pathTemplate, indices...), typedConfigs, typedConfigsMap, seenConfigs, logger)
			}
		}
	}
}

// JSON path üzerinden işlem yaparak config değerlerini ekler
func processPath(jsonStringStr, path string, typedConfigs *[]*models.TypedConfig, typedConfigsMap map[string]*models.TypedConfig, seenConfigs map[string]struct{}, logger *logrus.Logger) {
	singleTypedConfig := GetTypedConfigValue(jsonStringStr, path+".value", logger)

	if singleTypedConfig != nil {
		// Benzersizlik kontrolü için benzersiz bir anahtar oluşturuyoruz (örneğin, type_url ve name kullanarak)
		uniqueKey := fmt.Sprintf("%s|%s|%s", singleTypedConfig.Gtype, singleTypedConfig.Name, path)

		// Benzersizlik kontrolü
		if _, exists := seenConfigs[uniqueKey]; !exists {
			*typedConfigs = append(*typedConfigs, singleTypedConfig)
			typedConfigsMap[path] = singleTypedConfig
			seenConfigs[uniqueKey] = struct{}{} // Eklendi olarak işaretle
		} else {
			logger.Debugf("Duplicate typed_config detected for key: %s", uniqueKey)
		}
	}
}

// Birden fazla %d kombinasyonlarını üreten yardımcı fonksiyon
func generateIndexCombinations(jsonStringStr string, arrayPaths []models.ArrayPath) [][]interface{} {
	var combinations [][]interface{}

	if len(arrayPaths) == 0 {
		return combinations // Eğer array tanımlı değilse, boş dön
	}

	indices := make([]interface{}, len(arrayPaths))
	generateCombinations(arrayPaths, indices, 0, &combinations, jsonStringStr)

	return combinations
}

func generateCombinations(arrayPaths []models.ArrayPath, indices []interface{}, level int, combinations *[][]interface{}, jsonStringStr string) {
	if level == len(indices) {
		// Level en üst seviyeye ulaştığında kombinasyonu ekle
		*combinations = append(*combinations, append([]interface{}(nil), indices...))
		return
	}

	// Path'leri dinamik olarak doldur
	parentPath := fillIndices(arrayPaths[level].ParentPath, indices[:level])

	// Mevcut seviyedeki diziyi al
	currentArray := gjson.Get(jsonStringStr, parentPath)
	if !currentArray.IsArray() {
		// Eğer mevcut path bir array değilse, işlemi sonlandır
		fmt.Printf("Warning: Expected array at path %s, but did not find an array.\n", parentPath)
		return
	}

	// Mevcut dizinin boyutunu kontrol ediyoruz ve her eleman için recursive çağrı yapıyoruz
	for i := 0; i < len(currentArray.Array()); i++ {
		indices[level] = i

		// Bir sonraki seviyeye geç ve dizileri işle
		generateCombinations(arrayPaths, indices, level+1, combinations, jsonStringStr)
	}
}

func fillIndices(pathTemplate string, indices []interface{}) string {
	// Doldurulmamış tüm %d yer tutucularını mevcut indeksler ile doldurur
	return fmt.Sprintf(pathTemplate, indices...)
}
