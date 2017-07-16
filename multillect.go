package multillect

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/arteev/go-translate"
	"github.com/arteev/go-translate/translator"
)

const (
	PROVIDER_CODE = "multillect"
	URL           = "http://api.multillect.com/translate/json/1.0"
)

type ProviderMultillect struct {
	apikey    string
	accountID string
}

type response struct {
	Result *struct {
		Translated string
	}
	Error *struct {
		Code    int
		Message string
	}
}

func (p *ProviderMultillect) urlroute() string {
	return URL + "/" + p.accountID
}

func (ProviderMultillect) GetLangs(code string) ([]*translator.Language, error) {
	return nil, translator.ErrUnsupported
}

func (ProviderMultillect) Detect(text string) (*translator.Language, error) {
	return nil, translator.ErrUnsupported
}
func DecodeDirection(direction string) (string, string) {
	vdir := strings.Split(direction, "-")

	if len(vdir) <= 1 {
		return "", direction
	} else {
		return vdir[0], vdir[1]
	}
}

func (p *ProviderMultillect) Translate(text, direction string) *translator.Result {

	from, to := DecodeDirection(direction)
	r, err := http.PostForm(p.urlroute(), url.Values{
		"method":  {"translate/api/translate"},
		"from":    {from},
		"to":      {to},
		"text":    {text},
		"options": {"1"},
		"sig":     {p.apikey},
	})
	if err != nil {
		return &translator.Result{Err: err}
	}
	defer r.Body.Close()
	var resp response
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return &translator.Result{Err: err}
	}
	/*if err:=p.DecodeApiError(response.Code,200,response.Message);err!=nil{
		return &translator.Result{Err:err}
	}*/
	if resp.Error != nil {
		if resp.Error.Message == "Invalid signature" {
			return &translator.Result{Err: translator.ErrWrongApiKey}
		}
		return &translator.Result{Err: fmt.Errorf("Bad request: (%d) %s", resp.Error.Code, resp.Error.Message)}
	}
	if resp.Result == nil {
		return &translator.Result{Err: errors.New("Response empty")}
	}
	return &translator.Result{
		Text: resp.Result.Translated,
	}
}

func (ProviderMultillect) Name() string {
	return PROVIDER_CODE
}

type transfact struct{}

func (transfact) NewInstance(opts map[string]interface{}) translator.Translator {
	result := &ProviderMultillect{}
	//TODO: check apikey, AccountId
	result.apikey = opts["apikey"].(string) //Secret key
	result.accountID = opts["AccountId"].(string)
	if result.accountID == "" {
		panic(errors.New("API Multillect need user ID"))
	}
	return translator.Translator(result)
}

func init() {
	translate.Register(PROVIDER_CODE, &transfact{})
}
