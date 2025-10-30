package lang

import "strings"

var create = map[string]map[string]string{
	"en": {"message": "RESOURCE {resource} CREATED SUCCESSFULLY"},
	"es": {"message": "RECURSO {resource} CREADO EXITOSAMENTE"},
	"pt": {"message": "RECURSO {resource} CRIADO COM SUCESSO"},
}

var retrieve = map[string]map[string]string{
	"en": {"message": "RESOURCE {resource} RETRIEVED SUCCESSFULLY"},
	"es": {"message": "RECURSO {resource} RECUPERADO EXITOSAMENTE"},
	"pt": {"message": "RECURSO {resource} RECUPERADO COM SUCESSO"},
}

var update = map[string]map[string]string{
	"en": {"message": "RESOURCE {resource} UPDATED SUCCESSFULLY"},
	"es": {"message": "RECURSO {resource} ACTUALIZADO EXITOSAMENTE"},
	"pt": {"message": "RECURSO {resource} ATUALIZADO COM SUCESSO"},
}

var delete = map[string]map[string]string{
	"en": {"message": "RESOURCE {resource} DELETED SUCCESSFULLY"},
	"es": {"message": "RECURSO {resource} ELIMINADO EXITOSAMENTE"},
	"pt": {"message": "RECURSO {resource} DELETADO COM SUCESSO"},
}

var failCreate = map[string]map[string]string{
	"en": {"message": "FAILED TO CREATE {resource}"},
	"es": {"message": "FALLO AL CREAR EL {resource}"},
	"pt": {"message": "FALHA AO CRIAR O {resource}"},
}

var failRetrieve = map[string]map[string]string{
	"en": {"message": "FAILED TO RETRIEVE {resource}"},
	"es": {"message": "FALLO AL RECUPERAR EL {resource}"},
	"pt": {"message": "FALHA AO RECUPERAR O {resource}"},
}

var failUpdate = map[string]map[string]string{
	"en": {"message": "FAILED TO UPDATE {resource}"},
	"es": {"message": "FALLO AL ACTUALIZAR EL {resource}"},
	"pt": {"message": "FALHA AO ATUALIZAR O {resource}"},
}

var failDelete = map[string]map[string]string{
	"en": {"message": "FAILED TO DELETE {resource}"},
	"es": {"message": "FALLO AL ELIMINAR EL {resource}"},
	"pt": {"message": "FALHA AO DELETAR O {resource}"},
}

func BuildMessage(resource string, messageId int, lang string) string {

	var template = ""

	if lang == "" {
		lang = "en"
	}

	switch messageId {
	// Resource Created Successfully
	case 1003:
		template = create[lang]["message"]
	// Resource Retrieved Successfully
	case 1004:
		template = retrieve[lang]["message"]
	// Resource Updated Successfully
	case 1005:
		template = update[lang]["message"]
	// Resource Deleted Successfully
	case 1006:
		template = delete[lang]["message"]
	// Failed to Create Resource
	case 2011:
		template = failCreate[lang]["message"]
	// Failed to Retrieve Resource
	case 2012:
		template = failRetrieve[lang]["message"]
	// Failed to Update Resource
	case 2013:
		template = failUpdate[lang]["message"]
	// Failed to Delete Resource
	case 2014:
		template = failDelete[lang]["message"]
	}

	result := strings.Replace(template, "{resource}", resource, 1)

	return strings.ToUpper(result)
}
