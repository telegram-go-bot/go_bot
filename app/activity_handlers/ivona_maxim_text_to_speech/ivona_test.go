package text_to_speech

import "testing"

func TestGetUrl(t *testing.T) {
	url, parsed := getUrl(`Ссылка для скачивания звукового файла (WAV): http:\\sport.ulrg.ru:4055\vk\9D676154-77FD-4ECC-848A-06570C575186.wav
	Ссылка для скачивания звукового файла (mp3): http:\\sport.ulrg.ru:4055\vk\9D676154-77FD-4ECC-848A-06570C575186.mp3
	Ссылка будет доступна в течение 60 минут
	   
	Введите текст, который будем озвучивать либо выберите команду`)
	if !parsed {
		t.Fatal("getUrl parsing failed")
	}

	expectedUrl := `http:\\sport.ulrg.ru:4055\vk\9D676154-77FD-4ECC-848A-06570C575186.mp3`
	if url != expectedUrl {
		t.Fatalf("getUrl returned invalid url. \nExpected: %s,\nActual: %s", expectedUrl, url)
	}
}
