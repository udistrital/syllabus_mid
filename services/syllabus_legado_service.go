package services

import (
	"encoding/base64"
	"fmt"
	"net/http"
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
		return requestresponse.APIResponseDTO(false, http.StatusBadRequest, nil, paramsError.Error())
	}

	// Get map of params
	paramsMap, paramsMapError := paramsString2Map(paramsPlans)
	if paramsMapError != nil {
		logs.Error(paramsMapError.Error())
		return requestresponse.APIResponseDTO(false, http.StatusBadRequest, nil, paramsMapError.Error())
	}

	// Get params
	planEstudioString, planEstudioOK := paramsMap["planEstudio"]
	if !planEstudioOK {
		err := fmt.Errorf("params without planEstudio")
		logs.Error(err.Error())
		return requestresponse.APIResponseDTO(false, http.StatusBadRequest, nil, err.Error())
	}
	planEstudioId, planEstudioError := strconv.ParseInt(planEstudioString.(string), 10, 64)
	if planEstudioError != nil {
		logs.Error(planEstudioError.Error())
		return requestresponse.APIResponseDTO(false, http.StatusBadRequest, nil, planEstudioError.Error())
	}

	proyectoCurricularString, proyectoCurricularOK := paramsMap["codProyecto"]
	if !proyectoCurricularOK {
		err := fmt.Errorf("params without Proyecto Curricular")
		logs.Error(err.Error())
		return requestresponse.APIResponseDTO(false, http.StatusBadRequest, nil, err.Error())
	}
	proyectoCurricularId, proyectoCurricularError := strconv.ParseInt(proyectoCurricularString.(string), 10, 64)
	if proyectoCurricularError != nil {
		logs.Error(proyectoCurricularError.Error())
		return requestresponse.APIResponseDTO(false, http.StatusBadRequest, nil, proyectoCurricularError.Error())
	}

	espacioAcademicoString, espacioAcademicoOK := paramsMap["codEspacio"]
	if !espacioAcademicoOK {
		err := fmt.Errorf("params without Espacio Académico")
		logs.Error(err.Error())
		return requestresponse.APIResponseDTO(false, http.StatusBadRequest, nil, err.Error())
	}
	espacioAcademicoId, espacioAcademicoError := strconv.ParseInt(espacioAcademicoString.(string), 10, 64)
	if espacioAcademicoError != nil {
		logs.Error(espacioAcademicoError.Error())
		return requestresponse.APIResponseDTO(false, http.StatusBadRequest, nil, espacioAcademicoError.Error())
	}

	// Query
	syllabusErr := request.GetJson(beego.AppConfig.String("SyllabusService")+
		fmt.Sprintf("syllabus?query=espacio_academico_id:%v,proyecto_curricular_id:%v,plan_estudios_id:%v,syllabus_actual:true&limit=1",
			espacioAcademicoId, proyectoCurricularId, planEstudioId),
		&syllabusResponse)
	if syllabusErr != nil || syllabusResponse["Success"] == false {
		if syllabusErr == nil {
			syllabusErr = fmt.Errorf("SyllabusService: %v", syllabusResponse["Message"])
		}
		logs.Error(syllabusErr.Error())
		return requestresponse.APIResponseDTO(false, http.StatusNotFound, nil, syllabusErr.Error())
	}
	syllabusList, syllabusListOK := syllabusResponse["Data"].([]interface{})
	if !syllabusListOK || len(syllabusList) == 0 {
		err := fmt.Errorf("SyllabusService: No syllabus found")
		logs.Error(err.Error())
		return requestresponse.APIResponseDTO(false, http.StatusNotFound, nil, err.Error())
	}
	syllabusItem, syllabusItemOK := syllabusList[0].(map[string]interface{})
	if !syllabusItemOK {
		err := fmt.Errorf("SyllabusService: syllabus con formato inesperado")
		logs.Error(err.Error())
		return requestresponse.APIResponseDTO(false, http.StatusNotFound, nil, err.Error())
	}
	syllabusData = syllabusItem

	spaceData, spaceErr := utils.GetAcademicSpaceData(
		int(planEstudioId),
		int(proyectoCurricularId),
		int(espacioAcademicoId))

	if spaceErr != nil {
		logs.Error(spaceErr.Error())
		return requestresponse.APIResponseDTO(false, http.StatusBadGateway, nil, spaceErr.Error())
	}

	projectData, projectErr := utils.GetProyectoCurricular(int(proyectoCurricularId))

	if projectErr != nil {
		logs.Error(projectErr.Error())
		return requestresponse.APIResponseDTO(false, http.StatusBadGateway, nil, projectErr.Error())
	}

	idOikos, idOikosOK := projectData["id_oikos"].(string)
	if !idOikosOK {
		idOikos = fmt.Sprintf("%v", projectData["id_oikos"])
	}
	facultyData, facultyErr := utils.GetFacultadDelProyectoC(idOikos)

	idiomas := ""
	if idiomaIDs, idiomaOK := syllabusData["idioma_espacio_id"].([]interface{}); idiomaOK {
		idiomasStr, idiomaErr := utils.GetIdiomas(idiomaIDs)
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
		return requestresponse.APIResponseDTO(true, http.StatusOK, syllabusTemplateData, "Syllabus OK")
	} else {
		err := fmt.Errorf(
			"SyllabusService: Incomplete data. Facultad y/o Idioma")
		logs.Error(err.Error())
		return requestresponse.APIResponseDTO(false, http.StatusBadRequest, nil, err.Error())
	}

}

func decodeParamsPlan(encryptedParamsString string) (string, error) {
	seed := beego.AppConfig.String("SyllabusSeed")
	if len(seed) != 32 {
		return "", fmt.Errorf("wrong decryption: SyllabusSeed inválido (%d caracteres, se esperan 32)", len(seed))
	}
	key := []byte(seed)
	iv := make([]byte, 32)

	// Preprocessing of the encrypted params: base64 URL-safe, con o sin padding
	normalized := strings.ReplaceAll(encryptedParamsString, "-", "+")
	normalized = strings.ReplaceAll(normalized, "_", "/")
	normalized = strings.TrimRight(normalized, "=")

	encryptedBase64String, err := base64.RawStdEncoding.DecodeString(normalized)
	if err != nil {
		return "", fmt.Errorf("wrong decryption: base64 inválido: %v", err)
	}

	decrypted, err := mcrypt.Decrypt(key, iv, encryptedBase64String, "rijndael-256", "ecb")
	if err != nil {
		return "", fmt.Errorf("wrong decryption: %v", err)
	}
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
		paramSplit := strings.SplitN(param, "=", 2)
		if len(paramSplit) != 2 || paramSplit[0] == "" {
			return nil, fmt.Errorf("wrong params, parameters without value")
		}
		paramsMap[paramSplit[0]] = paramSplit[1]
	}
	return paramsMap, nil
}
