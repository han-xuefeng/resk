package base

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	vtzh "github.com/go-playground/validator/v10/translations/zh"
	log "github.com/sirupsen/logrus"
	"study-gin/resk/infra"
)

var validate *validator.Validate

var translator ut.Translator

//验证器
func Validate() *validator.Validate {
	return validate
}

// 翻译器
func Transtate() ut.Translator {
	return translator
}

type ValidatorStarter struct {
	infra.BaseStarter
}

func (v *ValidatorStarter) Init(ctx infra.StarterContext) {
	validate = validator.New()
	//创建消息国际化通用翻译器
	cn := zh.New()
	uni := ut.New(cn, cn)
	var found bool
	translator, found = uni.GetTranslator("zh")
	if found {
		err := vtzh.RegisterDefaultTranslations(validate, translator)
		if err != nil {
			log.Error(err)
		}
	} else {
		log.Error("Not found translator: zh")
	}
}

func ValidateStruct(s interface{}) (err error) {
	err = Validate().Struct(s)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			log.Error("验证错误", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, e := range errs {
				log.Error(e.Translate(Transtate()))
			}
		}
		return err
	}
	return nil
}