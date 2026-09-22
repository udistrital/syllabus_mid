package services

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/mfpierre/go-mcrypt"
	"github.com/udistrital/syllabus_mid/utils"
	"github.com/udistrital/utils_oas/request"
	"github.com/udistrital/utils_oas/requestresponse"
)

func GetSyllabusLegacy(encodedParamsPlan string) requestresponse.APIResponse {
	var syllabusTemplateData map[string]interface{}
	var syllabusData map[string]interface{}
	var syllabusResponse map[string]interface{}

	paramsPlans, paramsError := decodeParamsPlan(encodedParamsPlan)
	if paramsError != nil {
		logs.Error(paramsError.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, paramsError.Error())
	}

	// Get map of params
	paramsMap, paramsMapError := paramsString2Map(paramsPlans)
	if paramsMapError != nil {
		logs.Error(paramsMapError.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, paramsMapError.Error())
	}
	fmt.Println(paramsMap)

	// Get params
	planEstudioString, planEstudioOK := paramsMap["planEstudio"]
	if !planEstudioOK {
		err := fmt.Errorf("params without planEstudio")
		logs.Error(err.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}
	planEstudioId, planEstudioError := strconv.ParseInt(planEstudioString.(string), 10, 64)
	if planEstudioError != nil {
		logs.Error(planEstudioError.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, planEstudioError.Error())
	}

	proyectoCurricularString, proyectoCurricularOK := paramsMap["codProyecto"]
	if !proyectoCurricularOK {
		err := fmt.Errorf("params without Proyecto Curricular")
		logs.Error(err.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}
	proyectoCurricularId, proyectoCurricularError := strconv.ParseInt(proyectoCurricularString.(string), 10, 64)
	if proyectoCurricularError != nil {
		logs.Error(proyectoCurricularError.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, proyectoCurricularError.Error())
	}

	espacioAcademicoString, espacioAcademicoOK := paramsMap["codEspacio"]
	if !espacioAcademicoOK {
		err := fmt.Errorf("params without Espacio Académico")
		logs.Error(err.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}
	espacioAcademicoId, espacioAcademicoError := strconv.ParseInt(espacioAcademicoString.(string), 10, 64)
	if espacioAcademicoError != nil {
		logs.Error(espacioAcademicoError.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, espacioAcademicoError.Error())
	}

	// Query
	syllabusErr := request.GetJson(beego.AppConfig.String("SyllabusService")+
		fmt.Sprintf("syllabus?query=espacio_academico_id:%v,proyecto_curricular_id:%v,plan_estudios_id:%v,syllabus_actual:true",
			espacioAcademicoId, proyectoCurricularId, planEstudioId),
		&syllabusResponse)
	if syllabusErr != nil || syllabusResponse["Success"] == false {
		if syllabusErr == nil {
			syllabusErr = fmt.Errorf("SyllabusService: %v", syllabusResponse["Message"])
		}
		logs.Error(syllabusErr.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, syllabusErr.Error())
	}
	syllabusList := syllabusResponse["Data"].([]interface{})
	if len(syllabusList) == 0 {
		err := fmt.Errorf("SyllabusService: No syllabus found")
		logs.Error(err.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}
	syllabusData = syllabusList[0].(map[string]interface{})

	spaceData, spaceErr := utils.GetAcademicSpaceData(
		int(planEstudioId),
		int(proyectoCurricularId),
		int(espacioAcademicoId))

	if spaceErr != nil {
		logs.Error(spaceErr.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, spaceErr.Error())
	}

	projectData, projectErr := utils.GetProyectoCurricular(int(proyectoCurricularId))

	if projectErr != nil {
		logs.Error(projectErr.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, projectErr.Error())
	}

	facultyData, facultyErr := utils.GetFacultadDelProyectoC(projectData["id_oikos"].(string))
	idiomas := ""

	if syllabusData["idioma_espacio_id"] != nil {
		idiomasStr, idiomaErr := utils.GetIdiomas(syllabusData["idioma_espacio_id"].([]interface{}))
		if idiomaErr == nil {
			idiomas = idiomasStr
		}
	}

	if syllabusData["plan_estudios_id"] == nil {
		syllabusData["plan_estudios_id"] = planEstudioString
	}

	if facultyErr == nil {
		syllabusTemplateData = utils.GetSyllabusTemplateData(
			spaceData, syllabusData,
			facultyData, projectData, idiomas)
		fmt.Println(syllabusTemplateData)
		return requestresponse.APIResponseDTO(true, 200, syllabusTemplateData, "Syllabus OK")
	} else {
		err := fmt.Errorf(
			"SyllabusService: Incomplete data. Facultad y/o Idioma")
		logs.Error(err.Error())
		return requestresponse.APIResponseDTO(false, 404, nil, err.Error())
	}

}

func decodeParamsPlan(encryptedParamsString string) (string, error) {
	key := []byte(beego.AppConfig.String("SyllabusSeed"))
	iv := make([]byte, 32)

	// Preprocessing of the encrypted params

	encryptedString := strings.Replace(encryptedParamsString, "-", "+", -1)
	encryptedString = strings.Replace(encryptedString, "_", "/", -1)
	complement := len(encryptedString) % 4
	if complement != 0 {
		encryptedString += strings.Repeat("=", complement)
	}
	encryptedBase64String, _ := base64.StdEncoding.DecodeString(encryptedString)

	decrypted, _ := mcrypt.Decrypt(key, iv, []byte(encryptedBase64String), "rijndael-256", "ecb")
	if len(decrypted) == 0 {
		return "", fmt.Errorf("wrong decryption")
	}
	return string(decrypted), nil
}

func paramsString2Map(paramsString string) (map[string]interface{}, error) {
	paramsMap := make(map[string]interface{})
	paramsList := strings.Split(paramsString, "&")

	if len(paramsList) < 3 {
		return nil, fmt.Errorf("wrong params, incomplete parameters")
	}
	for _, param := range paramsList {
		paramSplit := strings.Split(param, "=")
		if len(paramSplit) != 2 {
			return nil, fmt.Errorf("wrong params, parameters without value")
		}
		paramsMap[paramSplit[0]] = paramSplit[1]
	}
	return paramsMap, nil
}
