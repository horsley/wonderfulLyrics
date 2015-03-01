package xiami

import (
	//"io/ioutil"
	//"net/http"
	"testing"
)

func TestMain(t *testing.T) {
	url := getSongUrl(1773862701)
	t.Log(url)
	//resp, err := http.Get(url)
	//if err != nil {
	//	t.Error(err)
	//}
	//defer resp.Body.Close()

	//t.Log(resp.ContentLength)
	//bin, _ := ioutil.ReadAll(resp.Body)
	//ioutil.WriteFile("test.mp3", bin, 0666)
}
