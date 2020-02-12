package xmlapi

import (
	g "balabanovds/go-stats/globs"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func xmlapiRequest(ip string, query []byte) ([]byte, error) {
	url := fmt.Sprintf("https://%s:8443%s",
		ip,
		"/xmlapi/invoke")
	//data := []byte(strings.TrimSpace(query))

	req, err := http.NewRequest("POST", url, bytes.NewReader(query))
	if err != nil {
		g.Debugf("ERROR: Create request object. %v", err.Error())
		return nil, errors.New(err.Error())
	}

	req.Header.Set("Content-type", "text/xml; charset=ISO-8859-1")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		g.Debugf("%s:ERROR Dispatch request. %v", tagHTTP, err.Error())
		return nil, errors.New(err.Error())
	}
	// io.Copy(os.Stdout, res.Body)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil,
			errors.New("Service on server " + ip +
				" returned status: " + strconv.Itoa(res.StatusCode))
	}

	rawResponce, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return rawResponce, nil
}
