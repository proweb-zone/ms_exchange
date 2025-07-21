package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"ms_exchange/internal/dto"
	"ms_exchange/pkg/transport"
	"net/http"
	"time"
)

func HandleError(w http.ResponseWriter, result any, errMsg string) {
	transport.Http(
		w, http.StatusBadRequest, dto.ErrorResponse{
			Error:  errMsg,
			Result: result,
		},
	)
}

func ParseBody(w http.ResponseWriter, r *http.Request) []byte {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		HandleError(w, nil, "ошибка чтения запроса")
		return nil
	}
	defer r.Body.Close()

	return body
}

func DoCurl(
	ctx context.Context,
	method,
	url string,
	headers map[string]string,
	reqBody any,
) ([]byte, error) {
	var jsonBody []byte
	var result []byte
	var err error

	if reqBody != nil {
		jsonBody, err = json.Marshal(reqBody)
		if err != nil {
			return result, errors.New("ошибка сериализации тела запроса")
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return result, errors.New("ошибка при формировании http-запроса")
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return result, errors.New("ошибка во время запроса к серверу")
	}
	defer resp.Body.Close()

	result, err = io.ReadAll(resp.Body)
	if err != nil {
		return result, errors.New("ошибка чтения ответа на http-запрос")
	}

	return result, err
}
